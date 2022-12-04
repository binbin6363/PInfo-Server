package users

import (
	"PInfo-server/api"
	"github.com/gin-gonic/gin"
	"net/http"
)

func usersSettingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Hello Welcome to PIM",
		"data": api.UserSettingRsp{
			SettingInfo: api.SettingInfo{
				KeyboardEventNotify: "",
				NotifyCueTone:       "",
				ThemeBagImg:         "",
				ThemeColor:          "",
				ThemeMode:           "",
			},
			UserInfo: api.UserDetailInfo{
				IsQiYe:   false,
				Gender:   1,
				Email:    "12123232@qq.com",
				Avatar:   "http://gimg2.baidu.com/image_search/src=http%3A%2F%2Fc-ssl.duitang.com%2Fuploads%2Fblog%2F202105%2F19%2F20210519212438_ced7e.jpg&refer=http%3A%2F%2Fc-ssl.duitang.com&app=2002&size=f9999,10000&q=a80&n=0&g=0n&fmt=auto?sec=1670909422&t=b2df9f38251d7ba756cbf3a789d241e2",
				Mobile:   "1762556212",
				Motto:    "kefu",
				NickName: "lanlan",
				Uid:      20221113,
			},
		},
	})
}

func usersDetailHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Hello Welcome to PIM",
		"data": api.UserDetailInfo{
			IsQiYe:   false,
			Gender:   1,
			Email:    "12123232@qq.com",
			Avatar:   "http://gimg2.baidu.com/image_search/src=http%3A%2F%2Fc-ssl.duitang.com%2Fuploads%2Fblog%2F202105%2F19%2F20210519212438_ced7e.jpg&refer=http%3A%2F%2Fc-ssl.duitang.com&app=2002&size=f9999,10000&q=a80&n=0&g=0n&fmt=auto?sec=1670909422&t=b2df9f38251d7ba756cbf3a789d241e2",
			Mobile:   "1762556212",
			Motto:    "kefu",
			NickName: "lanlan",
			Uid:      20221113,
		},
	})
}

// Routers .
func Routers(r *gin.Engine) {
	login := r.Group("/api/v1/users")
	{
		login.GET("/setting", usersSettingHandler)
		login.GET("/detail", usersDetailHandler)
	}
}
