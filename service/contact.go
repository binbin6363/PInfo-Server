package service

import (
	"PInfo-server/api"
	"PInfo-server/log"
	"PInfo-server/model"
	"PInfo-server/utils"
	"context"
	"time"

	"gorm.io/gorm"
)

func (s *Service) GetContactList(ctx context.Context, uid int64, status model.ContactType) (error, []*api.ContactInfo) {
	err, userContacts := s.dao.GetContactList(ctx, uid, int(status))
	if err != nil {
		log.InfoContextf(ctx, "query contact info failed, err:%v", err)
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
			CreatedAt:    utils.FormatTimeStr(contact.CreateTime),
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
		log.InfoContextf(ctx, "search user by user name failed, err:%v, user:%s", err, req.UserName)
		return err, rsp
	}

	rsp.Data = api.ContactSearchRsp{
		Uid: info.Uid,
	}

	log.InfoContextf(ctx, "search user ok, username:%s, rsp:%+v", req.UserName, rsp)
	return nil, rsp
}

// ContactDetail 获取用户详情
// 如果不在我的联系人记录里（包含发送好友申请），则返回非好友状态。
// 如果在我的联系人记录里，则判断状态，1代表已发送好友申请，2代表对方已通过，目前是好友关系。
// status由friendStatus+applyStatus组成。可选取值为：1【非好友，未申请】，2【等我审批】，3【好友，已发申请】，4【好友，已通过申请】
// friendStatus，是否是好友。1不是好友，2是好友。
// applyStatus，加好友申请状态。0未申请【匹配friendStatus=1】，1申请中【匹配friendStatus=2】，2申请完成【匹配friendStatus=2。
func (s *Service) ContactDetail(ctx context.Context, req *api.ContactDetailReq) (err error, rsp *api.CommRsp) {
	rsp = &api.CommRsp{
		Code:    0,
		Message: "ok",
		Data:    nil,
	}

	err, userInfo := s.dao.GetUserInfoByUid(ctx, req.ContactId)
	if err != nil {
		rsp.Code = 400
		rsp.Message = "内部错误"
		log.InfoContextf(ctx, "get user detail by uid failed, err:%v, user:%d", err, req.Uid)
		return err, rsp
	}

	friendStatus := 1 // 是否是好友。1不是好友，2是好友
	applyStatus := 0  // 加好友申请状态。0未申请，1申请中，2申请完成
	friendRemark := userInfo.NickName
	err, contactInfo := s.dao.GetContactDetailInfo(ctx, req.Uid, req.ContactId)
	if err == gorm.ErrRecordNotFound {
		friendStatus = 1
		applyStatus = 0
	} else if err == nil {
		if contactInfo.Status == 1 { // 标记删除，删除了就是status=1
			friendStatus = 1
			applyStatus = 0
		} else if contactInfo.Status == 2 { // 对方发起的好友申请等待我处理中
			friendStatus = 1
			applyStatus = 2
		} else if contactInfo.Status == 3 { // 我发起的好友申请等待对方处理中
			friendStatus = 2
			applyStatus = 1
		} else if contactInfo.Status == 4 { // 已通过好友申请，目前是好友状态
			friendStatus = 2
			applyStatus = 2
		}
	}
	//if err == nil {
	//	friendStatus = contactInfo.Status
	//	friendRemark = contactInfo.FriendRemark
	//	applyStatus = 1
	//	// 已经是好友了，申请状态标记为2
	//	if friendStatus == 2 {
	//		applyStatus = 2
	//	}
	//}
	//
	//// 查看对方是否已同意我的请求，如果对方没同意，修复为申请中状态
	//err, peerInfo := s.dao.GetContactDetailInfo(ctx, req.ContactId, req.Uid)
	//if err == nil {
	//	if peerInfo.Status == 1 {
	//		applyStatus = 1
	//	}
	//}
	//if err == gorm.ErrRecordNotFound {
	//	applyStatus &= 1
	//} else if err == nil && peerInfo != nil {
	//	applyStatus = peerInfo.Status
	//}

	rsp.Data = api.ContactDetailRsp{
		Gender:         userInfo.Gender,
		Avatar:         userInfo.Avatar,
		Motto:          userInfo.Motto,
		NickName:       userInfo.UserName,
		Uid:            req.ContactId,
		FriendStatus:   friendStatus,
		FriendApply:    applyStatus,
		NickNameRemark: friendRemark,
	}

	log.InfoContextf(ctx, "search user ok, username:%d, rsp:%+v", req.Uid, rsp)
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
	if err != nil && err != gorm.ErrRecordNotFound {
		rsp.Code = 400
		rsp.Message = "内部错误"
		log.InfoContextf(ctx, "get forward contact failed, err:%v, %d => %d", err, req.Uid, req.ContactID)
		return err, rsp
	}
	if forwardInfo != nil && forwardInfo.ContactID != 0 && forwardInfo.Status != 0 {
		log.InfoContextf(ctx, "[forward] [%d] is [%d] friend, status:%d", req.ContactID, req.Uid, forwardInfo.Status)
		isMyContact = true
	}

	// 检查我是不是对方的好友，反向好友关系
	isHisContact := false
	err, reverseInfo := s.dao.GetContactDetailInfo(ctx, req.ContactID, req.Uid)
	if err != nil && err != gorm.ErrRecordNotFound {
		rsp.Code = 400
		rsp.Message = "内部错误"
		log.InfoContextf(ctx, "get reverse contact failed, err:%v, %d => %d", err, req.ContactID, req.Uid)
		return err, rsp
	}
	if reverseInfo != nil && reverseInfo.ContactID != 0 && reverseInfo.Status != 0 {
		log.InfoContextf(ctx, "[reverse] [%d] is [%d] friend, status:%d", req.Uid, req.ContactID, reverseInfo.Status)
		isHisContact = true
	}

	if isMyContact && isHisContact {
		log.InfoContextf(ctx, "already each other contact friend, uid:%d, contact id:%d, status:%d", req.Uid, req.ContactID, reverseInfo.Status)
		return nil, rsp
	}
	nowTimme := time.Now().Unix()
	// 将对方添加到我的好友关系中，设置为已是好友status=4
	if !isMyContact {
		contactInfo := &model.Contacts{
			Uid:        req.Uid,
			ContactID:  req.ContactID,
			Status:     int(model.ContactFriend),
			RemarkName: req.RemarkName,
			CreateTime: nowTimme,
			UpdateTime: nowTimme,
			Sequence:   nowTimme,
		}

		err = s.dao.SetContactInfo(ctx, contactInfo)
		if err != nil {
			rsp.Code = 400
			rsp.Message = "内部错误"
			log.InfoContextf(ctx, "SetContactInfo mine failed, err:%v, user:%d", err, req.ContactID)
			return err, rsp
		}
	}

	// 将我加到对方的好友关系中，设置为待通过status=1
	if !isHisContact {
		contactInfo := &model.Contacts{
			Uid:        req.ContactID,
			ContactID:  req.Uid,
			Status:     int(model.ContactWaitMeApply),
			RemarkName: req.RemarkName,
			CreateTime: nowTimme,
			UpdateTime: nowTimme,
			Sequence:   nowTimme, // 序列号，保证递增即可
		}

		err = s.dao.SetContactInfo(ctx, contactInfo)
		if err != nil {
			rsp.Code = 400
			rsp.Message = "内部错误"
			log.InfoContextf(ctx, "SetContactInfo peer failed, err:%v, user:%d", err, req.ContactID)
			return err, rsp
		}
	}

	log.InfoContextf(ctx, "add each other contact ok, req:%+v, rsp:%+v", req, rsp)
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
		log.InfoContextf(ctx, "get my contact by contact failed, err:%v, %d => %d", err, req.Uid, req.ApplyId)
		return err, rsp
	}
	if forwardInfo.Status == int(model.ContactFriend) {
		log.InfoContextf(ctx, "contact[%d] is my[%d] friend already, status:%d", req.ApplyId, req.Uid, forwardInfo.Status)
		return nil, rsp
	}

	// 设置我的联系人好友状态
	nowTimme := time.Now().Unix()
	contactInfo := &model.Contacts{
		Uid:        req.Uid,
		ContactID:  req.ApplyId,
		Status:     int(model.ContactFriend),
		RemarkName: req.Remark,
		UpdateTime: nowTimme,
		Sequence:   nowTimme, // 序列号，保证递增即可
	}
	err = s.dao.SetContactInfo(ctx, contactInfo)
	if err != nil {
		rsp.Code = 400
		rsp.Message = "内部错误"
		log.InfoContextf(ctx, "SetContactInfo mine failed, err:%v, user:%d", err, req.ApplyId)
		return err, rsp
	}

	log.InfoContextf(ctx, "apply contact ok, req:%+v, rsp:%+v", req, rsp)
	return nil, rsp
}

