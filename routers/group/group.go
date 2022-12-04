package group

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func groupInvitesHandler(c *gin.Context) {
	log.Printf("unimplemented\n")

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Hello Welcome to PIM",
		"data":    nil,
	})
}

func groupListHandler(c *gin.Context) {
	log.Printf("unimplemented\n")

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Hello Welcome to PIM",
		"data":    nil,
	})
}

func getGroupDetailHandler(c *gin.Context) {
	log.Printf("unimplemented\n")

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Hello Welcome to PIM",
		"data":    nil,
	})
}

func createGroupHandler(c *gin.Context) {
	log.Printf("unimplemented\n")

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Hello Welcome to PIM",
		"data":    nil,
	})
}

func modifyGroupHandler(c *gin.Context) {
	log.Printf("unimplemented\n")

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Hello Welcome to PIM",
		"data":    nil,
	})
}

func inviteGroupHandler(c *gin.Context) {
	log.Printf("unimplemented\n")

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Hello Welcome to PIM",
		"data":    nil,
	})
}

func removeGroupMemberHandler(c *gin.Context) {
	log.Printf("unimplemented\n")

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Hello Welcome to PIM",
		"data":    nil,
	})
}

func dismissGroupHandler(c *gin.Context) {
	log.Printf("unimplemented\n")

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Hello Welcome to PIM",
		"data":    nil,
	})
}

func secedeGroupHandler(c *gin.Context) {
	log.Printf("unimplemented\n")

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Hello Welcome to PIM",
		"data":    nil,
	})
}

func remarkGroupHandler(c *gin.Context) {
	log.Printf("unimplemented\n")

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Hello Welcome to PIM",
		"data":    nil,
	})
}

func getGroupMemberListHandler(c *gin.Context) {
	log.Printf("unimplemented\n")

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Hello Welcome to PIM",
		"data":    nil,
	})
}

func getGroupNoticesHandler(c *gin.Context) {
	log.Printf("unimplemented\n")

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Hello Welcome to PIM",
		"data":    nil,
	})
}

func editGroupNoticeHandler(c *gin.Context) {
	log.Printf("unimplemented\n")
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Hello Welcome to PIM",
		"data":    nil,
	})
}

// Routers .
func Routers(r *gin.Engine) {
	login := r.Group("/api/v1/group/")
	{
		login.GET("/member/invites", groupInvitesHandler)
		login.GET("/list", groupListHandler)
		login.GET("/detail", getGroupDetailHandler)
		login.POST("/create", createGroupHandler)
		login.POST("/setting", modifyGroupHandler)
		login.POST("/invite", inviteGroupHandler)
		login.POST("/member/remove", removeGroupMemberHandler)
		login.POST("/dismiss", dismissGroupHandler)
		login.POST("/secede", secedeGroupHandler)
		login.POST("/member/remark", remarkGroupHandler)
		//login.GET("/member/invites", getGroupMemberListHandler)
		login.GET("/notice/list", getGroupNoticesHandler)
		login.POST("/notice/edit", editGroupNoticeHandler)

	}
}
