package service

import (
	"bytes"
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
		return map[string]interface{}{"proxies": map[string]interface{}{}}, nil
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
		return []interface{}{}, nil
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
		return []interface{}{}, nil
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
		return "rule", nil
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
		return fmt.Errorf("API request failed")
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

	if result := cs.ApplyConfig(); result.Error != nil {
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

func (cs *CoreService) ApplyConfig() ApplyConfigResult {
	merger := NewMergerService()
	mergedConfig, err := merger.Merge()
	if err != nil {
		mergedConfig = merger.GetMinimalConfig()
	}

	yamlConfig, err := merger.GenerateYAML(mergedConfig)
	if err != nil {
		return ApplyConfigResult{Error: err}
	}
	cs.mu.Lock()
	cs.lastConfig = mergedConfig
	cs.mu.Unlock()

	if !cs.needsRestart(cs.getLastConfig(), mergedConfig) {
		err = cs.HotReload(yamlConfig)
		return ApplyConfigResult{HotReloaded: true, Error: err}
	} else {
		err = cs.SoftRestart(yamlConfig)
		return ApplyConfigResult{HotReloaded: true, Error: err}
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
