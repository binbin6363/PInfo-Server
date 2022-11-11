package im

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func addFriendHandler(c *gin.Context) {
	result := make(map[string]string)
	result["name"] = "add friend"
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "Hello Welcome to barter-economy",
		"data":result,
	})
}

func sendMsgHandler(c *gin.Context) {
	result := make(map[string]string)
	result["name"] = "send msg"
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "Hello Welcome to barter-economy",
		"data":result,
	})
}

func listFriendHandler(c *gin.Context) {
	result := make(map[string]string)
	result["name"] = "list friend"
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "Hello Welcome to barter-economy",
		"data":result,
	})
}

func delFriendHandler(c *gin.Context) {
	result := make(map[string]string)
	result["name"] = "del friend"
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "Hello Welcome to barter-economy",
		"data":result,
	})
}

// test add get
// Routers
func Routers(r *gin.Engine) {
	users := r.Group("/im")
	{
		users.GET("/friend/list", listFriendHandler)
		users.PUT("/friend/add", addFriendHandler)
		users.DELETE("/friend/del", delFriendHandler)

		users.POST("/send_msg", sendMsgHandler)

	}
}
