package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"clash-server/internal/ws"
)

// delayTestServer 模拟 Mihomo 的延迟测试 API。
type delayTestServer struct {
	*httptest.Server
	// nodeDelays 每个节点的模拟响应耗时(毫秒)
	nodeDelays map[string]int
	// groupMembers 代理组内的节点列表
	groupMembers []string
	// failNodes 需要返回错误的节点
	failNodes map[string]bool

	inflight    int32
	maxInflight int32
	totalCalls  int32
	mu          sync.Mutex
	callOrder   []string
}

func newDelayTestServer(nodeDelays map[string]int, groupMembers []string, failNodes map[string]bool) *delayTestServer {
	s := &delayTestServer{
		nodeDelays:   nodeDelays,
		groupMembers: groupMembers,
		failNodes:    failNodes,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/proxies/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		// /proxies/{group}
		if r.URL.Query().Get("timeout") == "" {
			name := path[len("/proxies/"):]
			s.mu.Lock()
			members := s.groupMembers
			s.mu.Unlock()
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"all": members, "name": name})
			return
		}
		// /proxies/{name}/delay
		name := path[len("/proxies/"):]
		name = name[:len(name)-len("/delay")]

		cur := atomic.AddInt32(&s.inflight, 1)
		for {
			max := atomic.LoadInt32(&s.maxInflight)
			if cur <= max || atomic.CompareAndSwapInt32(&s.maxInflight, max, cur) {
				break
			}
		}
		defer atomic.AddInt32(&s.inflight, -1)
		atomic.AddInt32(&s.totalCalls, 1)

		s.mu.Lock()
		s.callOrder = append(s.callOrder, name)
		s.mu.Unlock()

		delay := s.nodeDelays[name]
		if delay > 0 {
			time.Sleep(time.Duration(delay) * time.Millisecond)
		}
		if s.failNodes[name] {
			w.WriteHeader(http.StatusBadGateway)
			_, _ = w.Write([]byte(`{"message":"delay test failed"}`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"delay": delay})
	})
	s.Server = httptest.NewServer(mux)
	return s
}

// newTestCoreService 构造指向测试服务器且处于运行状态的 CoreService。
func newTestCoreService(t *testing.T, addr string) *CoreService {
	t.Helper()
	u, err := url.Parse(addr)
	if err != nil {
		t.Fatalf("解析测试服务器地址失败: %v", err)
	}
	port, err := strconv.Atoi(u.Port())
	if err != nil {
		t.Fatalf("解析测试服务器端口失败: %v", err)
	}

	t.Setenv("CS_CORE_API_HOST", u.Hostname())
	t.Setenv("CS_CORE_API_PORT", strconv.Itoa(port))

	cs := NewCoreService()
	cs.running = true
	return cs
}

// TestRunDelayTestConcurrent 验证节点被并发测试: 总耗时接近最慢节点而非所有节点之和。
func TestRunDelayTestConcurrent(t *testing.T) {
	nodeDelays := map[string]int{
		"fast": 100,
		"mid":  400,
		"slow": 800,
	}
	members := []string{"fast", "mid", "slow"}

	srv := newDelayTestServer(nodeDelays, members, nil)
	defer srv.Close()

	cs := newTestCoreService(t, srv.URL)

	var mu sync.Mutex
	var results []DelayTestResult

	start := time.Now()
	err := cs.RunDelayTest(context.Background(), DelayTestRequest{
		Group:       "G1",
		Timeout:     5000,
		Concurrency: 10,
	}, func(r DelayTestResult) {
		mu.Lock()
		results = append(results, r)
		mu.Unlock()
	})
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("RunDelayTest 返回错误: %v", err)
	}
	if len(results) != 3 {
		t.Fatalf("期望 3 个结果, 实际 %d 个", len(results))
	}

	// 并发时总耗时应接近最慢节点(800ms); 串行则需 1300ms
	if elapsed > 1200*time.Millisecond {
		t.Errorf("总耗时 %v 过长, 期望接近最慢节点的 800ms (说明未并发)", elapsed)
	}
	t.Logf("并发总耗时: %v (最慢节点 800ms)", elapsed)
}

