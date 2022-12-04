package auth

import (
	"PInfo-server/api"
	"PInfo-server/config"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func loginHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Hello Welcome to PIM",
		"data": api.LoginRsp{
			Type:   "Bearer",
			Token:  "1nlkasdhasdcasvderveq",
			Expire: 7200,
			UserInfo: api.UserBasicInfo{
				Uid:       212332324324,
				NickName:  "polite",
				Signature: "我的签名",
				Avatar:    "www.baidu.com",
			},
		},
	})

	// 通知websocket
	req := map[string]interface{}{
		"platform": "102",
		"uid":      20221113,
	}
	bytesData, _ := json.Marshal(req)
	url := fmt.Sprintf("http://%s/notice/auth/login", config.AppConfig().ConnInfo.Addr)
	_, err := http.Post(url, "application/json; charset=utf-8", bytes.NewReader(bytesData))
	if err != nil {
		log.Printf("post conn failed, err:%+v\n", err)
	} else {
		log.Printf("post conn success, req:%+v\n", req)
	}
}

func logoutHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Bye, PIM",
		"data":    nil,
	})
}

// registerHandler 注册服务接口
func registerHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Bye, PIM",
		"data":    nil,
	})
}

// refreshTokenHandler 刷新登录Token服务接口
func refreshTokenHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Bye, PIM",
		"data":    nil,
	})
}

// forgetHandler 找回密码服务
func forgetHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Bye, PIM",
		"data":    nil,
	})
}

// Routers .
func Routers(r *gin.Engine) {
	login := r.Group("/api/v1/auth")
	{
		login.POST("/login", loginHandler)
		login.POST("/logout", logoutHandler)
		login.POST("/register", registerHandler)
		login.POST("/refresh-token", refreshTokenHandler)
		login.POST("/forget", forgetHandler)
	}
}
