package service

import (
	"PInfo-server/api"
	"PInfo-server/config"
	"PInfo-server/log"
	"context"
	"io"
	"os"
	"path"
)

func (s *Service) UploadAvatar(ctx context.Context, req *api.UploadReq) (err error, rsp *api.CommRsp) {
	rsp = &api.CommRsp{
		Code:    0,
		Message: "ok",
		Data:    nil,
	}

	uploadRsp := &api.UploadRsp{}
	uploadOK := false
	resourceRoot := config.AppConfig().ServerInfo.ResourceRoot
	remoteUrlRoot := config.AppConfig().ServerInfo.RemoteUrlRoot
	for _, fileHeaders := range req.Form.File {
		for _, file := range fileHeaders {
			inFile, err := file.Open()
			if err != nil {
				log.Errorf("open infile failed, path:%s, err:%v", file.Filename, err)
				continue
			}
			defer inFile.Close()
			fullPath := path.Join(resourceRoot, file.Filename)
			out, err := os.Create(fullPath)
			if err != nil {
				log.Errorf("open outfile failed, path:%s, err:%v", fullPath, err)
				continue
			}
			defer out.Close()
			_, err = io.Copy(out, inFile)
			if err != nil {
				log.Errorf("copy file failed, infile:%s, outfile:%s, err:%v", file.Filename, fullPath, err)
				continue
			}
			uploadOK = true
			uploadRsp.Avatar = path.Join(remoteUrlRoot, file.Filename)
			log.Infof("upload file ok, src:%s, target path:%s, size:%d", file.Filename, fullPath, file.Size)
		}
	}

	if uploadOK {
		rsp.Data = uploadRsp
	} else {
		rsp.Code = 5000
		rsp.Message = "文件上传失败"
	}

	log.Infof("UploadAvatar ok, req:%+v, rsp:%+v", req, uploadRsp)
	return nil, rsp
}
