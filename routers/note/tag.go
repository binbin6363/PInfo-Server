package note

import (
	"PInfo-server/api"
	"PInfo-server/log"
	"PInfo-server/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func tagListHandler(c *gin.Context) {
	log.InfoContextf(c, "tagListHandler")

	req := &api.TagListReq{}
	if err := c.ShouldBind(req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "参数错误",
			"data":    nil,
		})
		return
	}

	rsp, err := service.DefaultService.TagList(c, req)
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
	log.InfoContextf(c, "done tagListHandler")
}

func tagEditorHandler(c *gin.Context) {
	log.InfoContextf(c, "tagEditorHandler")

	req := &api.TagEditReq{}
	if err := c.ShouldBind(req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "参数错误",
			"data":    nil,
		})
		return
	}

	rsp, err := service.DefaultService.TagEdit(c, req)
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
	log.InfoContextf(c, "done tagEditorHandler")
}

func tagDeleteHandler(c *gin.Context) {
	log.InfoContextf(c, "tagDeleteHandler")

	req := &api.TagDeleteReq{}
	if err := c.ShouldBind(req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "参数错误",
			"data":    nil,
		})
		return
	}

	rsp, err := service.DefaultService.TagDelete(c, req)
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
	log.InfoContextf(c, "done tagDeleteHandler")
}
