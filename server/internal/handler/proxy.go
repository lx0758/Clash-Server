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
	if err != nil || proxies == nil {
		c.JSON(http.StatusOK, response.Success(gin.H{"proxies": map[string]interface{}{}}))
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
	if err != nil || proxies == nil {
		c.JSON(http.StatusOK, response.Success(gin.H{"proxy": nil}))
		return
	}
	proxyList, ok := proxies["proxies"].(map[string]interface{})
	if !ok {
		c.JSON(http.StatusOK, response.Success(gin.H{"proxy": nil}))
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
		c.JSON(http.StatusOK, response.Success(gin.H{"delay": 0}))
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

func (h *ProxyHandler) GetMode(c *gin.Context) {
	mode, err := h.coreService.GetMode()
	if err != nil {
		c.JSON(http.StatusOK, response.Success(gin.H{"mode": "rule"}))
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
