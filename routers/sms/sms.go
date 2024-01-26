package sms

import (
	"PInfo-server/routers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func init() {
	routers.Register(Routers)
}

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
