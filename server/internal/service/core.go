package service

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"time"

	"clash-server/internal/config"
	"clash-server/internal/ws"

	"github.com/gorilla/websocket"
)

var RestartRequiredKeys = []string{
	"mixed-port",
	"external-controller",
	"secret",
	"bind-address",
	"tun",
	"interface-name",
	"routing-mark",
	"ipv6",
}

type CoreStatus struct {
	Running bool
	Version string
	Error   string
}

type ApplyConfigResult struct {
	Restarted   bool
	HotReloaded bool
	Error       error
}

type CoreService struct {
	running        bool
	autoRestart    bool
	version        string
	lastError      string
	lastConfig     map[string]interface{}
	mu             sync.RWMutex
	cmd            *exec.Cmd
	client         *http.Client
	onStatusChange func(status CoreStatus)
	hub            *ws.Hub
	wsConns        map[string]*websocket.Conn
	wsStopChan     chan struct{}
	wsDialer       websocket.Dialer
	// lifecycleCtx 跟随核心生命周期, 核心停止时取消, 用于终止进行中的测延迟任务
	lifecycleCtx    context.Context
	lifecycleCancel context.CancelFunc
}

var globalCoreService *CoreService
var coreServiceOnce sync.Once

func GetCoreService() *CoreService {
	coreServiceOnce.Do(func() {
		globalCoreService = &CoreService{
			client:   &http.Client{Timeout: 10 * time.Second},
			wsConns:  make(map[string]*websocket.Conn),
			wsDialer: websocket.Dialer{HandshakeTimeout: 10 * time.Second},
		}
	})
	return globalCoreService
}

func NewCoreService() *CoreService {
	return &CoreService{
		client:   &http.Client{Timeout: 10 * time.Second},
		wsConns:  make(map[string]*websocket.Conn),
		wsDialer: websocket.Dialer{HandshakeTimeout: 10 * time.Second},
	}
}

func (cs *CoreService) SetHub(hub *ws.Hub) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	cs.hub = hub
}

// lifecycleContext 返回跟随核心生命周期的上下文, 核心停止时该上下文被取消。
func (cs *CoreService) lifecycleContext() context.Context {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	if cs.lifecycleCtx == nil {
		cs.lifecycleCtx, cs.lifecycleCancel = context.WithCancel(context.Background())
	}
	return cs.lifecycleCtx
}

// cancelLifecycle 取消生命周期上下文, 终止所有进行中的测延迟任务。
func (cs *CoreService) cancelLifecycle() {
	cs.mu.Lock()
	cancel := cs.lifecycleCancel
	cs.lifecycleCtx = nil
	cs.lifecycleCancel = nil
	cs.mu.Unlock()
	if cancel != nil {
		cancel()
	}
}

func (cs *CoreService) SetOnStatusChange(callback func(status CoreStatus)) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	cs.onStatusChange = callback
}

func (cs *CoreService) Start() error {
	return cs.start()
}

func (cs *CoreService) Stop() error {
	return cs.stop()
}

func (cs *CoreService) Restart() error {
	err := cs.stop()
	if err != nil {
		return err
	}
	err = cs.start()
	if err != nil {
		return err
	}
	return nil
}

func (cs *CoreService) GetStatus() CoreStatus {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	return CoreStatus{
		Running: cs.running,
		Version: cs.version,
		Error:   cs.lastError,
	}
}

func (cs *CoreService) GetVersion() string {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	return cs.version
}

func (cs *CoreService) IsRunning() bool {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	return cs.running
}

