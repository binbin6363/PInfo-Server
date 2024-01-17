package note

import (
	"PInfo-server/api"
	"PInfo-server/log"
	"PInfo-server/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func articleListHandler(c *gin.Context) {
	log.Infof("unimplemented")

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Hello Welcome to PIM",
		"data":    nil,
	})
}
func articleDetailHandler(c *gin.Context) {
	log.Infof("unimplemented")

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Hello Welcome to PIM",
		"data":    nil,
	})
}

// 更新/新增 文章请求
func articleEditorHandler(c *gin.Context) {
	log.Infof("articleEditorHandler")

	req := &api.EditArticleReq{}

	rsp, err := service.DefaultService.EditArticle(c, req)
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
	log.Infof("done articleEditorHandler")
}

func articleDeleteHandler(c *gin.Context) {
	log.Infof("unimplemented")

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Hello Welcome to PIM",
		"data":    nil,
	})
}
func articleRecoverHandler(c *gin.Context) {
	log.Infof("unimplemented")

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Hello Welcome to PIM",
		"data":    nil,
	})
}
func articleAsteriskHandler(c *gin.Context) {
	log.Infof("unimplemented")

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Hello Welcome to PIM",
		"data":    nil,
	})
}
func articleMoveHandler(c *gin.Context) {
	log.Infof("unimplemented")

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Hello Welcome to PIM",
		"data":    nil,
	})
}
func articleUploadImgHandler(c *gin.Context) {
	log.Infof("unimplemented")

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Hello Welcome to PIM",
		"data":    nil,
	})
}
func articleTagHandler(c *gin.Context) {
	log.Infof("unimplemented")

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Hello Welcome to PIM",
		"data":    nil,
	})
}
