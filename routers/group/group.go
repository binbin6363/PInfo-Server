package group

import (
	"PInfo-server/api"
	"PInfo-server/service"
	"PInfo-server/utils"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"log"
	"net/http"
)

func groupMembersHandler(c *gin.Context) {
	groupMembersReq := &api.GroupMembersReq{}
	if err := c.ShouldBind(groupMembersReq); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "参数错误",
			"data":    nil,
		})
		return
	}

	groupMembersReq.Uid = utils.GetUid(c)
	groupMembersReq.GroupId = cast.ToInt64(c.Query("group_id"))

	// 获取群成员列表
	err, rsp := service.DefaultService.GetGroupMembers(c, groupMembersReq)
	if err != nil {
		utils.SendJsonRsp(c, rsp)
		return
	}

	utils.SendJsonRsp(c, rsp)
}

func groupListHandler(c *gin.Context) {
	log.Printf("unimplemented")

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Hello Welcome to PIM",
		"data":    nil,
	})
}

func getGroupDetailHandler(c *gin.Context) {
	log.Printf("unimplemented")

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Hello Welcome to PIM",
		"data":    nil,
	})
}

func createGroupHandler(c *gin.Context) {
	createGroupReq := &api.CreateGroupReq{}
	if err := c.ShouldBind(createGroupReq); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "参数错误",
			"data":    nil,
		})
		return
	}

	createGroupReq.Uid = utils.GetUid(c)
	_, createGroupReq.UserName = utils.GetUserName(c)

	// 创建群
	err, rsp := service.DefaultService.CreateGroup(c, createGroupReq)
	if err != nil {
		utils.SendJsonRsp(c, rsp)
		return
	}

	utils.SendJsonRsp(c, rsp)
}

func modifyGroupHandler(c *gin.Context) {
	log.Printf("unimplemented")

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Hello Welcome to PIM",
		"data":    nil,
	})
}

func inviteGroupHandler(c *gin.Context) {
	log.Printf("unimplemented")

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Hello Welcome to PIM",
		"data":    nil,
	})
}

func removeGroupMemberHandler(c *gin.Context) {
	log.Printf("unimplemented")

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Hello Welcome to PIM",
		"data":    nil,
	})
}

func dismissGroupHandler(c *gin.Context) {
	log.Printf("unimplemented")

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Hello Welcome to PIM",
		"data":    nil,
	})
}

func secedeGroupHandler(c *gin.Context) {
	log.Printf("unimplemented")

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Hello Welcome to PIM",
		"data":    nil,
	})
}

func remarkGroupHandler(c *gin.Context) {
	log.Printf("unimplemented")

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Hello Welcome to PIM",
		"data":    nil,
	})
}

func getGroupMemberListHandler(c *gin.Context) {
	log.Printf("unimplemented")

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Hello Welcome to PIM",
		"data":    nil,
	})
}

func getGroupNoticesHandler(c *gin.Context) {
	log.Printf("unimplemented")

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Hello Welcome to PIM",
		"data":    nil,
	})
}

func editGroupNoticeHandler(c *gin.Context) {
	log.Printf("unimplemented")
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Hello Welcome to PIM",
		"data":    nil,
	})
}

// Routers .
func Routers(r *gin.Engine) {
	group := r.Group("/api/v1/group/")
	{
		group.GET("/member/invites", groupMembersHandler)
		group.GET("/list", groupListHandler)
		group.GET("/detail", getGroupDetailHandler)
		group.POST("/create", createGroupHandler)
		group.POST("/setting", modifyGroupHandler)
		group.POST("/invite", inviteGroupHandler)
		group.POST("/member/remove", removeGroupMemberHandler)
		group.POST("/dismiss", dismissGroupHandler)
		group.POST("/secede", secedeGroupHandler)
		group.POST("/member/remark", remarkGroupHandler)
		//group.GET("/member/invites", getGroupMemberListHandler)
		group.GET("/notice/list", getGroupNoticesHandler)
		group.POST("/notice/edit", editGroupNoticeHandler)

	}
}
