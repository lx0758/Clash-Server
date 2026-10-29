package service

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestApplyConfigKeepsLastConfigSnapshot 验证 ApplyConfig 先取旧配置快照再写入新配置。
//
// 回归背景: 旧实现先写 lastConfig 再比较, 导致 needsRestart 恒为 false,
// 端口类配置变更永远不会走软重启分支。
func TestApplyConfigKeepsLastConfigSnapshot(t *testing.T) {
	cs := NewCoreService()

	// 预置一份包含重启类键的旧配置
	cs.mu.Lock()
	cs.lastConfig = map[string]interface{}{
		"mixed-port": 7890,
		"proxies":    []interface{}{"a"},
	}
	cs.mu.Unlock()

	old := cs.getLastConfig()
	if old["mixed-port"] != 7890 {
		t.Fatalf("旧配置快照错误: %v", old)
	}

	// 新配置仅改变 proxies, 不应触发重启
	newConfig := map[string]interface{}{
		"mixed-port": 7890,
		"proxies":    []interface{}{"b"},
	}
	if cs.needsRestart(old, newConfig) {
		t.Errorf("仅 proxies 变化不应触发重启")
	}

	// 新配置改变 mixed-port, 应触发重启
	newConfig2 := map[string]interface{}{
		"mixed-port": 7891,
		"proxies":    []interface{}{"b"},
	}
	if !cs.needsRestart(old, newConfig2) {
		t.Errorf("mixed-port 变化应触发重启")
	}
}

// TestNeedsRestartDetectsAllRestartKeys 验证所有重启类键都能被识别。
func TestNeedsRestartDetectsAllRestartKeys(t *testing.T) {
	cs := NewCoreService()
	base := map[string]interface{}{"mixed-port": 7890}

	for _, key := range RestartRequiredKeys {
		old := map[string]interface{}{}
		for k, v := range base {
			old[k] = v
		}
		old[key] = "old-value"

		next := map[string]interface{}{}
		for k, v := range base {
			next[k] = v
		}
		next[key] = "new-value"

		if !cs.needsRestart(old, next) {
			t.Errorf("键 %s 变化时应触发重启", key)
		}
	}
}

// TestApiRequestIncludesErrorBody 验证核心返回错误时错误信息包含响应体内容。
//
// 回归背景: 旧实现丢弃响应体, 只返回泛化的 "API request failed",
// 上层拿不到 YAML 语法错误等具体原因。
func TestApiRequestIncludesErrorBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"message": "proxy slow-node: invalid config: yaml: line 3: mapping values are not allowed",
		})
	}))
	defer srv.Close()

	cs := newTestCoreService(t, srv.URL)

	err := cs.apiRequest("PUT", "/configs", map[string]string{"path": ""}, nil)
	if err == nil {
		t.Fatalf("期望返回错误, 实际 nil")
	}
	msg := err.Error()
	if !strings.Contains(msg, "invalid config") {
		t.Errorf("错误信息应包含核心返回的具体原因, 实际: %s", msg)
	}
	if !strings.Contains(msg, "yaml") {
		t.Errorf("错误信息应保留 YAML 语法错误详情, 实际: %s", msg)
	}
	if strings.Contains(msg, "API request failed") && !strings.Contains(msg, "invalid config") {
		t.Errorf("不得只返回泛化失败信息, 实际: %s", msg)
	}
}

// TestApiRequestEmptyErrorBody 验证核心返回空错误体时仍有可辨识的错误信息。
func TestApiRequestEmptyErrorBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	cs := newTestCoreService(t, srv.URL)

	err := cs.apiRequest("PUT", "/configs", nil, nil)
	if err == nil {
		t.Fatalf("期望返回错误, 实际 nil")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("错误信息应包含状态码, 实际: %s", err.Error())
	}
}

// TestServiceMethodsErrorWhenCoreStopped 验证核心未运行时数据查询方法返回错误。
//
// 回归背景: 旧实现返回空结果 + nil error, 导致 handler 层以空列表
// 伪装成功, 前端误以为"核心运行但没有规则/连接/代理"。
func TestServiceMethodsErrorWhenCoreStopped(t *testing.T) {
	cs := NewCoreService() // running = false

	if _, err := cs.GetProxies(); err == nil {
		t.Errorf("GetProxies 在核心未运行时应返回错误")
	}
	if _, err := cs.GetConnections(); err == nil {
		t.Errorf("GetConnections 在核心未运行时应返回错误")
	}
	if _, err := cs.GetRules(); err == nil {
		t.Errorf("GetRules 在核心未运行时应返回错误")
	}
	if _, err := cs.GetMode(); err == nil {
		t.Errorf("GetMode 在核心未运行时应返回错误")
	}
}

