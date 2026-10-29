package service

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"sync"
	"sync/atomic"
	"time"

	"clash-server/internal/config"
	"clash-server/internal/ws"
)

// DefaultDelayConcurrency 逐节点测延迟的默认并发上限。
// Mihomo 核心本身并行处理测速请求, 这里的限流用于避免节点数量大时
// 对核心产生瞬时压力。
const DefaultDelayConcurrency = 10

// DelayTestResult 表示单个节点的测延迟结果。
type DelayTestResult struct {
	Group string
	Name  string
	Delay int
	Err   error
}

// DelayTestRequest 描述一次逐节点测延迟任务。
type DelayTestRequest struct {
	Group       string
	TestURL     string
	Timeout     int
	Concurrency int
}

// GetGroupProxyNames 返回指定代理组内的节点名称列表。
func (cs *CoreService) GetGroupProxyNames(groupName string) ([]string, error) {
	if !cs.IsRunning() {
		return nil, fmt.Errorf("core is not running")
	}
	var result struct {
		All []string `json:"all"`
	}
	path := fmt.Sprintf("/proxies/%s", url.PathEscape(groupName))
	if err := cs.apiRequest("GET", path, nil, &result); err != nil {
		return nil, err
	}
	return result.All, nil
}

// RunDelayTest 对目标组内的节点并发执行延迟测试, 每完成一个节点即通过
// onResult 回调。该方法阻塞直到全部节点测试完成或 ctx 被取消。
func (cs *CoreService) RunDelayTest(ctx context.Context, req DelayTestRequest, onResult func(DelayTestResult)) error {
	if !cs.IsRunning() {
		return fmt.Errorf("core is not running")
	}

	names, err := cs.GetGroupProxyNames(req.Group)
	if err != nil {
		return err
	}
	if len(names) == 0 {
		return fmt.Errorf("代理组 %s 内没有可测试的节点", req.Group)
	}

	concurrency := req.Concurrency
	if concurrency <= 0 {
		concurrency = config.GetCoreConfig().DelayConcurrency
	}
	if concurrency <= 0 {
		concurrency = DefaultDelayConcurrency
	}
	if concurrency > len(names) {
		concurrency = len(names)
	}

	// 带缓冲 channel 作为信号量, 限制同时打到核心的测速请求数
	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup
	// canceled 记录任务是否因上下文取消而中止
	var canceled atomic.Bool

	for _, name := range names {
		select {
		case <-ctx.Done():
			// 已取消, 等待在途请求收尾后返回
			canceled.Store(true)
			wg.Wait()
			return ctx.Err()
		default:
		}

		wg.Add(1)
		go func(proxyName string) {
			defer wg.Done()

			select {
			case sem <- struct{}{}:
				defer func() { <-sem }()
			case <-ctx.Done():
				canceled.Store(true)
				return
			}

			delay, err := cs.CheckDelay(proxyName, req.TestURL, req.Timeout)
			if err != nil {
				// 节点测试失败不影响同组其他节点, 单独回传错误
				onResult(DelayTestResult{Group: req.Group, Name: proxyName, Err: err})
				return
			}
			onResult(DelayTestResult{Group: req.Group, Name: proxyName, Delay: delay})
		}(name)
	}

	wg.Wait()
	if canceled.Load() {
		return ctx.Err()
	}
	return nil
}

// BroadcastDelayResult 将单个节点的延迟结果推送给订阅了 delay 类型的客户端。
func (cs *CoreService) BroadcastDelayResult(result DelayTestResult) {
	cs.mu.RLock()
	hub := cs.hub
	cs.mu.RUnlock()
	if hub == nil {
		return
	}

	data := ws.DelayData{Group: result.Group, Name: result.Name, Delay: result.Delay}
	if result.Err != nil {
		data.Error = result.Err.Error()
	}
	hub.BroadcastToType(ws.TypeDelay, data)
}

// BroadcastDelayDone 通知客户端指定代理组的测延迟任务已完成。
func (cs *CoreService) BroadcastDelayDone(group string) {
	cs.mu.RLock()
	hub := cs.hub
	cs.mu.RUnlock()
	if hub == nil {
		return
	}
	hub.BroadcastToType(ws.TypeDelay, ws.DelayData{Group: group, Done: true})
}

// StartDelayTest 启动一次逐节点测延迟任务并把结果推送给客户端。
// 该方法立即返回, 任务在后台进行, 生命周期跟随核心。
func (cs *CoreService) StartDelayTest(req DelayTestRequest) error {
	if !cs.IsRunning() {
		return fmt.Errorf("core is not running")
	}
	names, err := cs.GetGroupProxyNames(req.Group)
	if err != nil {
		return err
	}
	if len(names) == 0 {
		return fmt.Errorf("代理组 %s 内没有可测试的节点", req.Group)
	}

	ctx := cs.lifecycleContext()

	go func() {
		start := time.Now()
		runErr := cs.RunDelayTest(ctx, req, cs.BroadcastDelayResult)
		if runErr != nil {
			log.Printf("[CoreService] Delay test for group %s stopped: %v", req.Group, runErr)
		}
		log.Printf("[CoreService] Delay test for group %s finished in %v", req.Group, time.Since(start))
		// 无论正常结束还是被取消, 都要通知前端收敛状态
		cs.BroadcastDelayDone(req.Group)
	}()

	return nil
}
