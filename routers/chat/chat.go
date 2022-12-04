package chat

import (
	"PInfo-server/api"
	"PInfo-server/log"
	"PInfo-server/service"
	"PInfo-server/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

// talkListHandler 获取聊天列表服务接口
func talkListHandler(c *gin.Context) {
	req := &api.TalkListReq{
		Uid: utils.GetUid(c),
	}
	_, req.UserName = utils.GetUserName(c)
	// 拉取会话列表
	err, rsp := service.DefaultService.GetConversationList(c, req)
	if err != nil {
		utils.SendJsonRsp(c, rsp)
		return
	}

	utils.SendJsonRsp(c, rsp)

}

// createHandler 聊天列表创建服务接口
func createHandler(c *gin.Context) {
	req := &api.CreateTalkReq{
		Uid: utils.GetUid(c),
	}
	_, req.UserName = utils.GetUserName(c)
	if err := c.ShouldBind(req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "参数错误",
			"data":    nil,
		})
		return
	}

	// 创建会话
	err, rsp := service.DefaultService.CreateConversation(c, req)
	if err != nil {
		utils.SendJsonRsp(c, rsp)
		return
	}

	utils.SendJsonRsp(c, rsp)

}

// deleteHandler 删除聊天列表服务接口
func deleteHandler(c *gin.Context) {
	log.InfoContextf(c, "unimplemented")
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "unimplemented",
		"data":    nil,
	})
}

// toppingHandler 对话列表置顶服务接口
func toppingHandler(c *gin.Context) {
	log.InfoContextf(c, "unimplemented")
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "unimplemented",
		"data":    nil,
	})
}

// unreadClearHandler 清除聊天消息未读数服务接口
func unreadClearHandler(c *gin.Context) {
	log.InfoContextf(c, "unimplemented")
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "unimplemented",
		"data":    nil,
	})
}

// recordsHandler 获取聊天记录服务接口
func recordsHandler(c *gin.Context) {
	log.InfoContextf(c, "unimplemented")
	req := &api.MsgRecordsReq{}
	if err := c.ShouldBind(req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "参数错误",
			"data":    nil,
		})
		return
	}

	req.MinMsgId = cast.ToInt64(c.Query("record_id"))
	req.PeerId = cast.ToInt64(c.Query("receiver_id"))
	req.TalkType = cast.ToInt(c.Query("talk_type"))
	req.MsgType = cast.ToInt(c.Query("msg_type"))
	req.Limit = cast.ToInt(c.Query("limit"))
	if uid, ok := c.Get("uid"); ok {
		req.Uid = cast.ToInt64(uid)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "参数错误",
			"data":    nil,
		})
		return
	}

	// 拉取消息记录
	err, rsp := service.DefaultService.QueryMessage(c, req)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "内部错误",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    rsp,
	})
}

// recordsForwardHandler 获取转发会话记录详情列表服务接口
func recordsForwardHandler(c *gin.Context) {
	log.InfoContextf(c, "unimplemented")
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "unimplemented",
		"data":    nil,
	})
}

// disturbHandler 对话列表置顶服务接口
func disturbHandler(c *gin.Context) {
	log.InfoContextf(c, "unimplemented")
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "unimplemented",
		"data":    nil,
	})
}

// recordsHistoryHandler 查找用户聊天记录服务接口
func recordsHistoryHandler(c *gin.Context) {
	log.InfoContextf(c, "unimplemented")
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "unimplemented",
		"data":    nil,
	})
}

// searchMsgHandler 搜索用户聊天记录服务接口
func searchMsgHandler(c *gin.Context) {
	log.InfoContextf(c, "unimplemented")
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "unimplemented",
		"data":    nil,
	})
}

// ServeGetRecordsContext .
func ServeGetRecordsContext(c *gin.Context) {
	log.InfoContextf(c, "unimplemented")
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "unimplemented",
		"data":    nil,
	})
}

// sendTextMsgHandler 发送文本消息服务接口
func sendTextMsgHandler(c *gin.Context) {

	req := &api.SendTextMsgReq{}
	if err := c.ShouldBind(req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "参数错误",
			"data":    nil,
		})
		return
	}

	if uid, ok := c.Get("uid"); ok {
		req.Uid = cast.ToInt64(uid)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "参数错误",
			"data":    nil,
		})
		return
	}

	// 消息去重，暂时先基于mysql做去重
	err, rsp := service.DefaultService.SendTextMessage(c, req)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    500,
			"message": "内部错误",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    rsp,
	})

}

