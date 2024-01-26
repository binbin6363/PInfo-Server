package note

import (
	"PInfo-server/routers"

	"github.com/gin-gonic/gin"
)

func init() {
	routers.Register(Routers)
}

// Routers .
func Routers(r *gin.Engine) {
	note := r.Group("/api/v1/note")
	{
		// -------- 笔记分类相关 --------
		note.GET("/class/list", classListHandler)      // 查询用户文集分类服务接口
		note.POST("/class/editor", classEditorHandler) // 添加或编辑文集分类服务接口
		note.POST("/class/delete", classDeleteHandler) // 删除笔记分类服务接口
		note.POST("/class/sort", classSortHandler)     // 笔记分类排序服务接口
		note.POST("/class/merge", classMergeHandler)   // 合并笔记分类服务接口

		// -------- 笔记相关 --------
		note.GET("/article/list", articleListHandler)               // 查询用户文集分类服务接口
		note.GET("/article/detail", articleDetailHandler)           // 查询用户文集分类服务接口
		note.POST("/article/editor", articleEditorHandler)          // 编辑笔记服务接口
		note.POST("/article/delete", articleDeleteHandler)          // 删除笔记服务接口
		note.POST("/article/recover", articleRecoverHandler)        // 恢复笔记服务接口
		note.POST("/article/asterisk", articleAsteriskHandler)      // 设置标记星号笔记服务接口
		note.POST("/article/move", articleMoveHandler)              // 移动笔记服务接口
		note.POST("/article/upload/image", articleUploadImgHandler) // 笔记图片上传服务接口
		note.POST("/article/tag", articleTagHandler)                // 更新笔记标签服务接口

		// -------- 笔记标签相关 --------
		note.GET("/tag/list", tagListHandler)      // 获取笔记表标签服务接口
		note.POST("/tag/editor", tagEditorHandler) // 添加或编辑笔记标签服务接口
		note.POST("/tag/delete", tagDeleteHandler) // 删除笔记标签服务接口

		// -------- 笔记附件相关 --------
		note.POST("/annex/upload", annexUploadHandler)           // 笔记附件上传服务接口
		note.POST("/annex/delete", annexDeletedHandler)          // 删除笔记附件服务接口
		note.POST("/annex/recover", annexRecoverHandler)         // 恢复笔记附件服务接口
		note.GET("/annex/recover/list", annexRecoverListHandler) // 笔记附件回收站列表服务接口
	}
}