func (cs *CoreService) GetProxies() (map[string]interface{}, error) {
	if !cs.IsRunning() {
		// 核心未运行时返回错误, 避免调用方把"核心没跑"误读为"没有代理"
		return nil, fmt.Errorf("core is not running")
	}
	var result map[string]interface{}
	if err := cs.apiRequest("GET", "/proxies", nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (cs *CoreService) SelectProxy(group, name string) error {
	if !cs.IsRunning() {
		return fmt.Errorf("core is not running")
	}
	return cs.apiRequest("PUT", fmt.Sprintf("/proxies/%s", url.PathEscape(group)), map[string]string{"name": name})
}

func (cs *CoreService) GetConnections() ([]interface{}, error) {
	if !cs.IsRunning() {
		return nil, fmt.Errorf("core is not running")
	}
	var result struct {
		Connections []interface{} `json:"connections"`
	}
	if err := cs.apiRequest("GET", "/connections", nil, &result); err != nil {
		return nil, err
	}
	return result.Connections, nil
}

func (cs *CoreService) CloseConnection(id string) error {
	if !cs.IsRunning() {
		return nil
	}
	return cs.apiRequest("DELETE", fmt.Sprintf("/connections/%s", url.PathEscape(id)), nil)
}

func (cs *CoreService) CloseAllConnections() error {
	if !cs.IsRunning() {
		return nil
	}
	return cs.apiRequest("DELETE", "/connections", nil)
}

func (cs *CoreService) GetTraffic() (map[string]interface{}, error) {
	if !cs.IsRunning() {
		return nil, nil
	}
	var result map[string]interface{}
	if err := cs.apiRequest("GET", "/traffic", nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (cs *CoreService) GetRules() ([]interface{}, error) {
	if !cs.IsRunning() {
		return nil, fmt.Errorf("core is not running")
	}
	var result struct {
		Rules []interface{} `json:"rules"`
	}
	if err := cs.apiRequest("GET", "/rules", nil, &result); err != nil {
		return nil, err
	}
	return result.Rules, nil
}

func (cs *CoreService) CheckDelay(name string, testURL string, timeout int) (int, error) {
	if !cs.IsRunning() {
		return 0, fmt.Errorf("core is not running")
	}
	path := fmt.Sprintf("/proxies/%s/delay?timeout=%d", url.PathEscape(name), timeout)
	if testURL != "" {
		path += fmt.Sprintf("&url=%s", url.QueryEscape(testURL))
	}
	var result struct {
		Delay int `json:"delay"`
	}
	if err := cs.apiRequest("GET", path, nil, &result); err != nil {
		return 0, err
	}
	return result.Delay, nil
}

func (cs *CoreService) CheckGroupDelay(groupName string, testURL string, timeout int) error {
	if !cs.IsRunning() {
		return fmt.Errorf("core is not running")
	}
	path := fmt.Sprintf("/group/%s/delay?timeout=%d", url.PathEscape(groupName), timeout)
	if testURL != "" {
		path += fmt.Sprintf("&url=%s", url.QueryEscape(testURL))
	}
	return cs.apiRequest("GET", path, nil)
}

func (cs *CoreService) GetMode() (string, error) {
	if !cs.IsRunning() {
		return "", fmt.Errorf("core is not running")
	}
	var result struct {
		Mode string `json:"mode"`
	}
	if err := cs.apiRequest("GET", "/configs", nil, &result); err != nil {
		return "rule", err
	}
	return result.Mode, nil
}

func (cs *CoreService) SetMode(mode string) error {
	if !cs.IsRunning() {
		return fmt.Errorf("core is not running")
	}
	return cs.apiRequest("PATCH", "/configs", map[string]string{"mode": mode})
}

func (cs *CoreService) HotReload(yamlConfig string) error {
	payload := map[string]string{"path": "", "payload": yamlConfig}
	return cs.apiRequest("PUT", "/configs?force=true", payload, nil)
}

func (cs *CoreService) SoftRestart(yamlConfig string) error {
	payload := map[string]string{"path": "", "payload": yamlConfig}
	return cs.apiRequest("POST", "/restart", payload, nil)
}

func (cs *CoreService) apiRequest(method, path string, body interface{}, result ...interface{}) error {
	coreCfg := config.GetCoreConfig()
	apiURL := fmt.Sprintf("http://%s:%d%s", coreCfg.APIHost, coreCfg.APIPort, path)
	log.Printf("[CoreService] API request: %s %s", method, apiURL)
	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			log.Printf("[CoreService] Failed to marshal request body: %v", err)
			return err
		}
		reqBody = bytes.NewReader(data)
	}
	req, err := http.NewRequest(method, apiURL, reqBody)
	if err != nil {
		log.Printf("[CoreService] Failed to create request: %v", err)
		return err
	}
	if coreCfg.APISecret != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", coreCfg.APISecret))
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := cs.client.Do(req)
	if err != nil {
		log.Printf("[CoreService] API request failed: %v", err)
		return err
	}
	defer resp.Body.Close()
	log.Printf("[CoreService] API response: %s %s -> %d", method, path, resp.StatusCode)
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		// 读取核心返回的错误正文, 便于上层透传具体原因
		// (如 YAML 语法错误), 而不是只给出泛化的失败信息
		errBody, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			log.Printf("[CoreService] Failed to read error response body: %v", readErr)
		}
		detail := strings.TrimSpace(string(errBody))
		if detail == "" {
			return fmt.Errorf("核心 API 请求失败: %s %s 返回 %d", method, path, resp.StatusCode)
		}
		return fmt.Errorf("核心 API 请求失败: %s %s 返回 %d: %s", method, path, resp.StatusCode, detail)
	}
	if len(result) > 0 && result[0] != nil {
		return json.NewDecoder(resp.Body).Decode(result[0])
	}
	return nil
}

