package emoticon

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func emoticonListHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "ok",
		"data":    nil,
	})
}
func emoticonSystemListHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "ok",
		"data":    nil,
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
