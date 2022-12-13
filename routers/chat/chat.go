package chat

import (
	"PInfo-server/api"
	"PInfo-server/service"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"log"
	"net/http"
)

// talkListHandler 获取聊天列表服务接口
func talkListHandler(c *gin.Context) {
	req := &api.TalkListReq{}
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

	if req.Uid == 10000 {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "Hello Welcome to PIM",
			"data": api.TalkListRsp{
				TalkList: []api.TalkInfo{
					{
						ID:         10001,
						Type:       1,
						ReceiverId: 10001,
						IsTop:      0,
						IsDisturb:  0,
						IsOnline:   0,
						IsRobot:    0,
						Name:       "anjintang",
						Avatar:     "",
						RemarkName: "anjin",
						UnreadNum:  0,
						MsgText:    "hello anjin",
						UpdatedAt:  "2022-11-12 12:00:00",
					},
				},
			},
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "Hello Welcome to PIM",
			"data": api.TalkListRsp{
				TalkList: []api.TalkInfo{
					{
						ID:         10000,
						Type:       1,
						ReceiverId: 10000,
						IsTop:      1,
						IsDisturb:  0,
						IsOnline:   0,
						IsRobot:    0,
						Name:       "politewang",
						Avatar:     "",
						RemarkName: "polite",
						UnreadNum:  1,
						MsgText:    "hello polite",
						UpdatedAt:  "2022-11-11 12:00:00",
					},
				},
			},
		})
	}
	/*
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "success",
			"data": api.TalkListRsp{
				TalkList: []api.TalkInfo{
					{
						ID:         10001,
						Type:       1,
						ReceiverId: 10001,
						IsTop:      0,
						IsDisturb:  0,
						IsOnline:   0,
						IsRobot:    0,
						Name:       "anjintang",
						Avatar:     "",
						RemarkName: "anjin",
						UnreadNum:  1,
						MsgText:    "hello polite",
						UpdatedAt:  "2022-11-11 12:00:00",
					},
					{
						ID:         10000,
						Type:       1,
						ReceiverId: 10000,
						IsTop:      0,
						IsDisturb:  0,
						IsOnline:   0,
						IsRobot:    0,
						Name:       "politewang",
						Avatar:     "",
						RemarkName: "polite",
						UnreadNum:  0,
						MsgText:    "hello anjin",
						UpdatedAt:  "2022-11-12 12:00:00",
					},
				},
			},
		})

	*/
}

// createHandler 聊天列表创建服务接口
func createHandler(c *gin.Context) {
	log.Printf("unimplemented\n")
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "unimplemented",
		"data":    nil,
	})
}

// deleteHandler 删除聊天列表服务接口
func deleteHandler(c *gin.Context) {
	log.Printf("unimplemented\n")
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "unimplemented",
		"data":    nil,
	})
}

// toppingHandler 对话列表置顶服务接口
func toppingHandler(c *gin.Context) {
	log.Printf("unimplemented\n")
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "unimplemented",
		"data":    nil,
	})
}

// unreadClearHandler 清除聊天消息未读数服务接口
func unreadClearHandler(c *gin.Context) {
	log.Printf("unimplemented\n")
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "unimplemented",
		"data":    nil,
	})
}

// recordsHandler 获取聊天记录服务接口
func recordsHandler(c *gin.Context) {
	log.Printf("unimplemented\n")
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

	/*
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "unimplemented",
			"data": api.MsgRecordsRsp{
				Limit:       30,
				MaxRecordId: 12000,
				Rows: []api.MessageRow{
					{
						Id:         1200,
						Sequence:   2,
						TalkType:   1,
						MsgType:    1,
						UserId:     20221110,
						PeerId: 20221113,
						Nickname:   "jack",
						Avatar:     "https://im.gzydong.club/public/media/image/avatar/20221124/ea1bf7400e61fad835ad72c2c9e985b1_200x200.png",
						IsRevoke:   0,
						IsMark:     0,
						IsRead:     1,
						Content:    "last msg",
						CreatedAt:  "2022-12-08 08:50:45",
					}, {
						Id:         1123,
						Sequence:   1,
						TalkType:   1,
						MsgType:    1,
						UserId:     20221110,
						PeerId: 20221113,
						Nickname:   "jack",
						Avatar:     "https://im.gzydong.club/public/media/image/avatar/20221124/ea1bf7400e61fad835ad72c2c9e985b1_200x200.png",
						IsRevoke:   0,
						IsMark:     0,
						IsRead:     1,
						Content:    "new msg",
						CreatedAt:  "2022-12-08 09:50:45",
					},
				},
			},
		})
	*/
}

