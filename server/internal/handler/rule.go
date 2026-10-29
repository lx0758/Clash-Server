package handler

import (
	"net/http"

	"clash-server/internal/service"
	"clash-server/pkg/response"

	"github.com/gin-gonic/gin"
)

type RuleHandler struct {
	coreService *service.CoreService
}

func NewRuleHandler() *RuleHandler {
	return &RuleHandler{
		coreService: service.GetCoreService(),
	}
}

func (h *RuleHandler) GetRules(c *gin.Context) {
	rules, err := h.coreService.GetRules()
	if err != nil {
		// 核心不可用时返回业务错误码, 不得以空规则列表伪装成功
		c.JSON(http.StatusOK, response.CoreError("核心未运行，无法获取规则列表: "+err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.Success(gin.H{"rules": rules}))
}
