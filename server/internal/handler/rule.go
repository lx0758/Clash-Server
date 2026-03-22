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
		c.JSON(http.StatusOK, response.Success(gin.H{"rules": []interface{}{}}))
		return
	}
	c.JSON(http.StatusOK, response.Success(gin.H{"rules": rules}))
}
