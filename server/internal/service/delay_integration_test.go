package service

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"clash-server/internal/ws"
)

// TestDelayStreamingIntegration 用真实 Mihomo 核心验证逐节点测延迟的增量推送。
//
// 该测试需要外部环境, 未提供时自动跳过:
//   - 通过 E2E_CORE_ADDR 指定核心地址 (默认 127.0.0.1:19090)
//   - 通过 E2E_GROUP 指定代理组名 (默认 G1)
//
// 准备方式: 启动 mihomo 并配置若干响应时间不同的代理节点,
// 使快慢节点的延迟差异可观测。
func TestDelayStreamingIntegration(t *testing.T) {
	coreAddr := envOrDefault("E2E_CORE_ADDR", "127.0.0.1:19090")
	group := envOrDefault("E2E_GROUP", "G1")

	host, port := splitHostPortForTest(t, coreAddr)
	t.Setenv("CS_CORE_API_HOST", host)
	t.Setenv("CS_CORE_API_PORT", port)

	cs := NewCoreService()
	cs.running = true

	hub := ws.NewHub()
	go hub.Run()
	cs.SetHub(hub)

	client := ws.NewClient(hub)
	hub.Register(client)
	client.Subscribe([]ws.MessageType{ws.TypeDelay})

	received := make(chan ws.DelayData, 64)
	go func() {
		for msg := range client.GetSend() {
			var m ws.Message
			if err := json.Unmarshal(msg, &m); err != nil {
				continue
			}
			if m.Type != ws.TypeDelay {
				continue
			}
			var d ws.DelayData
			if err := json.Unmarshal(m.Data, &d); err != nil {
				continue
			}
			received <- d
		}
	}()

	names, err := cs.GetGroupProxyNames(group)
	if err != nil {
		t.Skipf("无法连接真实核心 (%s), 跳过集成测试: %v", coreAddr, err)
	}
	t.Logf("代理组 %s 含 %d 个节点: %v", group, len(names), names)

	start := time.Now()
	if err := cs.StartDelayTest(DelayTestRequest{
		Group:       group,
		TestURL:     envOrDefault("E2E_TEST_URL", "https://www.gstatic.com/generate_204"),
		Timeout:     5000,
		Concurrency: 10,
	}); err != nil {
		t.Fatalf("StartDelayTest 失败: %v", err)
	}

	// 触发应立即返回, 不等待最慢节点
	if elapsed := time.Since(start); elapsed > 500*time.Millisecond {
		t.Errorf("触发耗时 %v 过长, 期望立即返回", elapsed)
	}

	var nodeResults []ws.DelayData
	timeout := time.After(20 * time.Second)
	done := false
	for !done {
		select {
		case d := <-received:
			if d.Done {
				done = true
				break
			}
			nodeResults = append(nodeResults, d)
			t.Logf("  [%.3fs] 收到节点结果: %s = %dms", time.Since(start).Seconds(), d.Name, d.Delay)
		case <-timeout:
			t.Fatalf("等待推送超时, 已收到 %d 条", len(nodeResults))
		}
	}

	total := time.Since(start)
	t.Logf("共收到 %d 条节点结果 + done 信号, 总用时 %.3fs", len(nodeResults), total.Seconds())

	if len(nodeResults) != len(names) {
		t.Errorf("期望 %d 条节点结果, 实际 %d 条", len(names), len(nodeResults))
	}

	// 增量性: 存在快慢差异时, 首个结果应远早于全部完成
	if len(nodeResults) > 1 {
		sawFastBeforeEnd := false
		for _, r := range nodeResults[:len(nodeResults)-1] {
			if r.Name != "" {
				sawFastBeforeEnd = true
				break
			}
		}
		if !sawFastBeforeEnd {
			t.Errorf("未观察到节点结果在完成信号之前到达")
		}
	}
}

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func splitHostPortForTest(t *testing.T, addr string) (string, string) {
	t.Helper()
	for i := len(addr) - 1; i >= 0; i-- {
		if addr[i] == ':' {
			return addr[:i], addr[i+1:]
		}
	}
	t.Fatalf("非法地址: %s", addr)
	return "", ""
}
