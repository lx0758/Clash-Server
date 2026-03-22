package handler

import (
	"net/http"
	"strconv"

	"clash-server/internal/config"
	"clash-server/internal/repository"
	"clash-server/internal/service"
	"clash-server/pkg/response"

	"github.com/gin-gonic/gin"
)

type ConfigHandler struct {
	coreService *service.CoreService
	coreRepo    *repository.ConfigRepository
}

func NewConfigHandler() *ConfigHandler {
	return &ConfigHandler{
		coreService: service.GetCoreService(),
		coreRepo:    repository.NewConfigRepository(),
	}
}

func (h *ConfigHandler) Get(c *gin.Context) {
	serverCfg := config.GetServerConfig()
	coreCfg := config.GetCoreConfig()
	c.JSON(http.StatusOK, response.Success(gin.H{
		"server": gin.H{
			"host":     serverCfg.Host,
			"port":     serverCfg.Port,
			"database": serverCfg.Database,
		},
		"core": gin.H{
			"mode":       coreCfg.Mode,
			"api_host":   coreCfg.APIHost,
			"api_port":   coreCfg.APIPort,
			"api_secret": coreCfg.APISecret,
			"mixed_port": coreCfg.MixedPort,
			"allow_lan":  coreCfg.AllowLan,
			"log_level":  coreCfg.LogLevel,
			"ipv6":       coreCfg.IPv6,
		},
	}))
}

type ConfigUpdateRequest struct {
	Core *CoreConfigUpdate `json:"core"`
}

type CoreConfigUpdate struct {
	APIHost   *string `json:"api_host"`
	APIPort   *int    `json:"api_port"`
	APISecret *string `json:"api_secret"`
	MixedPort *int    `json:"mixed_port"`
	AllowLan  *bool   `json:"allow_lan"`
	Mode      *string `json:"mode"`
	LogLevel  *string `json:"log_level"`
	IPv6      *bool   `json:"ipv6"`
}

func (h *ConfigHandler) Update(c *gin.Context) {
	var req ConfigUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, response.BadRequest("配置格式错误，请检查输入"))
		return
	}
	if req.Core == nil {
		c.JSON(http.StatusOK, response.BadRequest("缺少 core 配置"))
		return
	}
	kv := make(map[string]string)
	if req.Core.APIHost != nil {
		kv["core:api_host"] = *req.Core.APIHost
	}
	if req.Core.APIPort != nil {
		kv["core:api_port"] = strconv.Itoa(*req.Core.APIPort)
	}
	if req.Core.APISecret != nil {
		kv["core:api_secret"] = *req.Core.APISecret
	}
	if req.Core.MixedPort != nil {
		kv["core:mixed_port"] = strconv.Itoa(*req.Core.MixedPort)
	}
	if req.Core.AllowLan != nil {
		kv["core:allow_lan"] = strconv.FormatBool(*req.Core.AllowLan)
	}
	if req.Core.Mode != nil {
		kv["core:mode"] = *req.Core.Mode
	}
	if req.Core.LogLevel != nil {
		kv["core:log_level"] = *req.Core.LogLevel
	}
	if req.Core.IPv6 != nil {
		kv["core:ipv6"] = strconv.FormatBool(*req.Core.IPv6)
	}
	if len(kv) > 0 {
		if err := h.coreRepo.SetMulti(kv); err != nil {
			c.JSON(http.StatusOK, response.InternalError("保存配置失败，请稍后重试"))
			return
		}
	}
	coreCfg := config.GetCoreConfig()
	coreErr := h.applyConfig()
	c.JSON(http.StatusOK, response.SuccessWithCoreError(gin.H{
		"core": gin.H{
			"mode":       coreCfg.Mode,
			"api_host":   coreCfg.APIHost,
			"api_port":   coreCfg.APIPort,
			"api_secret": coreCfg.APISecret,
			"mixed_port": coreCfg.MixedPort,
			"allow_lan":  coreCfg.AllowLan,
			"log_level":  coreCfg.LogLevel,
			"ipv6":       coreCfg.IPv6,
		},
	}, coreErr))
}

func (h *ConfigHandler) applyConfig() string {
	result := h.coreService.ApplyConfig()
	if result.Error != nil {
		return result.Error.Error()
	}
	return ""
}
