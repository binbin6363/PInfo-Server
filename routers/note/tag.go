package note

import (
	"PInfo-server/log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func tagListHandler(c *gin.Context) {
	log.InfoContextf(c, "unimplemented")

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Hello Welcome to PIM",
		"data":    nil,
	})
}

func tagEditorHandler(c *gin.Context) {
	log.InfoContextf(c, "unimplemented")

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Hello Welcome to PIM",
		"data":    nil,
	})
}
func tagDeleteHandler(c *gin.Context) {
	log.InfoContextf(c, "unimplemented")

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Hello Welcome to PIM",
		"data":    nil,
	})
}