func (cs *CoreService) notifyStatusChange() {
	cs.mu.RLock()
	callback := cs.onStatusChange
	cs.mu.RUnlock()
	if callback != nil {
		go callback(cs.GetStatus())
	}
}

func (cs *CoreService) start() error {
	cs.mu.Lock()
	if cs.running {
		cs.mu.Unlock()
		if err := cs.stop(); err != nil {
			return err
		}
	} else {
		cs.mu.Unlock()
	}

	time.Sleep(300 * time.Millisecond)

	cwd, err := os.Getwd()
	if err != nil {
		cs.mu.Lock()
		cs.lastError = fmt.Sprintf("failed to get working directory: %v", err)
		cs.mu.Unlock()
		return err
	}
	workDirPath := filepath.Join(cwd, "clash")
	binPath := getBinaryPath(workDirPath)
	if binPath == "" {
		cs.mu.Lock()
		cs.lastError = fmt.Sprintf("core binary not found: %s", binPath)
		cs.mu.Unlock()
		return fmt.Errorf("core binary not found")
	}

	coreCfg := config.GetCoreConfig()
	minimalConfig := fmt.Sprintf("external-controller: %s:%d", coreCfg.APIHost, coreCfg.APIPort)
	base64Config := base64.StdEncoding.EncodeToString([]byte(minimalConfig))

	cmd := exec.Command(
		binPath,
		"-d", workDirPath,
		"-config", base64Config,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	setProcessDeathSignal(cmd)

	if err := cmd.Start(); err != nil {
		cs.mu.Lock()
		cs.lastError = fmt.Sprintf("failed to start core: %v", err)
		cs.mu.Unlock()
		return err
	}

	cs.mu.Lock()
	cs.cmd = cmd
	cs.running = true
	cs.lastError = ""
	cs.autoRestart = true
	cs.mu.Unlock()

	go func() {
		waitErr := cmd.Wait()
		cs.mu.Lock()
		wasRunning := cs.running
		cs.running = false
		if waitErr != nil {
			cs.lastError = fmt.Sprintf("core exited with error: %v", waitErr)
		}
		shouldRestart := cs.autoRestart && wasRunning
		cs.mu.Unlock()

		cs.stopWSProxies()
		cs.notifyStatusChange()

		if shouldRestart {
			time.Sleep(3 * time.Second)
			cs.mu.RLock()
			autoRestart := cs.autoRestart
			cs.mu.RUnlock()
			if autoRestart {
				cs.Start()
			}
		}
	}()

	for i := 0; i < 10; i++ {
		time.Sleep(500 * time.Millisecond)
		if version, err := cs.getVersion(); err == nil {
			cs.mu.Lock()
			cs.version = version
			cs.mu.Unlock()
			break
		}
	}

	if result := cs.applyConfigToCore(); result.Error != nil {
		log.Printf("[CoreService] Failed to apply config: %v", result.Error)
	}

	go cs.startWSProxies()

	cs.notifyStatusChange()

	return nil
}

func (cs *CoreService) stop() error {
	cs.mu.Lock()
	if !cs.running || cs.cmd == nil {
		cs.mu.Unlock()
		return nil
	} else {
		cs.mu.Unlock()
	}

	cs.mu.Lock()
	cs.autoRestart = false
	cs.mu.Unlock()

	cs.stopWSProxies()
	// 核心停止, 终止进行中的测延迟任务, 避免 goroutine 泄漏
	cs.cancelLifecycle()

	if err := cs.cmd.Process.Signal(os.Interrupt); err != nil {
		cs.cmd.Process.Kill()
	}

	done := make(chan struct{})
	go func() {
		cs.cmd.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		cs.cmd.Process.Kill()
		<-done
	}

	cs.mu.Lock()
	cs.running = false
	cs.mu.Unlock()

	return nil
}

func (cs *CoreService) getVersion() (string, error) {
	var result struct {
		Version string `json:"version"`
	}
	if err := cs.apiRequest("GET", "/version", nil, &result); err != nil {
		return "", err
	}
	return result.Version, nil
}

func (cs *CoreService) startWSProxies() {
	cs.mu.RLock()
	hub := cs.hub
	cs.mu.RUnlock()

	if hub == nil {
		return
	}

	cs.mu.Lock()
	cs.wsStopChan = make(chan struct{})
	cs.mu.Unlock()

	go cs.runWSEndpoint("traffic", ws.TypeTraffic)
	go cs.runWSEndpoint("connections", ws.TypeConnections)
	go cs.runWSEndpoint("logs", ws.TypeLogs)
	go cs.runWSEndpoint("memory", ws.TypeMemory)
}

func (cs *CoreService) stopWSProxies() {
	cs.mu.Lock()
	if cs.wsStopChan != nil {
		select {
		case <-cs.wsStopChan:
		default:
			close(cs.wsStopChan)
		}
	}
	for _, conn := range cs.wsConns {
		conn.Close()
	}
	cs.wsConns = make(map[string]*websocket.Conn)
	cs.mu.Unlock()
}

func (cs *CoreService) runWSEndpoint(endpoint string, msgType ws.MessageType) {
	for {
		cs.mu.RLock()
		stopChan := cs.wsStopChan
		cs.mu.RUnlock()

		select {
		case <-stopChan:
			return
		default:
		}

		if err := cs.connectWSEndpoint(endpoint, msgType); err != nil {
			log.Printf("[CoreService] WS %s connection error: %v", endpoint, err)
		}

		cs.mu.RLock()
		stopChan = cs.wsStopChan
		cs.mu.RUnlock()

		select {
		case <-stopChan:
			return
		case <-time.After(3 * time.Second):
		}
	}
}

func (cs *CoreService) connectWSEndpoint(endpoint string, msgType ws.MessageType) error {
	coreCfg := config.GetCoreConfig()
	wsURL := fmt.Sprintf("ws://%s:%d/%s", coreCfg.APIHost, coreCfg.APIPort, endpoint)

	headers := http.Header{}
	if coreCfg.APISecret != "" {
		headers.Set("Authorization", "Bearer "+coreCfg.APISecret)
	}

	conn, _, err := cs.wsDialer.Dial(wsURL, headers)
	if err != nil {
		return err
	}

	cs.mu.Lock()
	cs.wsConns[endpoint] = conn
	cs.mu.Unlock()

	defer func() {
		cs.mu.Lock()
		delete(cs.wsConns, endpoint)
		cs.mu.Unlock()
		conn.Close()
	}()

	cs.mu.RLock()
	stopChan := cs.wsStopChan
	hub := cs.hub
	cs.mu.RUnlock()

	for {
		select {
		case <-stopChan:
			return nil
		default:
		}

		_, message, err := conn.ReadMessage()
		if err != nil {
			return err
		}

		var rawMsg json.RawMessage
		if err := json.Unmarshal(message, &rawMsg); err != nil {
			continue
		}

		if hub != nil {
			hub.BroadcastToTypeRaw(msgType, rawMsg)
		}
	}
}

func (cs *CoreService) needsRestart(oldConfig, newConfig map[string]interface{}) bool {
	for _, key := range RestartRequiredKeys {
		if !reflect.DeepEqual(oldConfig[key], newConfig[key]) {
			return true
		}
	}
	return false
}

// mergeConfig 合并当前订阅与自定义配置, 返回合并结果与其 YAML 表示。
//
// 合并失败时返回错误, 不静默回退最小配置。
func (cs *CoreService) mergeConfig() (map[string]interface{}, string, error) {
	merger := NewMergerService()
	mergedConfig, err := merger.Merge()
	if err != nil {
		// 合并失败时不得以最小配置或其他替代配置静默下发
		return nil, "", fmt.Errorf("合并配置失败: %w", err)
	}

	yamlConfig, err := merger.GenerateYAML(mergedConfig)
	if err != nil {
		return nil, "", fmt.Errorf("生成配置失败: %w", err)
	}
	return mergedConfig, yamlConfig, nil
}

// saveLastConfig 记录当前生效的配置, 用于下次变更比对。
func (cs *CoreService) saveLastConfig(mergedConfig map[string]interface{}) {
	cs.mu.Lock()
	cs.lastConfig = mergedConfig
	cs.mu.Unlock()
}

// ApplyConfig 合并当前配置并下发到核心。
//
// 语义约定:
//   - 合并失败: 显式返回错误, 核心继续运行原有配置, 不静默回退最小配置
//   - 需要重启的键有变化: 走软重启, 软重启失败降级为杀进程重启
//   - 其余情况: 走热重载
func (cs *CoreService) ApplyConfig() ApplyConfigResult {
	mergedConfig, yamlConfig, err := cs.mergeConfig()
	if err != nil {
		return ApplyConfigResult{Error: err}
	}

	// 先取旧配置快照完成比较, 再写入新配置。
	// 若先写入再比较, 比较的是同一份 map, needsRestart 恒为 false。
	oldConfig := cs.getLastConfig()
	needRestart := cs.needsRestart(oldConfig, mergedConfig)

	cs.saveLastConfig(mergedConfig)

	if !needRestart {
		err = cs.HotReload(yamlConfig)
		return ApplyConfigResult{HotReloaded: true, Error: err}
	}

	// 软重启路径, 失败时降级为杀进程重启
	err = cs.restartWithFallback(yamlConfig)
	return ApplyConfigResult{Restarted: true, Error: err}
}

// applyConfigToCore 合并当前配置并热重载到核心, 不做重启判定。
//
// 用于核心刚启动的场景: 此时核心是全新进程, 无需重启, 只需加载配置。
// 独立于 ApplyConfig 是为了避免与降级重启路径形成递归调用
// (restartWithFallback -> Restart -> start -> ApplyConfig -> restartWithFallback)。
func (cs *CoreService) applyConfigToCore() ApplyConfigResult {
	mergedConfig, yamlConfig, err := cs.mergeConfig()
	if err != nil {
		return ApplyConfigResult{Error: err}
	}
	cs.saveLastConfig(mergedConfig)
	err = cs.HotReload(yamlConfig)
	return ApplyConfigResult{HotReloaded: true, Error: err}
}

// restartWithFallback 先应用配置再重启核心, 软重启失败或未实际生效时降级为杀进程重启。
//
// 注意顺序: 必须先热重载再软重启。核心的 /restart 接口以内存中的当前配置
// 重启, 若不先下发新配置, 重启后仍运行旧配置。
//
// 降级重启的核心以最小配置启动, 因此重启成功后需重新应用完整配置。
// 两条路径都失败时返回非空错误, 不得让调用方误以为核心已生效。
func (cs *CoreService) restartWithFallback(yamlConfig string) error {
	if err := cs.HotReload(yamlConfig); err != nil {
		return fmt.Errorf("应用配置失败: %w", err)
	}

	err := cs.SoftRestart(yamlConfig)
	if err == nil {
		// 软重启接口返回成功不代表变更已生效:
		// 端口类配置需重新绑定套接字, 软重启不会真正生效
		if verifyErr := cs.verifyRestartApplied(); verifyErr != nil {
			log.Printf("[CoreService] Soft restart returned success but change not applied (%v), falling back to kill restart", verifyErr)
			err = verifyErr
		} else {
			return nil
		}
	}
	log.Printf("[CoreService] Soft restart failed, falling back to kill restart: %v", err)

	if restartErr := cs.Restart(); restartErr != nil {
		return fmt.Errorf("软重启失败 (%v), 降级重启也失败: %w", err, restartErr)
	}

	// 重启后核心运行的是最小配置, 需重新应用完整配置
	if reloadErr := cs.HotReload(yamlConfig); reloadErr != nil {
		return fmt.Errorf("软重启失败 (%v), 降级重启成功但重新应用配置失败: %w", err, reloadErr)
	}
	return nil
}

// restartVerifyKeys 需要重启才能生效、且可从核心读取以校验的配置项。
var restartVerifyKeys = append([]string{"allow-lan"}, RestartRequiredKeys...)

// verifyRestartApplied 校验需要重启才生效的配置项是否已在核心侧实际生效。
//
// 软重启对端口类配置只是重置而非重新绑定, 接口返回成功但变更并未生效,
// 因此必须校验实际状态, 未生效时需要降级为杀进程重启。
func (cs *CoreService) verifyRestartApplied() error {
	// 核心重启后需要一点时间才能提供查询接口
	var current map[string]interface{}
	var lastErr error
	for i := 0; i < 5; i++ {
		if err := cs.apiRequest("GET", "/configs", nil, &current); err != nil {
			lastErr = err
			time.Sleep(300 * time.Millisecond)
			continue
		}
		lastErr = nil
		break
	}
	if lastErr != nil {
		return fmt.Errorf("读取核心配置失败: %w", lastErr)
	}

	target := cs.getLastConfig()
	for _, key := range restartVerifyKeys {
		want, ok := target[key]
		if !ok || want == nil {
			// 该键本次未参与变更, 无需校验
			continue
		}
		got, has := current[key]
		if !has {
			// 核心不返回该键(如 secret), 无法确认是否生效, 保守判定为未生效
			return fmt.Errorf("%s 无法确认生效状态", key)
		}
		if !reflect.DeepEqual(normalizeConfigValue(want), normalizeConfigValue(got)) {
			return fmt.Errorf("%s 未生效: 期望 %v, 实际 %v", key, want, got)
		}
	}
	return nil
}

// normalizeConfigValue 归一化配置值的数值类型, 消除 JSON 往返产生的
// int 与 float64 差异, 使配置比较结果可靠。
func normalizeConfigValue(v interface{}) interface{} {
	switch t := v.(type) {
	case int:
		return float64(t)
	case int8:
		return float64(t)
	case int16:
		return float64(t)
	case int32:
		return float64(t)
	case int64:
		return float64(t)
	case uint:
		return float64(t)
	case uint8:
		return float64(t)
	case uint16:
		return float64(t)
	case uint32:
		return float64(t)
	case uint64:
		return float64(t)
	case float32:
		return float64(t)
	default:
		return v
	}
}

func (cs *CoreService) getLastConfig() map[string]interface{} {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	if cs.lastConfig != nil {
		return cs.lastConfig
	}
	return make(map[string]interface{})
}
