package login

import (
	"PInfo-server/api"
	"github.com/gin-gonic/gin"
	"net/http"
)

func loginHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Hello Welcome to PIM",
		"data": api.LoginRsp{
			Token:  "1nlkasdhasdcasvderveq",
			Expire: 7200,
			UserInfo: api.UserInfo{
				Uid:       212332324324,
				NickName:  "polite",
				Signature: "我的签名",
				Avatar:    "www.baidu.com",
			},
		},
	})
}

// Routers .
func Routers(r *gin.Engine) {
	login := r.Group("/api/v1/auth")
	{
		login.POST("/login", loginHandler)
	}
}
