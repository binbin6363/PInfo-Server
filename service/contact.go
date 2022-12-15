package service

import (
	"PInfo-server/model"
	"context"
	"log"
	"time"

	"PInfo-server/api"
)

func (s *Service) GetContactList(ctx context.Context, uid int64, status int) (error, []*api.ContactInfo) {
	err, userContacts := s.dao.GetContactList(ctx, uid, status)
	if err != nil {
		log.Printf("query contact info failed, err:%v\n", err)
		return err, nil
	}

	var contactList []*api.ContactInfo = nil
	for _, contact := range userContacts {
		contactInfo := &api.ContactInfo{
			Id:           contact.ContactID,
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
	err, contactInfo := s.dao.GetContactDetailInfo(ctx, req.Uid, req.ContactId)
	if err != nil {
		rsp.Code = 400
		rsp.Message = "内部错误"
		log.Printf("get user detail by uid failed, err:%v, user:%d\n", err, req.Uid)
		return err, rsp
	}

	applyStatus := 1
	if contactInfo.Status == 2 {
		applyStatus = 0
	}
	rsp.Data = api.ContactDetailRsp{
		Gender:         contactInfo.Gender,
		Avatar:         contactInfo.Avatar,
		Motto:          contactInfo.Motto,
		NickName:       contactInfo.Nickname,
		Uid:            contactInfo.ContactID,
		FriendStatus:   contactInfo.Status,
		FriendApply:    applyStatus,
		NickNameRemark: contactInfo.FriendRemark,
	}

	log.Printf("[INFO] search user ok, username:%d, rsp:%+v\n", req.Uid, rsp)
	return nil, rsp
}

func (s *Service) AddContact(ctx context.Context, req *api.AddContactReq) (err error, rsp *api.CommRsp) {
	rsp = &api.CommRsp{
		Code:    0,
		Message: "ok",
		Data:    nil,
	}

	// 检查对方是不是我的好友，正向好友关系
	isMyContact := false
	err, forwardInfo := s.dao.GetContactDetailInfo(ctx, req.Uid, req.ContactID)
	if err != nil {
		rsp.Code = 400
		rsp.Message = "内部错误"
		log.Printf("get forward contact failed, err:%v, %d => %d\n", err, req.Uid, req.ContactID)
		return err, rsp
	}
	if forwardInfo.ContactID != 0 && forwardInfo.Status != 0 {
		log.Printf("[forward] [%d] is [%d] friend, status:%d", req.ContactID, req.Uid, forwardInfo.Status)
		isMyContact = true
	}

	// 检查我是不是对方的好友，反向好友关系
	isHisContact := false
	err, reverseInfo := s.dao.GetContactDetailInfo(ctx, req.ContactID, req.Uid)
	if err != nil {
		rsp.Code = 400
		rsp.Message = "内部错误"
		log.Printf("get reverse contact failed, err:%v, %d => %d\n", err, req.ContactID, req.Uid)
		return err, rsp
	}
	if reverseInfo.ContactID != 0 && reverseInfo.Status != 0 {
		log.Printf("[reverse] [%d] is [%d] friend, status:%d", req.Uid, req.ContactID, reverseInfo.Status)
		isHisContact = true
	}

	if isMyContact && isHisContact {
		log.Printf("already each other contact friend, uid:%d, contact id:%d, status:%d", req.Uid, req.ContactID, reverseInfo.Status)
		return nil, rsp
	}
	nowTimme := time.Now().Unix()
	// 将对方添加到我的好友关系中，设置为已是好友status=2
	if !isMyContact {
		contactInfo := &model.Contacts{
			Uid:        req.Uid,
			ContactID:  req.ContactID,
			Status:     2,
			RemarkName: req.RemarkName,
			CreateTime: nowTimme,
			UpdateTime: nowTimme,
			Sequence:   nowTimme,
		}

		err = s.dao.SetContactInfo(ctx, contactInfo)
		if err != nil {
			rsp.Code = 400
			rsp.Message = "内部错误"
			log.Printf("SetContactInfo mine failed, err:%v, user:%d\n", err, req.ContactID)
			return err, rsp
		}
	}

	// 将我加到对方的好友关系中，设置为待通过status=1
	if !isHisContact {
		contactInfo := &model.Contacts{
			Uid:        req.ContactID,
			ContactID:  req.Uid,
			Status:     1,
			RemarkName: req.RemarkName,
			CreateTime: nowTimme,
			UpdateTime: nowTimme,
			Sequence:   nowTimme, // 序列号，保证递增即可
		}

		err = s.dao.SetContactInfo(ctx, contactInfo)
		if err != nil {
			rsp.Code = 400
			rsp.Message = "内部错误"
			log.Printf("SetContactInfo peer failed, err:%v, user:%d\n", err, req.ContactID)
			return err, rsp
		}
	}

	log.Printf("[INFO] add each other contact ok, req:%+v, rsp:%+v\n", req, rsp)
	return nil, rsp
}

func (s *Service) ApplyAddContact(ctx context.Context, req *api.ApplyAddContactReq) (err error, rsp *api.CommRsp) {
	rsp = &api.CommRsp{
		Code:    0,
		Message: "ok",
		Data:    nil,
	}
	// 检查对方是不是我的好友，正向关系
	err, forwardInfo := s.dao.GetContactDetailInfo(ctx, req.Uid, req.ApplyId)
	if err != nil {
		rsp.Code = 400
		rsp.Message = "内部错误"
		log.Printf("get my contact by contact failed, err:%v, %d => %d\n", err, req.Uid, req.ApplyId)
		return err, rsp
	}
	if forwardInfo.ContactID == 0 || forwardInfo.Status == 2 {
		log.Printf("contact[%d] is my[%d] friend already, status:%d", req.ApplyId, req.Uid, forwardInfo.Status)
		return nil, rsp
	}

	// 设置我的联系人好友状态
	nowTimme := time.Now().Unix()
	contactInfo := &model.Contacts{
		Uid:        req.Uid,
		ContactID:  req.ApplyId,
		Status:     2,
		RemarkName: req.Remark,
		UpdateTime: nowTimme,
		Sequence:   nowTimme, // 序列号，保证递增即可
	}
	err = s.dao.SetContactInfo(ctx, contactInfo)
	if err != nil {
		rsp.Code = 400
		rsp.Message = "内部错误"
		log.Printf("SetContactInfo mine failed, err:%v, user:%d\n", err, req.ApplyId)
		return err, rsp
	}

	log.Printf("[INFO] apply contact ok, req:%+v, rsp:%+v\n", req, rsp)
	return nil, rsp
}

// EditContactInfo 编辑好友备注信息
func (s *Service) EditContactInfo(ctx context.Context, req *api.EditContactInfoReq) (err error, rsp *api.CommRsp) {
	rsp = &api.CommRsp{
		Code:    0,
		Message: "ok",
		Data:    nil,
	}
	// 先获取联系人属性信息
	err, contactInfo := s.dao.GetContactInfo(ctx, req.Uid, req.ContactId)
	if err != nil {
		rsp.Code = 400
		rsp.Message = "内部错误"
		log.Printf("get my contact by contact failed, err:%v, %d => %d\n", err, req.Uid, req.ContactId)
		return err, rsp
	}
	if contactInfo.Uid == 0 || contactInfo.ContactID == 0 {
		log.Printf("contact not exist, uid:%d, contact:%d", req.Uid, req.ContactId)
		rsp.Code = 4005
		rsp.Message = "contact not exist"
		return nil, rsp
	}

	contactInfo.RemarkName = req.RemarkName
	err = s.dao.SetContactInfo(ctx, contactInfo)
	if err != nil {
		rsp.Code = 400
		rsp.Message = "内部错误"
		log.Printf("EditContactInfo failed, err:%v, contact:%d\n", err, req.ContactId)
		return err, rsp
	}

	log.Printf("[INFO] edit contact info ok, req:%+v, rsp:%+v\n", req, rsp)
	return nil, rsp
}
