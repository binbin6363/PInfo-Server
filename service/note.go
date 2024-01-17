package service

import (
	"PInfo-server/api"
	"PInfo-server/log"
	"PInfo-server/utils"
	"context"
	"time"
)

func (s *Service) ClassList(ctx context.Context, req *api.ClassListReq) (*api.ClassListRsp, error) {
	rsp := &api.ClassListRsp{}

	rsp.ClassItems = append(rsp.ClassItems, api.ClassItem{
		Id:        1,
		ClassName: "默认分组",
		Count:     0,
		IsDefault: true,
		UpdatedAt: utils.FormatTimeStr(time.Now().Unix()),
	})
	log.Infof("done ClassList, rsp: %v", rsp)
	return rsp, nil
}
