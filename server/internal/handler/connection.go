package handler

import (
	"net/http"

	"clash-server/internal/service"
	"clash-server/pkg/response"

	"github.com/gin-gonic/gin"
)

type ConnectionHandler struct {
	coreService *service.CoreService
}

func NewConnectionHandler() *ConnectionHandler {
	return &ConnectionHandler{coreService: service.GetCoreService()}
}

func (h *ConnectionHandler) List(c *gin.Context) {
	connections, err := h.coreService.GetConnections()
	if err != nil {
		// 核心不可用时返回业务错误码, 不得以空连接列表伪装成功
		c.JSON(http.StatusOK, response.CoreError("核心未运行，无法获取连接列表: "+err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.Success(gin.H{"connections": connections}))
}

func (h *ConnectionHandler) Close(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		if err := h.coreService.CloseAllConnections(); err != nil {
			c.JSON(http.StatusOK, response.CoreError("关闭所有连接失败"))
			return
		}
	} else {
		if err := h.coreService.CloseConnection(id); err != nil {
			c.JSON(http.StatusOK, response.CoreError("关闭连接失败，连接可能已断开"))
			return
		}
	}
	c.JSON(http.StatusOK, response.Success(nil))
}

func (h *ConnectionHandler) CloseAll(c *gin.Context) {
	if err := h.coreService.CloseAllConnections(); err != nil {
		c.JSON(http.StatusOK, response.CoreError("关闭所有连接失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(nil))
}
