package chat

import (
	"PInfo-server/api"
	"PInfo-server/config"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// talkListHandler 获取聊天列表服务接口
func talkListHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": api.TalkListRsp{
			TalkList: []api.TalkInfo{
				{
					ID:         123456789,
					Type:       1,
					ReceiverId: 20221111,
					IsTop:      0,
					IsDisturb:  0,
					IsOnline:   0,
					IsRobot:    0,
					Name:       "jack",
					Avatar:     "",
					RemarkName: "jack",
					UnreadNum:  1,
					MsgText:    "hello",
					UpdatedAt:  "2022-11-11 12:00:00",
				},
				{
					ID:         123456788,
					Type:       1,
					ReceiverId: 20221110,
					IsTop:      0,
					IsDisturb:  0,
					IsOnline:   0,
					IsRobot:    0,
					Name:       "mark",
					Avatar:     "",
					RemarkName: "mark",
					UnreadNum:  0,
					MsgText:    "hello polite",
					UpdatedAt:  "2022-11-12 12:00:00",
				},
			},
		},
	})
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
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "unimplemented",
		"data":    nil,
	})
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

	//rsp := api.SendTextMsgRsp{
	//	Id:         10001,
	//	TalkType:   1,
	//	ReceiverId: 20221114,
	//	SenderId:   20221113,
	//	Name:       "20221114",
	//	RemarkName: "mark",
	//	Avatar:     "",
	//	IsDisturb:  0,
	//	IsTop:      0,
	//	IsOnline:   0,
	//	IsRobot:    0,
	//	UnreadNum:  1,
	//	Content:    "chat content",
	//	DraftText:  "",
	//	MsgText:    "chat MsgText",
	//	IndexName:  "",
	//	CreatedAt:  "20221119",
	//}

	rsp := api.SendTextMsgEvtRsp{
		Content: api.SendTextMsgContent{
			Data: api.SendTextMsgData{
				Id:         10001,
				Sequence:   100000,
				TalkType:   1,
				MsgType:    2,
				UserId:     20221113,
				ReceiverId: 20221111,
				Nickname:   "polite",
				Avatar:     "",
				IsMark:     0,
				IsRead:     0,
				Content:    "chat content",
				CreatedAt:  "20221119",
			},
			TalkType:   1,
			ReceiverId: 20221111,
			SenderId:   20221113,
		},
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    rsp,
	})

	req := &api.SendTextMsgReq{}
	if err := c.BindJSON(req); err != nil {
		log.Printf("param error, err:%+v", err)
		return
	}

	// 通知websocket
	bytesData, _ := json.Marshal(req)
	url := fmt.Sprintf("http://%s/notice/message/text", config.AppConfig().ConnInfo.Addr)
	_, err := http.Post(url, "application/json; charset=utf-8", bytes.NewReader(bytesData))
	if err != nil {
		log.Printf("post conn failed, err:%+v\n", err)
	} else {
		log.Printf("post conn success, req:%+v\n", req)
	}
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
