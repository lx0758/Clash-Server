package handler

import (
	"net/http"

	"clash-server/pkg/response"
	"clash-server/res"

	"github.com/gin-gonic/gin"
)

func GetVersion(c *gin.Context) {
	c.JSON(http.StatusOK, response.Success(gin.H{
		"version": res.Version,
	}))
}
