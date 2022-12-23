package contact

import (
	"PInfo-server/api"
	"PInfo-server/model"
	"PInfo-server/service"
	"PInfo-server/utils"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"log"
	"net/http"
)

// contactListHandler 获取好友列表服务接口
func contactListHandler(c *gin.Context) {
	err, contactInfos := service.DefaultService.GetContactList(c, utils.GetUid(c), model.ContactFriend)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusNotFound,
			"message": "内部错误，请重试",
			"data":    nil,
		})
		return
	}

	contactList := api.ContactListRsp{}
	for _, contact := range contactInfos {
		contactInfo := &api.ContactInfo{
			Id:           contact.Id,
			Nickname:     contact.Nickname,
			Gender:       contact.Gender,
			Motto:        contact.Motto,
			Avatar:       contact.Avatar,
			FriendRemark: contact.FriendRemark,
			IsOnline:     1,
		}
		contactList.ContactList = append(contactList.ContactList, contactInfo)
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "ok",
		"data":    contactList,
	})
}

// contactSearchHandler 搜索联系人
func contactSearchHandler(c *gin.Context) {
	req := &api.ContactSearchReq{}
	if err := c.ShouldBind(req); err != nil {
		utils.SendJsonRsp(c, &api.CommRsp{
			Code:    400,
			Message: "param invalid",
			Data:    nil,
		})
		return
	}

	req.UserName = c.Query("mobile")
	err, rsp := service.DefaultService.ContactSearch(c, req)
	if err != nil {
		utils.SendJsonRsp(c, rsp)
		return
	}

	utils.SendJsonRsp(c, rsp)
}

// deleteContactHandler 解除好友关系服务接口
func deleteContactHandler(c *gin.Context) {
	log.Printf("unimplemented")
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Hello Welcome to PIM",
		"data":    nil,
	})
}

// editContactRemarkHandler 修改好友备注服务接口
func editContactRemarkHandler(c *gin.Context) {
	req := &api.EditContactInfoReq{}
	if err := c.ShouldBind(req); err != nil {
		utils.SendJsonRsp(c, &api.CommRsp{
			Code:    400,
			Message: "param invalid",
			Data:    nil,
		})
		return
	}

	req.Uid = utils.GetUid(c)
	_, req.UserName = utils.GetUserName(c)
	err, rsp := service.DefaultService.EditContactInfo(c, req)
	if err != nil {
		utils.SendJsonRsp(c, rsp)
		return
	}

	utils.SendJsonRsp(c, rsp)

}

// contactDetailHandler 搜索用户信息服务接口
func contactDetailHandler(c *gin.Context) {

	req := &api.ContactDetailReq{}
	if err := c.ShouldBind(req); err != nil {
		utils.SendJsonRsp(c, &api.CommRsp{
			Code:    400,
			Message: "param invalid",
			Data:    nil,
		})
		return
	}

	req.Uid = utils.GetUid(c)
	req.ContactId = cast.ToInt64(c.Query("user_id"))
	err, rsp := service.DefaultService.ContactDetail(c, req)
	if err != nil {
		utils.SendJsonRsp(c, rsp)
		return
	}

	utils.SendJsonRsp(c, rsp)
}

// unreadNumHandler 查询好友申请未读数量服务接口
func unreadNumHandler(c *gin.Context) {
	log.Printf("unimplemented")
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
	// todo:
	err, contactInfos := service.DefaultService.GetContactList(c, utils.GetUid(c), model.ContactWaitMeApply)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusNotFound,
			"message": "内部错误，请重试",
			"data":    nil,
		})
		return
	}

	contactList := api.ContactListRsp{}
	for _, contact := range contactInfos {
		contactInfo := &api.ContactInfo{
			Id:        contact.Id,
			UserId:    contact.Id,
			Nickname:  contact.Nickname,
			Gender:    contact.Gender,
			Motto:     contact.Motto,
			Avatar:    contact.Avatar,
			Remark:    contact.Remark,
			IsOnline:  1,
			CreatedAt: contact.CreatedAt,
		}
		contactList.ContactList = append(contactList.ContactList, contactInfo)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "ok",
		"data":    contactList,
	})

}

// addContactHandler 好友申请服务接口
func addContactHandler(c *gin.Context) {

	req := &api.AddContactReq{}
	if err := c.ShouldBind(req); err != nil {
		utils.SendJsonRsp(c, &api.CommRsp{
			Code:    400,
			Message: "param invalid",
			Data:    nil,
		})
		return
	}

	req.Uid = utils.GetUid(c)
	_, req.UserName = utils.GetUserName(c)
	err, rsp := service.DefaultService.AddContact(c, req)
	if err != nil {
		utils.SendJsonRsp(c, rsp)
		return
	}

	utils.SendJsonRsp(c, rsp)

	/*
		log.Printf("unimplemented")
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "Hello Welcome to PIM",
			"data":    nil,
		})

	*/
}

// acceptContactHandler 处理好友申请服务接口，同意
func acceptContactHandler(c *gin.Context) {

	req := &api.ApplyAddContactReq{}
	if err := c.ShouldBind(req); err != nil {
		utils.SendJsonRsp(c, &api.CommRsp{
			Code:    400,
			Message: "param invalid",
			Data:    nil,
		})
		return
	}

	req.Uid = utils.GetUid(c)
	_, req.UserName = utils.GetUserName(c)
	err, rsp := service.DefaultService.ApplyAddContact(c, req)
	if err != nil {
		utils.SendJsonRsp(c, rsp)
		return
	}

	utils.SendJsonRsp(c, rsp)

}

// declineContactHandler 处理好友申请服务接口，拒绝
func declineContactHandler(c *gin.Context) {

	req := &api.ApplyAddContactReq{}
	if err := c.ShouldBind(req); err != nil {
		utils.SendJsonRsp(c, &api.CommRsp{
			Code:    400,
			Message: "param invalid",
			Data:    nil,
		})
		return
	}

	req.Uid = utils.GetUid(c)
	_, req.UserName = utils.GetUserName(c)
	err, rsp := service.DefaultService.DeclineAddContact(c, req)
	if err != nil {
		utils.SendJsonRsp(c, rsp)
		return
	}

	utils.SendJsonRsp(c, rsp)

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
		contact.POST("/apply/decline", declineContactHandler)
	}
}
