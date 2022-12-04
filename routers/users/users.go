package users

import (
	"PInfo-server/api"
	"PInfo-server/service"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

func usersSettingHandler(c *gin.Context) {
	username, ok := c.Get("username")
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusNotFound,
			"message": "未鉴权",
			"data":    nil,
		})
		return
	}

	err, userInfo := service.DefaultService.GetUserInfo(c, username.(string))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusNotFound,
			"message": "not found user",
			"data":    nil,
		})
		return
	}

	user := api.UserDetailInfo{
		IsQiYe:   false,
		Gender:   userInfo.Gender,
		Email:    userInfo.Email,
		Avatar:   userInfo.Avatar,
		Mobile:   userInfo.Phone,
		Motto:    userInfo.Motto,
		NickName: userInfo.NickName,
		Uid:      userInfo.Uid,
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "ok",
		"data": api.UserSettingRsp{
			SettingInfo: api.SettingInfo{},
			UserInfo:    user,
		},
	})
}

func usersDetailHandler(c *gin.Context) {
	username, ok := c.Get("username")
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusNotFound,
			"message": "未鉴权",
			"data":    nil,
		})
		return
	}

	err, userInfo := service.DefaultService.GetUserInfo(c, username.(string))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusNotFound,
			"message": "not found user",
			"data":    nil,
		})
		return
	}
	user := api.UserDetailInfo{
		IsQiYe:   false,
		Gender:   userInfo.Gender,
		Email:    userInfo.Email,
		Avatar:   userInfo.Avatar,
		Mobile:   userInfo.Phone,
		Motto:    userInfo.Motto,
		NickName: userInfo.NickName,
		Uid:      userInfo.Uid,
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "ok",
		"data":    user,
	})
}

func modifyUsersSettingHandler(c *gin.Context) {
	req := &api.ModifyUsersSettingReq{}
	err := c.ShouldBind(req)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": "请求参数错误",
			"data":    nil,
		})
		return
	}
	uid, ok := c.Get("uid")
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusNotFound,
			"message": "未鉴权",
			"data":    nil,
		})
		return
	}

	err = service.DefaultService.SetUserInfo(c, cast.ToInt64(uid), req)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusNotFound,
			"message": "not found user",
			"data":    nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "ok",
		"data":    nil,
	})
}

// Routers .
func Routers(r *gin.Engine) {
	login := r.Group("/api/v1/users")
	{
		login.GET("/setting", usersSettingHandler)
		login.GET("/detail", usersDetailHandler)
		login.POST("/change/detail", modifyUsersSettingHandler)
	}
}
