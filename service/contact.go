package service

import (
	"context"
	"log"

	"PInfo-server/api"
)

func (s *Service) GetContactList(ctx context.Context, uid int64) (error, []*api.ContactInfo) {
	err, userContacts := s.dao.GetContactList(ctx, uid)
	if err != nil {
		log.Printf("query contact info failed, err:%v\n", err)
		return err, nil
	}

	var contactList []*api.ContactInfo = nil
	for _, contact := range userContacts {
		contactInfo := &api.ContactInfo{
			Id:           contact.Uid,
			Nickname:     contact.Nickname,
			Gender:       contact.Gender,
			Motto:        contact.Motto,
			Avatar:       contact.Avatar,
			FriendRemark: contact.FriendRemark,
			IsOnline:     0,
		}
		contactList = append(contactList, contactInfo)
	}

	return nil, contactList
}

func (s *Service) ContactSearch(ctx context.Context, req *api.ContactSearchReq) (err error, rsp *api.CommRsp) {
	rsp = &api.CommRsp{
		Code:    0,
		Message: "ok",
		Data:    nil,
	}
	err, info := s.dao.GetUserInfoByUserName(ctx, req.UserName)
	if err != nil {
		rsp.Code = 400
		rsp.Message = "内部错误"
		log.Printf("search user by user name failed, err:%v, user:%s\n", err, req.UserName)
		return err, rsp
	}

	rsp.Data = api.ContactSearchRsp{
		Uid: info.Uid,
	}

	log.Printf("[INFO] search user ok, username:%s, rsp:%+v\n", req.UserName, rsp)
	return nil, rsp
}

func (s *Service) ContactDetail(ctx context.Context, req *api.ContactDetailReq) (err error, rsp *api.CommRsp) {
	rsp = &api.CommRsp{
		Code:    0,
		Message: "ok",
		Data:    nil,
	}
	err, info := s.dao.GetUserInfoByUid(ctx, req.Uid)
	if err != nil {
		rsp.Code = 400
		rsp.Message = "内部错误"
		log.Printf("get user detail by uid failed, err:%v, user:%d\n", err, req.Uid)
		return err, rsp
	}

	rsp.Data = api.ContactDetailRsp{
		Gender:         info.Gender,
		Email:          info.Email,
		Avatar:         info.Avatar,
		UserName:       info.UserName,
		Motto:          info.Motto,
		NickName:       info.NickName,
		Uid:            info.Uid,
		FriendStatus:   1,
		FriendApply:    0,
		NickNameRemark: info.NickName,
	}

	log.Printf("[INFO] search user ok, username:%d, rsp:%+v\n", req.Uid, rsp)
	return nil, rsp
}
