package note

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func noteClassListHandler(c *gin.Context) {
	log.Printf("unimplemented\n")

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Hello Welcome to PIM",
		"data":    nil,
	})
}

func noteTagListHandler(c *gin.Context) {
	log.Printf("unimplemented\n")

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Hello Welcome to PIM",
		"data":    nil,
	})
}

func noteArticleListHandler(c *gin.Context) {
	log.Printf("unimplemented\n")

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Hello Welcome to PIM",
		"data":    nil,
	})
}

// Routers .
func Routers(r *gin.Engine) {
	login := r.Group("/api/v1/note")
	{
		login.GET("/class/list", noteClassListHandler)
		login.GET("/tag/list", noteTagListHandler)
		login.GET("/article/list", noteArticleListHandler)

	}
}