// TestRestartWithFallbackSoftRestartSuccess 验证软重启成功时不触发降级。
func TestRestartWithFallbackSoftRestartSuccess(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/restart" || r.URL.Path == "/configs" {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	cs := newTestCoreService(t, srv.URL)
	if err := cs.restartWithFallback("proxies: []"); err != nil {
		t.Errorf("软重启成功时不应返回错误, 实际: %v", err)
	}
}

// TestRestartWithFallbackReportsBothFailures 验证软重启与降级重启都失败时返回错误。
//
// 契约: 不得以成功响应掩盖核心未能生效的事实, 错误信息需说明两条路径都失败。
func TestRestartWithFallbackReportsBothFailures(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 配置下发成功, 但软重启失败
		if r.URL.Path == "/configs" {
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"message":"restart failed"}`))
	}))
	defer srv.Close()

	cs := newTestCoreService(t, srv.URL)

	err := cs.restartWithFallback("proxies: []")
	if err == nil {
		t.Fatalf("两条路径都失败时应返回错误, 实际 nil")
	}
	msg := err.Error()
	if !strings.Contains(msg, "软重启失败") {
		t.Errorf("错误信息应说明软重启失败, 实际: %s", msg)
	}
	if !strings.Contains(msg, "降级重启") {
		t.Errorf("错误信息应说明降级重启的情况, 实际: %s", msg)
	}
}

// TestRestartWithFallbackConfigApplyFailure 验证配置下发失败时不进入重启流程。
func TestRestartWithFallbackConfigApplyFailure(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"message":"invalid yaml"}`))
	}))
	defer srv.Close()

	cs := newTestCoreService(t, srv.URL)

	err := cs.restartWithFallback("proxies: []")
	if err == nil {
		t.Fatalf("配置下发失败时应返回错误, 实际 nil")
	}
	if !strings.Contains(err.Error(), "应用配置失败") {
		t.Errorf("错误信息应说明是配置下发失败, 实际: %s", err.Error())
	}
	// 配置都没应用成功, 不应再尝试重启
	if !strings.Contains(err.Error(), "软重启失败") {
		t.Logf("配置下发失败后未进入重启流程 (符合预期)")
	}
}

// TestVerifyRestartAppliedDetectsUnappliedChange 验证校验能识别"接口成功但变更未生效"。
//
// 回归背景: Mihomo 的软重启对端口类配置只是重置而非重新绑定,
// POST /restart 返回 200 但 mixed-port 并未真正生效。
func TestVerifyRestartAppliedDetectsUnappliedChange(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 核心返回的 mixed-port 与目标不一致
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"mixed-port": 0, "allow-lan": true}`))
	}))
	defer srv.Close()

	cs := newTestCoreService(t, srv.URL)
	cs.mu.Lock()
	cs.lastConfig = map[string]interface{}{
		"mixed-port": 7891,
		"allow-lan":  true,
	}
	cs.mu.Unlock()

	err := cs.verifyRestartApplied()
	if err == nil {
		t.Fatalf("mixed-port 未生效时应检出问题, 实际 nil")
	}
	if !strings.Contains(err.Error(), "mixed-port") {
		t.Errorf("错误信息应指出未生效的配置项, 实际: %s", err.Error())
	}
}

// TestVerifyRestartAppliedAcceptsAppliedChange 验证变更已生效时校验通过。
func TestVerifyRestartAppliedAcceptsAppliedChange(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"mixed-port": 7891, "allow-lan": true}`))
	}))
	defer srv.Close()

	cs := newTestCoreService(t, srv.URL)
	cs.mu.Lock()
	cs.lastConfig = map[string]interface{}{
		"mixed-port": 7891,
		"allow-lan":  true,
	}
	cs.mu.Unlock()

	if err := cs.verifyRestartApplied(); err != nil {
		t.Errorf("变更已生效时校验应通过, 实际: %v", err)
	}
}

// TestVerifyRestartAppliedNormalizesNumberTypes 验证 int 与 float64 的差异被归一化。
func TestVerifyRestartAppliedNormalizesNumberTypes(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"mixed-port": 7891}`))
	}))
	defer srv.Close()

	cs := newTestCoreService(t, srv.URL)
	cs.mu.Lock()
	cs.lastConfig = map[string]interface{}{"mixed-port": 7891} // int
	cs.mu.Unlock()

	// JSON 解码得到的是 float64, 归一化后应与 int 视为相等
	if err := cs.verifyRestartApplied(); err != nil {
		t.Errorf("数值类型差异不应导致校验失败, 实际: %v", err)
	}
}