func (s *Service) DeclineAddContact(ctx context.Context, req *api.ApplyAddContactReq) (err error, rsp *api.CommRsp) {
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
		log.InfoContextf(ctx, "get my contact by contact failed, err:%v, %d => %d", err, req.Uid, req.ApplyId)
		return err, rsp
	}
	if forwardInfo.Status == int(model.ContactStranger) {
		log.InfoContextf(ctx, "contact[%d] is decline by me[%d] already, status:%d", req.ApplyId, req.Uid, forwardInfo.Status)
		return nil, rsp
	}

	// 设置我的联系人好友状态
	nowTimme := time.Now().Unix()
	contactInfo := &model.Contacts{
		Uid:        req.Uid,
		ContactID:  req.ApplyId,
		Status:     int(model.ContactStranger),
		RemarkName: req.Remark,
		UpdateTime: nowTimme,
		Sequence:   nowTimme, // 序列号，保证递增即可
	}
	err = s.dao.SetContactInfo(ctx, contactInfo)
	if err != nil {
		rsp.Code = 400
		rsp.Message = "内部错误"
		log.InfoContextf(ctx, "SetContactInfo mine failed, err:%v, user:%d", err, req.ApplyId)
		return err, rsp
	}

	log.InfoContextf(ctx, "decline contact ok, req:%+v, rsp:%+v", req, rsp)
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
		log.InfoContextf(ctx, "get my contact by contact failed, err:%v, %d => %d", err, req.Uid, req.ContactId)
		return err, rsp
	}
	if contactInfo.Uid == 0 || contactInfo.ContactID == 0 {
		log.InfoContextf(ctx, "contact not exist, uid:%d, contact:%d", req.Uid, req.ContactId)
		rsp.Code = 4005
		rsp.Message = "contact not exist"
		return nil, rsp
	}

	contactInfo.RemarkName = req.RemarkName
	err = s.dao.SetContactInfo(ctx, contactInfo)
	if err != nil {
		rsp.Code = 400
		rsp.Message = "内部错误"
		log.InfoContextf(ctx, "EditContactInfo failed, err:%v, contact:%d", err, req.ContactId)
		return err, rsp
	}

	log.InfoContextf(ctx, "edit contact info ok, req:%+v, rsp:%+v", req, rsp)
	return nil, rsp
}
