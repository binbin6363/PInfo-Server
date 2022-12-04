package conversion

import (
	"PInfo-server/api"
	"github.com/gin-gonic/gin"
	"net/http"
)

func talkListHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Hello Welcome to PIM",
		"data": api.TalkListRsp{
			TalkList: []api.TalkInfo{
				{
					ID:         123456789,
					Type:       1,
					ReceiverId: 20221113,
					IsTop:      1,
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
					ID:         20221114,
					Type:       1,
					ReceiverId: 20221113,
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

// Routers .
func Routers(r *gin.Engine) {
	talk := r.Group("/api/v1/talk")
	{
		talk.GET("/list", talkListHandler)
	}

}
