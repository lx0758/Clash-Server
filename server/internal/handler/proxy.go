package handler

import (
	"net/http"
	"strconv"

	"clash-server/internal/service"
	"clash-server/pkg/response"

	"github.com/gin-gonic/gin"
)

type ProxyHandler struct {
	coreService *service.CoreService
}

func NewProxyHandler() *ProxyHandler {
	return &ProxyHandler{coreService: service.GetCoreService()}
}

func (h *ProxyHandler) List(c *gin.Context) {
	proxies, err := h.coreService.GetProxies()
	if err != nil {
		// 核心不可用时返回业务错误码, 不得以空代理列表伪装成功
		c.JSON(http.StatusOK, response.CoreError("核心未运行，无法获取代理列表: "+err.Error()))
		return
	}
	if proxyList, ok := proxies["proxies"].(map[string]interface{}); ok {
		c.JSON(http.StatusOK, response.Success(gin.H{"proxies": proxyList}))
		return
	}
	c.JSON(http.StatusOK, response.Success(gin.H{"proxies": map[string]interface{}{}}))
}

func (h *ProxyHandler) Get(c *gin.Context) {
	name := c.Param("name")
	proxies, err := h.coreService.GetProxies()
	if err != nil {
		c.JSON(http.StatusOK, response.CoreError("核心未运行，无法获取代理信息: "+err.Error()))
		return
	}
	proxyList, ok := proxies["proxies"].(map[string]interface{})
	if !ok {
		c.JSON(http.StatusOK, response.NotFound("代理不存在"))
		return
	}
	if proxy, ok := proxyList[name]; ok {
		c.JSON(http.StatusOK, response.Success(gin.H{"proxy": proxy}))
		return
	}
	c.JSON(http.StatusOK, response.NotFound("代理不存在"))
}

type SelectProxyRequest struct {
	Name string `json:"name" binding:"required"`
}

func (h *ProxyHandler) Select(c *gin.Context) {
	group := c.Param("name")
	var req SelectProxyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, response.BadRequest("请选择要切换的代理节点"))
		return
	}
	if err := h.coreService.SelectProxy(group, req.Name); err != nil {
		c.JSON(http.StatusOK, response.CoreError("切换代理节点失败，核心未运行或代理组不存在"))
		return
	}
	c.JSON(http.StatusOK, response.Success(nil))
}

func (h *ProxyHandler) CheckDelay(c *gin.Context) {
	name := c.Param("name")
	testURL := c.DefaultQuery("url", "https://www.gstatic.com/generate_204")
	timeout, _ := strconv.Atoi(c.DefaultQuery("timeout", "5000"))

	delay, err := h.coreService.CheckDelay(name, testURL, timeout)
	if err != nil {
		// 延迟测试失败时返回原因, 不得以 delay:0 伪装成功
		c.JSON(http.StatusOK, response.CoreError("测试延迟失败: "+err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.Success(gin.H{"delay": delay}))
}

func (h *ProxyHandler) CheckGroupDelay(c *gin.Context) {
	groupName := c.Param("group")
	testURL := c.DefaultQuery("url", "https://www.gstatic.com/generate_204")
	timeout, _ := strconv.Atoi(c.DefaultQuery("timeout", "5000"))

	if err := h.coreService.CheckGroupDelay(groupName, testURL, timeout); err != nil {
		c.JSON(http.StatusOK, response.CoreError("测试延迟失败: "+err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.Success(nil))
}

// StartGroupDelayTest 启动逐节点测延迟任务, 立即返回。
// 各节点的延迟结果通过 WebSocket 的 delay 类型消息陆续推送,
// 全部完成后推送该组的 done 信号。
func (h *ProxyHandler) StartGroupDelayTest(c *gin.Context) {
	groupName := c.Param("group")
	testURL := c.DefaultQuery("url", "https://www.gstatic.com/generate_204")
	timeout, _ := strconv.Atoi(c.DefaultQuery("timeout", "5000"))
	concurrency, _ := strconv.Atoi(c.DefaultQuery("concurrency", "0"))

	names, err := h.coreService.GetGroupProxyNames(groupName)
	if err != nil {
		c.JSON(http.StatusOK, response.CoreError("无法获取代理组节点: "+err.Error()))
		return
	}
	if len(names) == 0 {
		c.JSON(http.StatusOK, response.BadRequest("代理组内没有可测试的节点"))
		return
	}

	// 测速任务的生命周期跟随核心, 而非跟随本请求:
	// 请求返回后任务仍需继续, 核心停止时任务被取消
	req := service.DelayTestRequest{
		Group:       groupName,
		TestURL:     testURL,
		Timeout:     timeout,
		Concurrency: concurrency,
	}
	if err := h.coreService.StartDelayTest(req); err != nil {
		c.JSON(http.StatusOK, response.CoreError("启动测延迟失败: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"group": groupName,
		"nodes": names,
		"total": len(names),
	}))
}

func (h *ProxyHandler) GetMode(c *gin.Context) {
	mode, err := h.coreService.GetMode()
	if err != nil {
		// 核心不可用时返回业务错误码, 不得以默认模式伪装成功
		c.JSON(http.StatusOK, response.CoreError("核心未运行，无法获取运行模式: "+err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.Success(gin.H{"mode": mode}))
}

type SetModeRequest struct {
	Mode string `json:"mode" binding:"required"`
}

func (h *ProxyHandler) SetMode(c *gin.Context) {
	var req SetModeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, response.BadRequest("请指定运行模式"))
		return
	}
	if req.Mode != "rule" && req.Mode != "global" && req.Mode != "direct" {
		c.JSON(http.StatusOK, response.BadRequest("无效的运行模式"))
		return
	}
	if err := h.coreService.SetMode(req.Mode); err != nil {
		c.JSON(http.StatusOK, response.CoreError("切换模式失败: "+err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.Success(gin.H{"mode": req.Mode}))
}
