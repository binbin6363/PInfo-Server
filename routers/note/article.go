package note

import (
	"PInfo-server/api"
	"PInfo-server/log"
	"PInfo-server/service"
	"net/http"

	"github.com/spf13/cast"

	"github.com/gin-gonic/gin"
)

func articleListHandler(c *gin.Context) {
	log.InfoContextf(c, "articleListHandler")

	req := &api.ArticleListReq{}

	req.Page = cast.ToInt(c.Query("page"))
	req.Keyword = c.Query("keyword")
	req.FindType = cast.ToInt(c.Query("find_type"))
	req.Cid = cast.ToInt64(c.Query("cid"))

	if err := c.ShouldBind(req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "参数错误",
			"data":    nil,
		})
		return
	}
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

	rsp, err := service.DefaultService.ArticleList(c, req)
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
	log.InfoContextf(c, "done articleListHandler")
}

func articleDetailHandler(c *gin.Context) {
	log.InfoContextf(c, "articleDetailHandler")

	req := &api.ArticleDetailReq{}
	req.ArticleId = cast.ToInt64(c.Query("article_id"))
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

	rsp, err := service.DefaultService.ArticleDetail(c, req)
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
	log.InfoContextf(c, "done articleDetailHandler")
}

// 更新/新增 文章请求
func articleEditorHandler(c *gin.Context) {
	log.InfoContextf(c, "articleEditorHandler")

	req := &api.ArticleEditReq{}
	if err := c.ShouldBind(req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "参数错误",
			"data":    nil,
		})
		return
	}
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

	log.DebugContextf(c, "show articleEditorHandler req: %v", req)
	rsp, err := service.DefaultService.ArticleEdit(c, req)
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
	log.InfoContextf(c, "done articleEditorHandler")
}

func articleDeleteHandler(c *gin.Context) {
	log.InfoContextf(c, "unimplemented")

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Hello Welcome to PIM",
		"data":    nil,
	})
}
func articleRecoverHandler(c *gin.Context) {
	log.InfoContextf(c, "unimplemented")

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Hello Welcome to PIM",
		"data":    nil,
	})
}
func articleAsteriskHandler(c *gin.Context) {
	log.InfoContextf(c, "unimplemented")

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Hello Welcome to PIM",
		"data":    nil,
	})
}
func articleMoveHandler(c *gin.Context) {
	log.InfoContextf(c, "articleMoveHandler")

	req := &api.ArticleMoveClassReq{} // 移动到另外分类
	if err := c.ShouldBind(req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "参数错误",
			"data":    nil,
		})
		return
	}

	rsp, err := service.DefaultService.ArticleMoveClass(c, req)
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
	log.InfoContextf(c, "done articleMoveHandler")
}
func articleUploadImgHandler(c *gin.Context) {
	log.InfoContextf(c, "unimplemented")

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Hello Welcome to PIM",
		"data":    nil,
	})
}
func articleTagHandler(c *gin.Context) {
	log.InfoContextf(c, "unimplemented")

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Hello Welcome to PIM",
		"data":    nil,
	})
}
