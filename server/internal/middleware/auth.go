package middleware

import (
	"net/http"

	"clash-server/pkg/response"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const (
	SESSION_KEY_USER_ID = "userId"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userId := session.Get(SESSION_KEY_USER_ID)
		if userId == nil {
			c.JSON(http.StatusOK, response.Unauthorized("请先登录"))
			c.Abort()
			return
		}
		c.Next()
	}
}

func LoginUser(c *gin.Context, userId uint) {
	session := sessions.Default(c)
	session.Set(SESSION_KEY_USER_ID, userId)
	_ = session.Save()
}

func LogoutUser(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete(SESSION_KEY_USER_ID)
	_ = session.Save()
}

func GetUserID(c *gin.Context) uint {
	session := sessions.Default(c)
	userId := session.Get(SESSION_KEY_USER_ID)
	userIdValue, ok := userId.(uint)
	if !ok {
		return 0
	}
	return userIdValue
}
