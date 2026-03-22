package handler

import (
	"net/http"

	"clash-server/internal/middleware"
	"clash-server/internal/repository"
	"clash-server/internal/service"
	"clash-server/pkg/response"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userRepo *repository.UserRepository
}

func NewUserHandler() *UserHandler {
	return &UserHandler{userRepo: repository.NewUserRepository()}
}

func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	userID := middleware.GetUserID(c)
	user, err := h.userRepo.FindByID(userID)
	if err != nil {
		c.JSON(http.StatusOK, response.NotFound("用户不存在"))
		return
	}
	c.JSON(http.StatusOK, response.Success(gin.H{"user": user}))
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

func (h *UserHandler) ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, response.BadRequest("请输入旧密码和新密码（新密码至少6位）"))
		return
	}
	userID := middleware.GetUserID(c)
	authService := service.NewAuthService()
	if err := authService.ChangePassword(userID, req.OldPassword, req.NewPassword); err != nil {
		c.JSON(http.StatusOK, response.InvalidPassword("旧密码错误，请重新输入"))
		return
	}
	c.JSON(http.StatusOK, response.Success(nil))
}
