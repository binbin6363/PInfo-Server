package sms

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func sendSmsCodeHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Hello Welcome to PIM",
		"data":    nil,
	})
}

// Routers .
func Routers(r *gin.Engine) {
	login := r.Group("/api/v1/common")
	{
		login.POST("/sms-code", sendSmsCodeHandler)
	}
}
