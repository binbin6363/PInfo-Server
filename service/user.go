package service

import (
	"PInfo-server/api"
	"PInfo-server/log"
	"PInfo-server/model"
	"PInfo-server/utils"
	"context"
	"errors"
	"time"
)

func (s *Service) GetUserInfo(ctx context.Context, username string) (error, *model.UserInfo) {
	err, userInfo := s.dao.GetUserInfoByUserName(ctx, username)
	if err != nil {
		log.Infof("query user info failed, err:%v", err)
		return err, nil
	}

	log.Infof("ok get user info:%+v", userInfo)
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
		log.Infof("query user info failed, err:%v", err)
		return err
	}

	log.Infof("ok modify user info:%+v", userInfo)
	return nil
}

func (s *Service) RegisterUser(ctx context.Context, uid int64, req *api.RegisterReq) (err error, rsp *api.CommRsp) {
	rsp = &api.CommRsp{
		Code:    0,
		Message: "ok",
		Data:    nil,
	}

	// 1. 参数校验
	if req.UserName == "" || req.Password == "" {
		rsp.Code = 400
		rsp.Message = "invalid param"
		return errors.New("invalid param"), rsp
	}

	// 2. 校验用户名的存在性
	err, exist := s.dao.CheckUserExist(ctx, req.UserName)
	if exist {
		rsp.Code = 4000
		rsp.Message = "用户名已被占用"
		return errors.New("用户名已被占用"), rsp
	}

	err, passHash := utils.EncryptPassword(req.Password)
	if err != nil {
		rsp.Code = 4001
		rsp.Message = "内部错误"
		log.Infof("user:%s pass hash failed, pass:%s", req.UserName, req.Password)
		return errors.New("内部错误"), rsp
	}
	log.Infof("user:%s pass hash:%s", req.UserName, passHash)

	userInfo := &model.UserInfo{
		UserName:   req.UserName,
		PassHash:   passHash,
		NickName:   req.NickName,
		Phone:      "",
		Email:      "",
		Avatar:     "",
		Gender:     1,
		UserTag:    "",
		Motto:      "",
		CreateTime: time.Now().Unix(),
		UpdateTime: time.Now().Unix(),
	}

	// 申请用户ID
	err, userInfo.Uid = s.dao.AllocNewUserID(ctx)
	if err != nil {
		log.Infof("alloc user id failed, err:%v, user info:%+v", err, userInfo)
		rsp.Code = 4002
		rsp.Message = "内部错误"
		return errors.New("内部错误"), rsp
	}

	err = s.dao.SetUserInfo(ctx, userInfo)
	if err != nil {
		log.Infof("register user failed, err:%v, user info:%+v", err, userInfo)
		rsp.Code = 4002
		rsp.Message = "内部错误"
		return errors.New("内部错误"), rsp
	}

	log.Infof("register user ok, user:%s, nick:%s, pass hash:%s",
		req.UserName, req.NickName, userInfo.PassHash)
	return nil, rsp
}
