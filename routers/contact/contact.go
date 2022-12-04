package contact

import (
	"PInfo-server/api"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// contactListHandler 获取好友列表服务接口
func contactListHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Hello Welcome to PIM",
		"data":    api.ContactListRsp{},
	})
}

// contactSearchHandler 搜索联系人
func contactSearchHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Hello Welcome to PIM",
		"data": api.UserDetailInfo{
			IsQiYe:   false,
			Gender:   1,
			Email:    "12123232@qq.com",
			Avatar:   "",
			Mobile:   "1762556212",
			Motto:    "kefu",
			NickName: "lanlan",
			Uid:      20221113,
		},
	})
}

// deleteContactHandler 解除好友关系服务接口
func deleteContactHandler(c *gin.Context) {
	log.Printf("unimplemented\n")
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Hello Welcome to PIM",
		"data":    nil,
	})
}

// editContactRemarkHandler 修改好友备注服务接口
func editContactRemarkHandler(c *gin.Context) {
	log.Printf("unimplemented\n")
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Hello Welcome to PIM",
		"data":    nil,
	})
}

// contactDetailHandler 搜索用户信息服务接口
func contactDetailHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Hello Welcome to PIM",
		"data": api.UserDetailInfo{
			IsQiYe:   false,
			Gender:   1,
			Email:    "12123232@qq.com",
			Avatar:   "",
			Mobile:   "1762556212",
			Motto:    "kefu",
			NickName: "lanlan",
			Uid:      20221113,
		},
	})
}

// unreadNumHandler 查询好友申请未读数量服务接口
func unreadNumHandler(c *gin.Context) {
	log.Printf("unimplemented\n")
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Hello Welcome to PIM",
		"data": api.UnReadNumRsp{
			UnreadNum: 1,
		},
	})
}

// recordsHandler 查询好友申请服务接口
func recordsHandler(c *gin.Context) {
	log.Printf("unimplemented\n")
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Hello Welcome to PIM",
		"data":    nil,
	})
}

// addContactHandler 好友申请服务接口
func addContactHandler(c *gin.Context) {
	log.Printf("unimplemented\n")
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Hello Welcome to PIM",
		"data":    nil,
	})
}

// acceptContactHandler 处理好友申请服务接口
func acceptContactHandler(c *gin.Context) {
	log.Printf("unimplemented\n")
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Hello Welcome to PIM",
		"data":    nil,
	})
}

// Routers .
func Routers(r *gin.Engine) {
	contact := r.Group("/api/v1/contact")
	{
		contact.GET("/list", contactListHandler)
		contact.GET("/search", contactSearchHandler)
		contact.POST("/delete", deleteContactHandler)
		contact.POST("/edit-remark", editContactRemarkHandler)
		contact.GET("/detail", contactDetailHandler)
		contact.GET("/apply/unread-num", unreadNumHandler)
		contact.GET("/apply/records", recordsHandler)
		contact.POST("/apply/create", addContactHandler)
		contact.POST("/apply/accept", acceptContactHandler)
	}
}
