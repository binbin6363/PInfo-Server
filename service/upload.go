package service

import (
	"PInfo-server/api"
	"PInfo-server/config"
	"PInfo-server/log"

	"context"
	"encoding/base64"
	"fmt"
	"time"
)

func (s *Service) UploadAvatar(ctx context.Context, req *api.UploadReq) (err error, rsp *api.CommRsp) {
	rsp = &api.CommRsp{
		Code:    0,
		Message: "ok",
		Data:    nil,
	}

	uploadRsp := &api.UploadRsp{}
	uploadOK := false
	bucket := config.AppConfig().CosInfo.AvatarBucket
	expireHour := config.AppConfig().CosInfo.Expire
	for _, fileHeaders := range req.Form.File {
		for _, file := range fileHeaders {
			inFile, err := file.Open()
			if err != nil {
				log.Errorf("open infile failed, path:%s, err:%v", file.Filename, err)
				continue
			}
			defer inFile.Close()
			key := fmt.Sprintf("avatar/%d/%s", req.Uid, file.Filename)
			err = s.dao.UploadFile(ctx, bucket, key, inFile, "", "")
			if err != nil {
				log.Errorf("UploadFile failed, path:%s, err:%v", key, err)
				uploadOK = false
				break
			} else {
				log.Infof("UploadFile ok, name:%s, size:%d", file.Filename, file.Size)
				uploadOK = true
			}

			if p, e := s.dao.GetPresignUrl(ctx, bucket, key, time.Duration(expireHour)); e == nil {
				uploadRsp.Avatar = p
				log.Infof("get presign ok, key:%s, url:%s", key, uploadRsp.Avatar)
			} else {
				log.Errorf("get presign fail, key:%s, err: %v", key, e)
			}
			break
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

func (s *Service) DownloadAvatar(ctx context.Context, req *api.DownloadReq) (err error, rsp *api.CommRsp) {
	rsp = &api.CommRsp{
		Code:    0,
		Message: "ok",
		Data:    nil,
	}

	data, err := s.dao.RawDownload(ctx, req.Url)
	if err != nil {
		log.Errorf("RawDownload err:%v", err)
		return err, rsp
	}

	rsp.Data = &api.DownloadRsp{
		Data: base64.StdEncoding.EncodeToString(data),
	}

	log.Infof("DownloadAvatar ok")
	return nil, rsp
}
