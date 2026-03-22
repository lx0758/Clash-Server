package handler

import (
	"net/http"
	"strconv"

	"clash-server/internal/model"
	"clash-server/internal/service"
	"clash-server/pkg/response"

	"github.com/gin-gonic/gin"
)

type SubRuleHandler struct {
	ruleService *service.RuleService
	coreService *service.CoreService
}

func NewSubRuleHandler() *SubRuleHandler {
	return &SubRuleHandler{
		ruleService: service.NewRuleService(),
		coreService: service.GetCoreService(),
	}
}

func (h *SubRuleHandler) List(c *gin.Context) {
	subscriptionID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, response.BadRequest("无效的订阅ID"))
		return
	}
	rules, err := h.ruleService.List(uint(subscriptionID))
	if err != nil {
		c.JSON(http.StatusOK, response.InternalError("获取规则列表失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(gin.H{"rules": rules}))
}

func (h *SubRuleHandler) Get(c *gin.Context) {
	subscriptionID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, response.BadRequest("无效的订阅ID"))
		return
	}
	ruleID, err := strconv.ParseUint(c.Param("rule_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, response.BadRequest("无效的规则ID"))
		return
	}
	rule, err := h.ruleService.GetBySubscription(uint(ruleID), uint(subscriptionID))
	if err != nil {
		c.JSON(http.StatusOK, response.NotFound("规则不存在"))
		return
	}
	c.JSON(http.StatusOK, response.Success(gin.H{"rule": rule}))
}

type CreateRuleRequest struct {
	Name     string `json:"name" binding:"required"`
	Type     string `json:"type" binding:"required"`
	Payload  string `json:"payload" binding:"required"`
	Proxy    string `json:"proxy" binding:"required"`
	Enabled  bool   `json:"enabled"`
	Mode     string `json:"mode"`
	Priority int    `json:"priority"`
}

func (h *SubRuleHandler) Create(c *gin.Context) {
	subscriptionID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, response.BadRequest("无效的订阅ID"))
		return
	}
	var req CreateRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, response.BadRequest("请填写完整的规则信息"))
		return
	}
	rule := &model.Rule{
		SubscriptionID: uint(subscriptionID),
		Name:           req.Name,
		Type:           req.Type,
		Payload:        req.Payload,
		Proxy:          req.Proxy,
		Enabled:        req.Enabled,
		Mode:           model.RuleMode(req.Mode),
		Priority:       req.Priority,
	}
	if rule.Mode == "" {
		rule.Mode = model.RuleModeAppend
	}
	if err := h.ruleService.Create(rule); err != nil {
		c.JSON(http.StatusOK, response.BadRequest("创建规则失败，请检查规则格式是否正确"))
		return
	}
	coreErr := h.applyConfig()
	c.JSON(http.StatusOK, response.SuccessWithCoreError(gin.H{"rule": rule}, coreErr))
}

func (h *SubRuleHandler) Update(c *gin.Context) {
	subscriptionID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, response.BadRequest("无效的订阅ID"))
		return
	}
	ruleID, err := strconv.ParseUint(c.Param("rule_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, response.BadRequest("无效的规则ID"))
		return
	}
	var req CreateRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, response.BadRequest("请填写完整的规则信息"))
		return
	}
	rule, err := h.ruleService.GetBySubscription(uint(ruleID), uint(subscriptionID))
	if err != nil {
		c.JSON(http.StatusOK, response.NotFound("规则不存在"))
		return
	}
	rule.Name = req.Name
	rule.Type = req.Type
	rule.Payload = req.Payload
	rule.Proxy = req.Proxy
	rule.Enabled = req.Enabled
	rule.Mode = model.RuleMode(req.Mode)
	rule.Priority = req.Priority
	if err := h.ruleService.Update(rule); err != nil {
		c.JSON(http.StatusOK, response.BadRequest("更新规则失败，请检查规则格式是否正确"))
		return
	}
	coreErr := h.applyConfig()
	c.JSON(http.StatusOK, response.SuccessWithCoreError(gin.H{"rule": rule}, coreErr))
}

func (h *SubRuleHandler) Delete(c *gin.Context) {
	subscriptionID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, response.BadRequest("无效的订阅ID"))
		return
	}
	ruleID, err := strconv.ParseUint(c.Param("rule_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, response.BadRequest("无效的规则ID"))
		return
	}
	if _, err := h.ruleService.GetBySubscription(uint(ruleID), uint(subscriptionID)); err != nil {
		c.JSON(http.StatusOK, response.NotFound("规则不存在"))
		return
	}
	if err := h.ruleService.Delete(uint(ruleID)); err != nil {
		c.JSON(http.StatusOK, response.InternalError("删除规则失败，请稍后重试"))
		return
	}
	coreErr := h.applyConfig()
	c.JSON(http.StatusOK, response.SuccessWithCoreError(nil, coreErr))
}

func (h *SubRuleHandler) applyConfig() string {
	result := h.coreService.ApplyConfig()
	if result.Error != nil {
		return result.Error.Error()
	}
	return ""
}
