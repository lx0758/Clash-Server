package service

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"clash-server/internal/config"
	"clash-server/internal/model"
	"clash-server/internal/repository"
)

// setupTestDB 初始化一个隔离的内存数据库并让全局 merger 指向它。
//
// 返回一个清理函数, 用于恢复原有的全局状态。
func setupTestDB(t *testing.T) func() {
	t.Helper()

	oldDB := model.DB
	oldMerger := mergerInstance

	if err := model.InitDB(":memory:"); err != nil {
		t.Fatalf("初始化测试数据库失败: %v", err)
	}

	// 让 merger 绑定到测试数据库。
	// 直接替换 mergerInstance 而不重置 sync.Once, 避免拷贝含锁的值。
	GetMergerService()
	mergerInstance = &MergerService{
		subRepo:           repository.NewSubscriptionRepository(),
		customizationRepo: repository.NewCustomizationRepository(),
	}

	cleanup := func() {
		model.DB = oldDB
		mergerInstance = oldMerger
	}
	return cleanup
}

// TestApplyConfigMergeFailureReturnsError 验证合并失败时 ApplyConfig 返回错误。
//
// 回归背景: 旧实现在合并失败时静默回退到最小配置并继续下发,
// 调用方以为配置已生效, 实际核心运行的是完全不同的配置。
func TestApplyConfigMergeFailureReturnsError(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	// 插入一个激活的订阅, 其内容是非法 YAML
	sub := &model.Subscription{
		Name:       "坏订阅",
		SourceType: model.SourceTypeLocal,
		Content:    "proxies:\n  - name: broken\n   bad-indent: [\n",
		Active:     true,
	}
	if err := repository.NewSubscriptionRepository().Create(sub); err != nil {
		t.Fatalf("创建测试订阅失败: %v", err)
	}

	cs := NewCoreService()

	// 预置 lastConfig, 用于验证合并失败时它保持不变
	cs.mu.Lock()
	cs.lastConfig = map[string]interface{}{"marker": "original"}
	cs.mu.Unlock()

	result := cs.ApplyConfig()
	if result.Error == nil {
		t.Fatalf("合并失败时 ApplyConfig 应返回错误, 实际 nil")
	}
	if !strings.Contains(result.Error.Error(), "合并配置失败") {
		t.Errorf("错误信息应说明是合并失败, 实际: %v", result.Error)
	}
	if result.HotReloaded || result.Restarted {
		t.Errorf("合并失败时不得声称配置已生效: %+v", result)
	}

	// 合并失败不得改变已持久化的配置状态
	if cs.getLastConfig()["marker"] != "original" {
		t.Errorf("合并失败时 lastConfig 应保持原样, 实际: %v", cs.getLastConfig())
	}
}

// TestApplyConfigMergeFailureDoesNotTouchCore 验证合并失败时不会调用核心接口。
func TestApplyConfigMergeFailureDoesNotTouchCore(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	sub := &model.Subscription{
		Name:       "坏订阅",
		SourceType: model.SourceTypeLocal,
		Content:    "proxies: [unclosed\n",
		Active:     true,
	}
	if err := repository.NewSubscriptionRepository().Create(sub); err != nil {
		t.Fatalf("创建测试订阅失败: %v", err)
	}

	// 指向一个会在被调用时立刻失败的测试服务器
	cs := newTestCoreService(t, "http://127.0.0.1:1")

	result := cs.ApplyConfig()
	if result.Error == nil {
		t.Fatalf("合并失败时 ApplyConfig 应返回错误, 实际 nil")
	}
	// 错误信息应来自合并阶段, 而不是核心请求阶段
	if strings.Contains(result.Error.Error(), "核心 API 请求失败") {
		t.Errorf("合并失败时不应发起核心请求, 实际错误: %v", result.Error)
	}
}

// TestApplyConfigSuccessDoesNotRestartOnProxyChange 验证仅节点变化时走热重载。
func TestApplyConfigSuccessDoesNotRestartOnProxyChange(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	sub := &model.Subscription{
		Name:       "好订阅",
		SourceType: model.SourceTypeLocal,
		Content:    "proxies:\n  - {name: node-a, type: http, server: 127.0.0.1, port: 1234}\n",
		Active:     true,
	}
	if err := repository.NewSubscriptionRepository().Create(sub); err != nil {
		t.Fatalf("创建测试订阅失败: %v", err)
	}

	// 预置一份旧配置: 核心相关键与当前一致, 仅 proxies 不同
	// 指向一个接受热重载的 mock 核心
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()
	cs := newTestCoreService(t, srv.URL)

	// 必须在设置核心地址的环境变量之后再读取核心配置, 否则快照与实际下发值不一致
	coreCfg := config.GetCoreConfig()

	cs.mu.Lock()
	cs.lastConfig = map[string]interface{}{
		"proxies":             []interface{}{map[string]interface{}{"name": "old-node"}},
		"mixed-port":          coreCfg.MixedPort,
		"allow-lan":           coreCfg.AllowLan,
		"mode":                coreCfg.Mode,
		"log-level":           coreCfg.LogLevel,
		"ipv6":                coreCfg.IPv6,
		"external-controller": fmt.Sprintf("%s:%d", coreCfg.APIHost, coreCfg.APIPort),
	}
	cs.mu.Unlock()

	result := cs.ApplyConfig()
	if result.Error != nil {
		t.Fatalf("ApplyConfig 应成功: %v", result.Error)
	}
	if result.Restarted {
		t.Errorf("仅节点变化不应触发重启")
	}
	if !result.HotReloaded {
		t.Errorf("仅节点变化应走热重载路径")
	}
}
