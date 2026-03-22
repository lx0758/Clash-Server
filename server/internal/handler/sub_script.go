package handler

import (
	"net/http"
	"strconv"

	"clash-server/internal/model"
	"clash-server/internal/service"
	"clash-server/pkg/response"

	"github.com/gin-gonic/gin"
)

type SubScriptHandler struct {
	scriptService *service.ScriptService
	coreService   *service.CoreService
}

func NewSubScriptHandler() *SubScriptHandler {
	return &SubScriptHandler{
		scriptService: service.NewScriptService(),
		coreService:   service.GetCoreService(),
	}
}

func (h *SubScriptHandler) List(c *gin.Context) {
	subscriptionID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, response.BadRequest("无效的订阅ID"))
		return
	}
	scripts, err := h.scriptService.List(uint(subscriptionID))
	if err != nil {
		c.JSON(http.StatusOK, response.InternalError("获取脚本列表失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(gin.H{"scripts": scripts}))
}

func (h *SubScriptHandler) Get(c *gin.Context) {
	subscriptionID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, response.BadRequest("无效的订阅ID"))
		return
	}
	scriptID, err := strconv.ParseUint(c.Param("script_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, response.BadRequest("无效的脚本ID"))
		return
	}
	script, err := h.scriptService.GetBySubscription(uint(scriptID), uint(subscriptionID))
	if err != nil {
		c.JSON(http.StatusOK, response.NotFound("脚本不存在"))
		return
	}
	c.JSON(http.StatusOK, response.Success(gin.H{"script": script}))
}

type CreateScriptRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Content     string `json:"content" binding:"required"`
	Enabled     bool   `json:"enabled"`
}

func (h *SubScriptHandler) Create(c *gin.Context) {
	subscriptionID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, response.BadRequest("无效的订阅ID"))
		return
	}
	var req CreateScriptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, response.BadRequest("请输入脚本名称和内容"))
		return
	}
	script := &model.Script{
		SubscriptionID: uint(subscriptionID),
		Name:           req.Name,
		Description:    req.Description,
		Content:        req.Content,
		Enabled:        req.Enabled,
	}
	if err := h.scriptService.Create(script); err != nil {
		c.JSON(http.StatusOK, response.InternalError("创建脚本失败，请稍后重试"))
		return
	}
	coreErr := h.applyConfig()
	c.JSON(http.StatusOK, response.SuccessWithCoreError(gin.H{"script": script}, coreErr))
}

func (h *SubScriptHandler) Update(c *gin.Context) {
	subscriptionID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, response.BadRequest("无效的订阅ID"))
		return
	}
	scriptID, err := strconv.ParseUint(c.Param("script_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, response.BadRequest("无效的脚本ID"))
		return
	}
	var req CreateScriptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, response.BadRequest("请输入脚本名称和内容"))
		return
	}
	script, err := h.scriptService.GetBySubscription(uint(scriptID), uint(subscriptionID))
	if err != nil {
		c.JSON(http.StatusOK, response.NotFound("脚本不存在"))
		return
	}
	script.Name = req.Name
	script.Description = req.Description
	script.Content = req.Content
	script.Enabled = req.Enabled
	if err := h.scriptService.Update(script); err != nil {
		c.JSON(http.StatusOK, response.InternalError("更新脚本失败，请稍后重试"))
		return
	}
	coreErr := h.applyConfig()
	c.JSON(http.StatusOK, response.SuccessWithCoreError(gin.H{"script": script}, coreErr))
}

func (h *SubScriptHandler) Delete(c *gin.Context) {
	subscriptionID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, response.BadRequest("无效的订阅ID"))
		return
	}
	scriptID, err := strconv.ParseUint(c.Param("script_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, response.BadRequest("无效的脚本ID"))
		return
	}
	if _, err := h.scriptService.GetBySubscription(uint(scriptID), uint(subscriptionID)); err != nil {
		c.JSON(http.StatusOK, response.NotFound("脚本不存在"))
		return
	}
	if err := h.scriptService.Delete(uint(scriptID)); err != nil {
		c.JSON(http.StatusOK, response.InternalError("删除脚本失败，请稍后重试"))
		return
	}
	coreErr := h.applyConfig()
	c.JSON(http.StatusOK, response.SuccessWithCoreError(nil, coreErr))
}

type TestScriptRequest struct {
	Config map[string]interface{} `json:"config"`
}

func (h *SubScriptHandler) Test(c *gin.Context) {
	subscriptionID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, response.BadRequest("无效的订阅ID"))
		return
	}
	scriptID, err := strconv.ParseUint(c.Param("script_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, response.BadRequest("无效的脚本ID"))
		return
	}
	var req TestScriptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, response.BadRequest("请提供测试配置数据"))
		return
	}
	result, err := h.scriptService.TestBySubscription(uint(scriptID), uint(subscriptionID), req.Config)
	if err != nil {
		c.JSON(http.StatusOK, response.BadRequest("脚本执行失败："+err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.Success(gin.H{"result": result}))
}

func (h *SubScriptHandler) applyConfig() string {
	result := h.coreService.ApplyConfig()
	if result.Error != nil {
		return result.Error.Error()
	}
	return ""
}
