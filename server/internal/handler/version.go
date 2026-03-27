package handler

import (
	"net/http"

	"clash-server/pkg/response"
	"clash-server/pkg/version"

	"github.com/gin-gonic/gin"
)

func GetVersion(c *gin.Context) {
	c.JSON(http.StatusOK, response.Success(gin.H{
		"version": version.Get(),
	}))
}