// recordsForwardHandler 获取转发会话记录详情列表服务接口
func recordsForwardHandler(c *gin.Context) {
	log.Printf("unimplemented\n")
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "unimplemented",
		"data":    nil,
	})
}

// disturbHandler 对话列表置顶服务接口
func disturbHandler(c *gin.Context) {
	log.Printf("unimplemented\n")
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "unimplemented",
		"data":    nil,
	})
}

// recordsHistoryHandler 查找用户聊天记录服务接口
func recordsHistoryHandler(c *gin.Context) {
	log.Printf("unimplemented\n")
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "unimplemented",
		"data":    nil,
	})
}

// searchMsgHandler 搜索用户聊天记录服务接口
func searchMsgHandler(c *gin.Context) {
	log.Printf("unimplemented\n")
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "unimplemented",
		"data":    nil,
	})
}

// ServeGetRecordsContext .
func ServeGetRecordsContext(c *gin.Context) {
	log.Printf("unimplemented\n")
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
	err, rsp := service.DefaultService.AddOneMessage(c, req)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "内部错误",
			"data":    nil,
		})
		return
	}
	//
	//rsp := api.SendTextMsgEvtNotice{
	//	Content: api.SendTextMsgContent{
	//		Data: api.SendTextMsgData{
	//			Id:         1201,
	//			Sequence:   3,
	//			TalkType:   1,
	//			MsgType:    1,
	//			UserId:     20221113,
	//			PeerId: 20221110,
	//			Nickname:   "jack",
	//			Avatar:     "https://im.gzydong.club/public/media/image/avatar/20221124/ea1bf7400e61fad835ad72c2c9e985b1_200x200.png",
	//			IsMark:     0,
	//			IsRead:     0,
	//			Content:    "new chat content",
	//			CreatedAt:  "2022-12-09 00:50:45",
	//		},
	//		TalkType:   1,
	//		PeerId: 20221111,
	//		SenderId:   20221113,
	//	},
	//}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    rsp,
	})

}

// sendCodeMsgHandler 发送代码块消息服务接口
func sendCodeMsgHandler(c *gin.Context) {
	log.Printf("unimplemented\n")
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "unimplemented",
		"data":    nil,
	})
}

// sendFileMsgHandler 发送聊天文件服务接口
func sendFileMsgHandler(c *gin.Context) {
	log.Printf("unimplemented\n")
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "unimplemented",
		"data":    nil,
	})
}

// sendImageMsgHandler 发送聊天图片服务接口
func sendImageMsgHandler(c *gin.Context) {
	log.Printf("unimplemented\n")
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "unimplemented",
		"data":    nil,
	})
}

// sendEmoticonMsgHandler 发送表情包服务接口
func sendEmoticonMsgHandler(c *gin.Context) {
	log.Printf("unimplemented\n")
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "unimplemented",
		"data":    nil,
	})
}

// forwardMsgHandler 转发消息服务接口
func forwardMsgHandler(c *gin.Context) {
	log.Printf("unimplemented\n")
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "unimplemented",
		"data":    nil,
	})
}

// revokeMsgHandler 撤回消息服务接口
func revokeMsgHandler(c *gin.Context) {
	log.Printf("unimplemented\n")
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "unimplemented",
		"data":    nil,
	})
}

// deleteMsgHandler 删除消息服务接口
func deleteMsgHandler(c *gin.Context) {
	log.Printf("unimplemented\n")
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "unimplemented",
		"data":    nil,
	})
}

// collectMsgHandler 收藏表情包服务接口
func collectMsgHandler(c *gin.Context) {
	log.Printf("unimplemented\n")
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "unimplemented",
		"data":    nil,
	})
}

// sendVoteMsgHandler 发送投票消息服务接口
func sendVoteMsgHandler(c *gin.Context) {
	log.Printf("unimplemented\n")
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "unimplemented",
		"data":    nil,
	})
}

// confirmVoteMsgHandler 确认投票消息服务接口
func confirmVoteMsgHandler(c *gin.Context) {
	log.Printf("unimplemented\n")
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
		login.GET("/records/history", recordsHistoryHandler)
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
