package upload

import (
	"PInfo-server/api"
	"PInfo-server/service"
	"PInfo-server/utils"
	"github.com/gin-gonic/gin"
)

func uploadAvatarHandler(c *gin.Context) {
	form, _ := c.MultipartForm()
	req := &api.UploadReq{
		Form: form,
	}

	err, rsp := service.DefaultService.UploadAvatar(c.Request.Context(), req)
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
}
