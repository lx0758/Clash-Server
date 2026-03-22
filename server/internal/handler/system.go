package handler

import (
	"net/http"

	"clash-server/internal/service"
	"clash-server/pkg/response"

	"github.com/gin-gonic/gin"
)

type SystemHandler struct {
	coreService *service.CoreService
	subService  *service.SubscriptionService
}

func NewSystemHandler() *SystemHandler {
	return &SystemHandler{
		coreService: service.GetCoreService(),
		subService:  service.NewSubscriptionService(),
	}
}

func (h *SystemHandler) GetInfo(c *gin.Context) {
	coreStatus := h.coreService.GetStatus()

	subs, err := h.subService.List()
	subCount := 0
	proxyCount := 0
	if err == nil {
		subCount = len(subs)
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"core": gin.H{
			"running": coreStatus.Running,
			"version": coreStatus.Version,
			"error":   coreStatus.Error,
		},
		"subscription": gin.H{
			"count":       subCount,
			"proxy_count": proxyCount,
		},
		"traffic": nil,
	}))
}