// TestRunDelayTestIncrementalOrder 验证快节点的先返回, 而非等待全部完成。
func TestRunDelayTestIncrementalOrder(t *testing.T) {
	nodeDelays := map[string]int{
		"fast": 50,
		"mid":  300,
		"slow": 700,
	}
	members := []string{"slow", "mid", "fast"}

	srv := newDelayTestServer(nodeDelays, members, nil)
	defer srv.Close()

	cs := newTestCoreService(t, srv.URL)

	var mu sync.Mutex
	var order []string
	var firstArrival time.Duration
	var lastArrival time.Duration
	start := time.Now()

	err := cs.RunDelayTest(context.Background(), DelayTestRequest{
		Group:       "G1",
		Timeout:     5000,
		Concurrency: 10,
	}, func(r DelayTestResult) {
		mu.Lock()
		order = append(order, r.Name)
		if len(order) == 1 {
			firstArrival = time.Since(start)
		}
		lastArrival = time.Since(start)
		mu.Unlock()
	})
	if err != nil {
		t.Fatalf("RunDelayTest 返回错误: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()
	if len(order) != 3 {
		t.Fatalf("期望 3 个结果, 实际 %d 个", len(order))
	}
	// 第一个结果应在快节点完成后即到达, 而非等待全部完成
	if firstArrival > 300*time.Millisecond {
		t.Errorf("首个结果在 %v 才到达, 期望快节点完成后立即到达 (约 50ms)", firstArrival)
	}
	if lastArrival-firstArrival < 400*time.Millisecond {
		t.Errorf("结果到达间隔 %v 过小, 期望节点间存在明显时间差", lastArrival-firstArrival)
	}
	t.Logf("结果到达顺序: %v, 首个 %v, 末个 %v", order, firstArrival, lastArrival)
}

// TestRunDelayTestConcurrencyLimit 验证并发数不超过配置上限。
func TestRunDelayTestConcurrencyLimit(t *testing.T) {
	nodeDelays := map[string]int{}
	members := []string{}
	for i := 0; i < 20; i++ {
		name := fmt.Sprintf("node-%02d", i)
		nodeDelays[name] = 200
		members = append(members, name)
	}

	srv := newDelayTestServer(nodeDelays, members, nil)
	defer srv.Close()

	cs := newTestCoreService(t, srv.URL)

	err := cs.RunDelayTest(context.Background(), DelayTestRequest{
		Group:       "G1",
		Timeout:     5000,
		Concurrency: 5,
	}, func(r DelayTestResult) {})
	if err != nil {
		t.Fatalf("RunDelayTest 返回错误: %v", err)
	}

	maxSeen := atomic.LoadInt32(&srv.maxInflight)
	if maxSeen > 5 {
		t.Errorf("最大并发数 %d 超过配置上限 5", maxSeen)
	}
	totalCalls := atomic.LoadInt32(&srv.totalCalls)
	if totalCalls != 20 {
		t.Errorf("期望发起 20 次测速调用, 实际 %d 次", totalCalls)
	}
	t.Logf("最大并发数: %d (上限 5), 总调用: %d", maxSeen, totalCalls)
}

// TestRunDelayTestNodeFailure 验证单节点失败不影响同组其他节点。
func TestRunDelayTestNodeFailure(t *testing.T) {
	nodeDelays := map[string]int{
		"ok1":   50,
		"bad":   50,
		"ok2":   50,
	}
	members := []string{"ok1", "bad", "ok2"}
	failNodes := map[string]bool{"bad": true}

	srv := newDelayTestServer(nodeDelays, members, failNodes)
	defer srv.Close()

	cs := newTestCoreService(t, srv.URL)

	var mu sync.Mutex
	var successCount, failCount int
	err := cs.RunDelayTest(context.Background(), DelayTestRequest{
		Group:       "G1",
		Timeout:     5000,
		Concurrency: 10,
	}, func(r DelayTestResult) {
		mu.Lock()
		defer mu.Unlock()
		if r.Err != nil {
			failCount++
			if r.Name != "bad" {
				t.Errorf("非预期节点报错: %s", r.Name)
			}
		} else {
			successCount++
		}
	})
	if err != nil {
		t.Fatalf("RunDelayTest 返回错误: %v", err)
	}
	if successCount != 2 {
		t.Errorf("期望 2 个成功结果, 实际 %d 个", successCount)
	}
	if failCount != 1 {
		t.Errorf("期望 1 个失败结果, 实际 %d 个", failCount)
	}
}

// TestRunDelayTestCancel 验证取消上下文后任务停止。
func TestRunDelayTestCancel(t *testing.T) {
	nodeDelays := map[string]int{}
	members := []string{}
	for i := 0; i < 10; i++ {
		name := fmt.Sprintf("node-%02d", i)
		nodeDelays[name] = 1000
		members = append(members, name)
	}

	srv := newDelayTestServer(nodeDelays, members, nil)
	defer srv.Close()

	cs := newTestCoreService(t, srv.URL)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	start := time.Now()
	err := cs.RunDelayTest(ctx, DelayTestRequest{
		Group:       "G1",
		Timeout:     5000,
		Concurrency: 2,
	}, func(r DelayTestResult) {})
	elapsed := time.Since(start)

	if err == nil {
		t.Errorf("期望取消后返回错误, 实际返回 nil")
	}
	// 取消后不应等待全部 10 个节点(每个 1 秒)跑完
	if elapsed > 5*time.Second {
		t.Errorf("取消后耗时 %v 过长", elapsed)
	}
	t.Logf("取消后耗时: %v, 错误: %v", elapsed, err)
}

// TestDelayResultBroadcast 验证结果通过 hub 按 delay 类型推送给订阅者。
func TestDelayResultBroadcast(t *testing.T) {
	hub := ws.NewHub()
	go hub.Run()

	cs := NewCoreService()
	cs.SetHub(hub)

	client := ws.NewClient(hub)
	hub.Register(client)
	client.Subscribe([]ws.MessageType{ws.TypeDelay})

	received := make(chan ws.DelayData, 8)
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

	cs.BroadcastDelayResult(DelayTestResult{Group: "G1", Name: "HK-01", Delay: 120})
	cs.BroadcastDelayResult(DelayTestResult{Group: "G1", Name: "JP-02", Err: fmt.Errorf("timeout")})
	cs.BroadcastDelayDone("G1")

	for i := 0; i < 3; i++ {
		select {
		case d := <-received:
			switch i {
			case 0:
				if d.Name != "HK-01" || d.Delay != 120 || d.Done {
					t.Errorf("第 1 条消息异常: %+v", d)
				}
			case 1:
				if d.Name != "JP-02" || d.Error == "" || d.Done {
					t.Errorf("第 2 条消息异常: %+v", d)
				}
			case 2:
				if !d.Done || d.Group != "G1" {
					t.Errorf("第 3 条消息应为完成信号: %+v", d)
				}
			}
		case <-time.After(2 * time.Second):
			t.Fatalf("等待第 %d 条推送超时", i+1)
		}
	}
}

// TestDelayResultNotBroadcastToUnsubscribed 验证未订阅 delay 类型的客户端收不到推送。
func TestDelayResultNotBroadcastToUnsubscribed(t *testing.T) {
	hub := ws.NewHub()
	go hub.Run()

	cs := NewCoreService()
	cs.SetHub(hub)

	client := ws.NewClient(hub)
	hub.Register(client)
	// 只订阅 traffic 类型
	client.Subscribe([]ws.MessageType{ws.TypeTraffic})

	received := make(chan ws.Message, 8)
	go func() {
		for msg := range client.GetSend() {
			var m ws.Message
			if err := json.Unmarshal(msg, &m); err == nil {
				received <- m
			}
		}
	}()

	cs.BroadcastDelayResult(DelayTestResult{Group: "G1", Name: "HK-01", Delay: 120})
	cs.BroadcastDelayDone("G1")

	select {
	case m := <-received:
		t.Errorf("未订阅 delay 的客户端收到了消息: %+v", m)
	case <-time.After(500 * time.Millisecond):
		// 符合预期: 未订阅则不推送
	}
}

// TestLifecycleCancelStopsDelayTest 验证核心停止时进行中的测延迟任务被终止。
func TestLifecycleCancelStopsDelayTest(t *testing.T) {
	nodeDelays := map[string]int{}
	members := []string{}
	for i := 0; i < 10; i++ {
		name := fmt.Sprintf("node-%02d", i)
		nodeDelays[name] = 1000
		members = append(members, name)
	}

	srv := newDelayTestServer(nodeDelays, members, nil)
	defer srv.Close()

	cs := newTestCoreService(t, srv.URL)

	// 启动一个长耗时的测速任务
	finished := make(chan error, 1)
	go func() {
		finished <- cs.RunDelayTest(cs.lifecycleContext(), DelayTestRequest{
			Group:       "G1",
			Timeout:     5000,
			Concurrency: 2,
		}, func(r DelayTestResult) {})
	}()

	time.Sleep(100 * time.Millisecond)
	// 核心停止 -> 生命周期上下文被取消
	cs.cancelLifecycle()

	select {
	case err := <-finished:
		if err == nil {
			t.Errorf("核心停止后测速任务应返回取消错误, 实际返回 nil")
		} else {
			t.Logf("核心停止后测速任务已终止: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatalf("核心停止后测速任务未在预期时间内终止")
	}
}
