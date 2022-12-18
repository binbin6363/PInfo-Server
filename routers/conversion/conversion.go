package conversion

/*
import (
	"PInfo-server/api"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"net/http"
)

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
				TalkList: []*api.TalkInfo{
					{
						ID:         10001,
						Type:       1,
						ReceiverId: 10001,
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
						Name:       "anjintang",
						Avatar:     "",
						RemarkName: "anjin",
						UnreadNum:  1,
						MsgText:    "hello",
						UpdatedAt:  "2022-11-11 12:00:00",
					},
				},
			},
		})
	}

}

// Routers .
func Routers(r *gin.Engine) {
	talk := r.Group("/api/v1/talk")
	{
		talk.GET("/list", talkListHandler)
	}

}


*/
