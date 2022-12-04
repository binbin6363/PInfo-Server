package note

import (
	"PInfo-server/log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func annexRecoverListHandler(c *gin.Context) {
	log.InfoContextf(c, "unimplemented")

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Hello Welcome to PIM",
		"data":    nil,
	})
}
func annexDeletedHandler(c *gin.Context) {
	log.InfoContextf(c, "unimplemented")

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Hello Welcome to PIM",
		"data":    nil,
	})
}
func annexRecoverHandler(c *gin.Context) {
	log.InfoContextf(c, "unimplemented")

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Hello Welcome to PIM",
		"data":    nil,
	})
}

func annexUploadHandler(c *gin.Context) {
	log.InfoContextf(c, "unimplemented")

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Hello Welcome to PIM",
		"data":    nil,
	})
}
