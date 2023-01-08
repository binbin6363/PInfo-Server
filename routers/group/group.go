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

func memberInviteHandler(c *gin.Context) {
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
	err, rsp := service.DefaultService.InviteGroupMember(c, groupMembersReq)
	if err != nil {
		utils.SendJsonRsp(c, rsp)
		return
	}

	utils.SendJsonRsp(c, rsp)
}

func groupListHandler(c *gin.Context) {
	groupListReq := &api.GroupListReq{}
	if err := c.ShouldBind(groupListReq); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "参数错误",
			"data":    nil,
		})
		return
	}

	groupListReq.Uid = utils.GetUid(c)

	// 获取群成员列表
	err, rsp := service.DefaultService.GetGroupList(c, groupListReq)
	if err != nil {
		utils.SendJsonRsp(c, rsp)
		return
	}

	utils.SendJsonRsp(c, rsp)
}

func getGroupDetailHandler(c *gin.Context) {
	groupDetailReq := &api.GroupDetailReq{}
	if err := c.ShouldBind(groupDetailReq); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "参数错误",
			"data":    nil,
		})
		return
	}

	groupDetailReq.Uid = utils.GetUid(c)
	_, groupDetailReq.UserName = utils.GetUserName(c)
	groupDetailReq.GroupId = utils.QueryInt64Val(c, "group_id")

	// 获取群详情
	err, rsp := service.DefaultService.GetGroupDetail(c, groupDetailReq)
	if err != nil {
		utils.SendJsonRsp(c, rsp)
		return
	}

	utils.SendJsonRsp(c, rsp)
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
	setGroupInfoReq := &api.SetGroupInfoReq{}
	if err := c.ShouldBind(setGroupInfoReq); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "参数错误",
			"data":    nil,
		})
		return
	}

	setGroupInfoReq.Uid = utils.GetUid(c)
	_, setGroupInfoReq.UserName = utils.GetUserName(c)

	// 邀请进群
	err, rsp := service.DefaultService.SetGroupInfo(c, setGroupInfoReq)
	if err != nil {
		utils.SendJsonRsp(c, rsp)
		return
	}

	utils.SendJsonRsp(c, rsp)
}

func inviteGroupHandler(c *gin.Context) {
	inviteGroupReq := &api.InviteGroupReq{}
	if err := c.ShouldBind(inviteGroupReq); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "参数错误",
			"data":    nil,
		})
		return
	}

	inviteGroupReq.Uid = utils.GetUid(c)
	_, inviteGroupReq.UserName = utils.GetUserName(c)

	// 邀请进群
	err, rsp := service.DefaultService.InviteGroup(c, inviteGroupReq)
	if err != nil {
		utils.SendJsonRsp(c, rsp)
		return
	}

	utils.SendJsonRsp(c, rsp)
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

// remarkGroupHandler 修改在群中的昵称
func remarkGroupHandler(c *gin.Context) {
	remarkNameInGroupReq := &api.RemarkNameInGroupReq{}
	if err := c.ShouldBind(remarkNameInGroupReq); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "参数错误",
			"data":    nil,
		})
		return
	}

	remarkNameInGroupReq.Uid = utils.GetUid(c)
	if remarkNameInGroupReq.GroupId == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "参数错误",
			"data":    nil,
		})
		return
	}

	// 获取群成员列表
	err, rsp := service.DefaultService.RemarkNameInGroup(c, remarkNameInGroupReq)
	if err != nil {
		utils.SendJsonRsp(c, rsp)
		return
	}

	utils.SendJsonRsp(c, rsp)
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

func getGroupMembersHandler(c *gin.Context) {
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
	if groupMembersReq.GroupId == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "参数错误",
			"data":    nil,
		})
		return
	}

	// 获取群成员列表
	err, rsp := service.DefaultService.GetGroupMembers(c, groupMembersReq)
	if err != nil {
		utils.SendJsonRsp(c, rsp)
		return
	}

	utils.SendJsonRsp(c, rsp)
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
		group.GET("/member/invites", memberInviteHandler)
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
		group.GET("/member/list", getGroupMembersHandler)
		group.POST("/notice/edit", editGroupNoticeHandler)

	}
}
