package auth

import (
	"PInfo-server/api"
	"PInfo-server/config"
	"PInfo-server/model"
	"PInfo-server/service"
	"PInfo-server/utils"
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func loginHandler(c *gin.Context) {
	loginReq := &api.LoginReq{}
	if err := c.ShouldBind(loginReq); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "参数错误",
			"data":    nil,
		})
		return
	}

	// 验证用户名
	err, detailInfo := service.DefaultService.GetUserInfo(c, loginReq.UserName)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "参数错误",
			"data":    nil,
		})
		return
	}

	userInfo := &model.UserInfo{
		UserName: loginReq.UserName,
		Uid:      detailInfo.Uid,
	}

	_, enc := utils.EncryptPassword(loginReq.PassWord)
	log.Printf("enc passwd hash:%s\n", enc)

	// 验证密码
	if !utils.CheckPasswordHash(loginReq.PassWord, detailInfo.PassHash) {
		c.JSON(http.StatusOK, gin.H{
			"code":    5000,
			"message": "用户名或密码错误",
			"data":    nil,
		})
		return
	}
	log.Printf("passwd check pass\n")

	// 生成token
	err, token := service.DefaultService.CreateJwt(context.TODO(), userInfo)
	if err != nil {
		log.Printf("gen token failed.")
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "参数错误",
			"data":    nil,
		})
		return
	}

	rsp := api.LoginRsp{
		Type:   "Bearer",
		Token:  token,
		Expire: config.AppConfig().ServerInfo.TokenExpire,
		UserInfo: api.UserBasicInfo{
			Uid:       detailInfo.Uid,
			NickName:  detailInfo.NickName,
			Signature: detailInfo.Motto,
			Avatar:    detailInfo.Avatar,
		},
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Hello Welcome to PIM",
		"data":    rsp,
	})

	/*
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
	*/
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

	req := &api.RegisterReq{}
	if err := c.ShouldBind(req); err != nil {
		utils.SendJsonRsp(c, &api.CommRsp{
			Code:    400,
			Message: "invalid param",
			Data:    nil,
		})
		return
	}

	err, rsp := service.DefaultService.RegisterUser(c, utils.GetUid(c), req)
	if err != nil {
		utils.SendJsonRsp(c, rsp)
		return
	}
	utils.SendJsonRsp(c, rsp)
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
