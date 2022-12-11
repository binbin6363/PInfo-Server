package service

import (
	"PInfo-server/api"
	"PInfo-server/model"
	"context"
	"log"
)

func (s *Service) GetUserInfo(ctx context.Context, username string) (error, *model.UserInfo) {
	err, userInfo := s.dao.GetUserInfoByUserName(ctx, username)
	if err != nil {
		log.Printf("query user info failed, err:%v\n", err)
		return err, nil
	}

	log.Printf("ok get user info:%+v", userInfo)
	return nil, userInfo
}

func (s *Service) SetUserInfo(ctx context.Context, uid int64, req *api.ModifyUsersSettingReq) error {
	err, userInfo := s.dao.GetUserInfoByUid(ctx, uid)
	if err != nil {
		return err
	}
	userInfo.NickName = req.NickName
	userInfo.Avatar = req.Avatar
	userInfo.Gender = req.Gender
	userInfo.Motto = req.Motto
	err = s.dao.SetUserInfo(ctx, userInfo)
	if err != nil {
		log.Printf("query user info failed, err:%v\n", err)
		return err
	}

	log.Printf("ok modify user info:%+v", userInfo)
	return nil
}