// sendCodeMsgHandler 发送代码块消息服务接口
func sendCodeMsgHandler(c *gin.Context) {
	log.InfoContextf(c, "unimplemented")
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "unimplemented",
		"data":    nil,
	})
}

// sendFileMsgHandler 发送聊天文件服务接口
func sendFileMsgHandler(c *gin.Context) {
	log.InfoContextf(c, "unimplemented")
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "unimplemented",
		"data":    nil,
	})
}

// sendImageMsgHandler 发送聊天图片服务接口
func sendImageMsgHandler(c *gin.Context) {
	form, _ := c.MultipartForm()
	if len(form.File) == 0 {
		log.InfoContextf(c, "no file specified!")
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "no file specified",
			"data":    nil,
		})
		return
	}

	log.InfoContextf(c, "show form")
	for key, values := range form.File {
		log.InfoContextf(c, "%s: %v\n", key, values)
	}

	req := &api.SendImageMsgReq{}

	for key, values := range form.Value {
		if key == "receiver_id" {
			req.ReceiverId = cast.ToInt64(values[0])
		} else if key == "talk_type" {
			req.TalkType = cast.ToInt(values[0])
		} else if key == "client_msg_id" {
			req.ClientMsgId = cast.ToInt64(values[0])
		}
	}
	req.Form = form
	if uid, ok := c.Get("uid"); ok {
		req.Uid = cast.ToInt64(uid)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "参数错误",
			"data":    nil,
		})
		return
	}
	// 上传文件到服务器
	log.Debugf("show req: %v", req)
	// uid 没解析出来
	rsp, err := service.DefaultService.SendImageMessage(c, req)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    500,
			"message": "内部错误",
			"data":    nil,
		})
		return
	}

	log.InfoContextf(c, "done send image message")
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    rsp,
	})
}

// sendEmoticonMsgHandler 发送表情包服务接口
func sendEmoticonMsgHandler(c *gin.Context) {
	log.InfoContextf(c, "unimplemented")
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "unimplemented",
		"data":    nil,
	})
}

// forwardMsgHandler 转发消息服务接口
func forwardMsgHandler(c *gin.Context) {
	log.InfoContextf(c, "unimplemented")
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "unimplemented",
		"data":    nil,
	})
}

// revokeMsgHandler 撤回消息服务接口
func revokeMsgHandler(c *gin.Context) {
	log.InfoContextf(c, "unimplemented")
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "unimplemented",
		"data":    nil,
	})
}

// deleteMsgHandler 删除消息服务接口
func deleteMsgHandler(c *gin.Context) {
	log.InfoContextf(c, "unimplemented")
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "unimplemented",
		"data":    nil,
	})
}

// collectMsgHandler 收藏表情包服务接口
func collectMsgHandler(c *gin.Context) {
	log.InfoContextf(c, "unimplemented")
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "unimplemented",
		"data":    nil,
	})
}

// sendVoteMsgHandler 发送投票消息服务接口
func sendVoteMsgHandler(c *gin.Context) {
	log.InfoContextf(c, "unimplemented")
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "unimplemented",
		"data":    nil,
	})
}

// confirmVoteMsgHandler 确认投票消息服务接口
func confirmVoteMsgHandler(c *gin.Context) {
	log.InfoContextf(c, "unimplemented")
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "unimplemented",
		"data":    nil,
	})
}

// Routers .
func Routers(r *gin.Engine) {
	login := r.Group("/api/v1/talk")
	{
		login.GET("/list", talkListHandler)
		login.POST("/create", createHandler)
		login.POST("/delete", deleteHandler)
		login.POST("/topping", toppingHandler)
		login.POST("/unread/clear", unreadClearHandler)
		login.GET("/records", recordsHandler)
		login.GET("/records/forward", recordsForwardHandler)
		login.POST("/disturb", disturbHandler)
		login.GET("/records/history", recordsHandler)
		login.GET("/search-chat-records", searchMsgHandler)
		login.POST("/message/text", sendTextMsgHandler)
		login.POST("/message/code", sendCodeMsgHandler)
		login.POST("/message/file", sendFileMsgHandler)
		login.POST("/message/image", sendImageMsgHandler)
		login.POST("/message/emoticon", sendEmoticonMsgHandler)
		login.POST("/message/forward", forwardMsgHandler)
		login.POST("/message/revoke", revokeMsgHandler)
		login.POST("/message/delete", deleteMsgHandler)
		login.POST("/message/collect", collectMsgHandler)
		login.POST("/message/vote", sendVoteMsgHandler)
		login.POST("/message/vote/handle", confirmVoteMsgHandler)

	}
}
