package emoticon

import (
	"PInfo-server/api"
	"net/http"

	"github.com/gin-gonic/gin"
)

func emoticonListHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "ok",
		"data": &api.EmoticonListRsp{
			CollectEmoticon: nil,
			SystemEmoticon:  nil,
		},
	})
}
func emoticonSystemListHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "ok",
		"data": &api.EmoticonListRsp{
			SystemEmoticon: nil,
		},
	})
}

// Routers .
func Routers(r *gin.Engine) {
	talk := r.Group("/api/v1/emoticon")
	{
		talk.GET("/list", emoticonListHandler)
		talk.GET("/system/list", emoticonSystemListHandler)

	}

}
