package note

import (
	"PInfo-server/api"
	"PInfo-server/log"
	"PInfo-server/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func classListHandler(c *gin.Context) {
	log.InfoContextf(c, "classListHandler")

	req := &api.ClassListReq{}

	rsp, err := service.DefaultService.ClassList(c, req)
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
	log.InfoContextf(c, "done classListHandler")
}

func classEditorHandler(c *gin.Context) {
	log.InfoContextf(c, "unimplemented")

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Hello Welcome to PIM",
		"data":    nil,
	})
}
func classDeleteHandler(c *gin.Context) {
	log.InfoContextf(c, "unimplemented")

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Hello Welcome to PIM",
		"data":    nil,
	})
}
func classSortHandler(c *gin.Context) {
	log.InfoContextf(c, "unimplemented")

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Hello Welcome to PIM",
		"data":    nil,
	})
}
func classMergeHandler(c *gin.Context) {
	log.InfoContextf(c, "unimplemented")

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Hello Welcome to PIM",
		"data":    nil,
	})
}
