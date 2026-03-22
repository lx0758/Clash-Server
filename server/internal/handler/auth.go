package handler

import (
	"net/http"

	"clash-server/internal/middleware"
	"clash-server/internal/service"
	"clash-server/pkg/response"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{authService: service.NewAuthService()}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type InitPasswordRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, response.BadRequest("请输入用户名和密码"))
		return
	}
	user, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusOK, response.Unauthorized("用户名或密码错误"))
		return
	}
	middleware.LoginUser(c, user.ID)
	c.JSON(http.StatusOK, response.Success(gin.H{"user": user}))
}

func (h *AuthHandler) Logout(c *gin.Context) {
	middleware.LogoutUser(c)
	c.JSON(http.StatusOK, response.Success(nil))
}

func (h *AuthHandler) CheckInit(c *gin.Context) {
	initialized, err := h.authService.IsInitialized()
	if err != nil {
		c.JSON(http.StatusOK, response.InternalError("检查初始化状态失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(gin.H{"initialized": initialized}))
}

func (h *AuthHandler) InitPassword(c *gin.Context) {
	var req InitPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, response.BadRequest("请输入有效的用户名和密码（密码至少6位）"))
		return
	}
	if err := h.authService.InitPassword(req.Username, req.Password); err != nil {
		c.JSON(http.StatusOK, response.BadRequest("初始化密码失败，请稍后重试"))
		return
	}
	c.JSON(http.StatusOK, response.Success(nil))
}
