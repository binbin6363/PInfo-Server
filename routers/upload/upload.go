package upload

import (
	"PInfo-server/api"
	"PInfo-server/log"
	"PInfo-server/routers"
	"PInfo-server/service"
	"PInfo-server/utils"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

func init() {
	routers.Register(Routers)
}

func uploadAvatarHandler(c *gin.Context) {
	form, _ := c.MultipartForm()
	req := &api.UploadReq{
		Form: form,
	}

	if uid, ok := c.Get("uid"); ok {
		req.Uid = cast.ToInt64(uid)
	} else {
		log.Errorf("invalid uid req")
		return
	}

	err, rsp := service.DefaultService.UploadAvatar(c.Request.Context(), req)
	if err != nil {
		utils.SendJsonRsp(c, rsp)
		return
	}

	utils.SendJsonRsp(c, rsp)
}

func downloadAvatar(c *gin.Context) {
	req := &api.DownloadReq{}

	req.Url = c.Query("url")
	if uid, ok := c.Get("uid"); ok {
		req.Uid = cast.ToInt64(uid)
	} else {
		log.Errorf("invalid uid req")
		return
	}
	err, rsp := service.DefaultService.DownloadAvatar(c.Request.Context(), req)
	if err != nil {
		utils.SendJsonRsp(c, rsp)
		return
	}

	utils.SendJsonRsp(c, rsp)
}

// Routers .
func Routers(r *gin.Engine) {
	login := r.Group("/api/v1/upload")
	{
		login.POST("/avatar", uploadAvatarHandler)
	}

	download := r.Group("/api/v1/download")
	{
		download.GET("/img", downloadAvatar)
	}
}
