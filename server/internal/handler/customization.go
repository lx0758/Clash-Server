package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"clash-server/internal/model"
	"clash-server/internal/service"
	"clash-server/pkg/response"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

type CustomizationHandler struct {
	service     *service.CustomizationService
	subService  *service.SubscriptionService
	coreService *service.CoreService
	merger      *service.MergerService
}

func NewCustomizationHandler() *CustomizationHandler {
	return &CustomizationHandler{
		service:     service.NewCustomizationService(),
		subService:  service.NewSubscriptionService(),
		coreService: service.GetCoreService(),
		merger:      service.GetMergerService(),
	}
}

func (h *CustomizationHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, response.BadRequest("无效的订阅ID"))
		return
	}

	if _, err := h.subService.Get(uint(id)); err != nil {
		c.JSON(http.StatusOK, response.NotFound("订阅不存在"))
		return
	}

	customization, err := h.service.GetBySubscriptionID(uint(id))
	if err != nil {
		c.JSON(http.StatusOK, response.Success(gin.H{"customization": nil}))
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{"customization": customization}))
}

type UpdateCustomizationRequest struct {
	ProxyInsert      string `json:"proxy_insert"`
	ProxyAppend      string `json:"proxy_append"`
	ProxyRemove      string `json:"proxy_remove"`
	ProxyGroupInsert string `json:"proxy_group_insert"`
	ProxyGroupAppend string `json:"proxy_group_append"`
	ProxyGroupRemove string `json:"proxy_group_remove"`
	RuleInsert       string `json:"rule_insert"`
	RuleAppend       string `json:"rule_append"`
	RuleRemove       string `json:"rule_remove"`
	GlobalOverride   string `json:"global_override"`
	Script           string `json:"script"`
}

func (h *CustomizationHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, response.BadRequest("无效的订阅ID"))
		return
	}

	if _, err := h.subService.Get(uint(id)); err != nil {
		c.JSON(http.StatusOK, response.NotFound("订阅不存在"))
		return
	}

	var req UpdateCustomizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, response.BadRequest("无效的请求数据"))
		return
	}

	if err := validateYAMLFields(req); err != nil {
		c.JSON(http.StatusOK, response.BadRequest(err.Error()))
		return
	}

	customization := &model.SubscriptionCustomization{
		SubscriptionID:   uint(id),
		ProxyInsert:      req.ProxyInsert,
		ProxyAppend:      req.ProxyAppend,
		ProxyRemove:      req.ProxyRemove,
		ProxyGroupInsert: req.ProxyGroupInsert,
		ProxyGroupAppend: req.ProxyGroupAppend,
		ProxyGroupRemove: req.ProxyGroupRemove,
		RuleInsert:       req.RuleInsert,
		RuleAppend:       req.RuleAppend,
		RuleRemove:       req.RuleRemove,
		GlobalOverride:   req.GlobalOverride,
		Script:           req.Script,
	}

	if err := h.service.Save(customization); err != nil {
		c.JSON(http.StatusOK, response.InternalError("保存自定义配置失败"))
		return
	}

	mergedConfig, err := h.merger.MergeForSubscription(uint(id))
	if err != nil {
		c.JSON(http.StatusOK, response.BadRequest("合并配置失败: "+err.Error()))
		return
	}

	if err := h.merger.Validate(mergedConfig); err != nil {
		c.JSON(http.StatusOK, response.BadRequest("配置验证失败: "+err.Error()))
		return
	}

	result := h.coreService.ApplyConfig()
	if result.Error != nil {
		c.JSON(http.StatusOK, response.Error(response.CodeConfigError, "核心应用失败: "+result.Error.Error()))
		return
	}
	c.JSON(http.StatusOK, response.Success(gin.H{"customization": customization}))
}

func restoreCustomization(svc *service.CustomizationService, subID uint, old *model.SubscriptionCustomization) {
	if old == nil {
		svc.DeleteBySubscriptionID(subID)
		return
	}
	svc.Save(old)
}

func validateYAMLFields(req UpdateCustomizationRequest) error {
	yamlFields := []struct {
		value string
		name  string
	}{
		{req.ProxyInsert, "插入节点"},
		{req.ProxyAppend, "追加节点"},
		{req.ProxyRemove, "移除节点"},
		{req.ProxyGroupInsert, "插入代理组"},
		{req.ProxyGroupAppend, "追加代理组"},
		{req.ProxyGroupRemove, "移除代理组"},
		{req.RuleInsert, "插入规则"},
		{req.RuleAppend, "追加规则"},
		{req.RuleRemove, "移除规则"},
		{req.GlobalOverride, "全局配置"},
	}

	for _, field := range yamlFields {
		if strings.TrimSpace(field.value) == "" {
			continue
		}
		var v interface{}
		if err := yaml.Unmarshal([]byte(field.value), &v); err != nil {
			return fmt.Errorf("%s YAML 格式错误: %w", field.name, err)
		}
	}

	return nil
}
