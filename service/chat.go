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

func (s *Service) GetConversationList(ctx context.Context, req *api.TalkListReq) (err error, rsp *api.CommRsp) {

	rsp = &api.CommRsp{
		Code:    0,
		Message: "ok",
		Data:    nil,
	}
	err, conversationList := s.dao.GetConversationList(ctx, req.Uid, req.Sequence)
	if err != nil {
		rsp.Code = 400
		rsp.Message = "内部错误"
		log.InfoContextf(ctx, "search user by user name failed, err:%v, user:%s", err, req.UserName)
		return err, rsp
	}

	listRsp := &api.TalkListRsp{}
	for _, con := range conversationList {
		talkInfo := &api.TalkInfo{
			ID:         con.ID,
			Type:       con.ConversationType,
			ReceiverId: con.ContactID,
			IsTop:      0,
			IsDisturb:  0,
			IsOnline:   1,
			IsRobot:    0,
			Name:       con.UserName,
			Avatar:     con.Avatar,
			RemarkName: con.RemarkName,
			UnreadNum:  con.Unread,
			MsgText:    con.MsgDigest,
			UpdatedAt:  utils.FormatTimeStr(con.UpdateTime),
		}
		if talkInfo.RemarkName == "" {
			talkInfo.RemarkName = con.ConversationName
		}
		listRsp.TalkList = append(listRsp.TalkList, talkInfo)
	}

	rsp.Data = listRsp

	log.InfoContextf(ctx, "get conversation req:%+v, rsp:%+v", req, rsp)
	return nil, rsp
}

func (s *Service) CreateConversation(ctx context.Context, req *api.CreateTalkReq) (err error, rsp *api.CommRsp) {

	rsp = &api.CommRsp{
		Code:    0,
		Message: "ok",
		Data:    nil,
	}
	err, conversationInfo := s.dao.GetConversation(ctx, req.Uid, req.ContactId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.InfoContextf(ctx, "conversation not exist, need create")
		} else {
			rsp.Code = 400
			rsp.Message = "内部错误"
			log.InfoContextf(ctx, "search user by user name failed, err:%v, user:%s", err, req.UserName)
			return err, rsp
		}
	}

	// 如果已存在，则直接返回信息
	if conversationInfo != nil && conversationInfo.Uid != 0 && conversationInfo.ContactID != 0 {
		err, userInfo := s.dao.GetUserInfoByUid(ctx, req.ContactId)
		if err != nil {
			rsp.Code = 400
			rsp.Message = "联系人获取失败"
			log.InfoContextf(ctx, "get contact info failed, err:%v, %d => %d", err, req.Uid, req.ContactId)
			return err, rsp
		}
		conversationAvatar := userInfo.Avatar
		rsp.Data = &api.TalkInfo{
			ID:         conversationInfo.ID,
			Type:       conversationInfo.ConversationType,
			ReceiverId: conversationInfo.ContactID,
			IsTop:      0,
			IsDisturb:  0,
			IsOnline:   0,
			IsRobot:    0,
			Name:       conversationInfo.ConversationName,
			Avatar:     conversationAvatar,
			RemarkName: conversationInfo.ConversationName,
			UnreadNum:  0,
			MsgText:    conversationInfo.MsgDigest,
			UpdatedAt:  utils.FormatTimeStr(conversationInfo.UpdateTime),
			CreatedAt:  utils.FormatTimeStr(conversationInfo.CreateTime),
		}

		return nil, rsp
	}

	// 不存在，则新建
	conversationName := "新建群聊天"
	conversationAvatar := ""
	if req.TalkType == 1 {
		// 获取联系人备注信息
		err, contactInfo := s.dao.GetContactInfo(ctx, req.Uid, req.ContactId)
		if err != nil {
			rsp.Code = 400
			rsp.Message = "联系人获取失败"
			log.InfoContextf(ctx, "get contact info failed, err:%v, %d => %d", err, req.Uid, req.ContactId)
			return err, rsp
		}
		conversationName = contactInfo.RemarkName
		err, userInfo := s.dao.GetUserInfoByUid(ctx, req.ContactId)
		if err != nil {
			rsp.Code = 400
			rsp.Message = "联系人获取失败"
			log.InfoContextf(ctx, "get contact info failed, err:%v, %d => %d", err, req.Uid, req.ContactId)
			return err, rsp
		}
		conversationAvatar = userInfo.Avatar
	} else if req.TalkType == 2 {
		err, groupInfo := s.dao.GetGroupInfo(ctx, req.ContactId)
		if err != nil {
			rsp.Code = 400
			rsp.Message = "群信息获取失败"
			log.InfoContextf(ctx, "get group info failed, err:%v, %d => %d", err, req.Uid, req.ContactId)
			return err, rsp
		}
		conversationName = groupInfo.GroupName
		conversationAvatar = groupInfo.GroupAvatar
	}

	nowTime := time.Now().Unix()

	conversationInfo = &model.Conversations{
		Uid:                req.Uid,
		ContactID:          req.ContactId,
		ConversationType:   req.TalkType,
		ConversationName:   conversationName,
		ConversationStatus: 1,
		Unread:             0,
		MsgDigest:          "",
		Sequence:           nowTime,
		CreateTime:         nowTime,
		UpdateTime:         nowTime,
	}
	err = s.dao.SetConversation(ctx, conversationInfo)
	if err != nil {
		rsp.Code = 400
		rsp.Message = "会话创建失败"
		log.InfoContextf(ctx, "create conversation failed, err:%v, %d => %d", err, req.Uid, req.ContactId)
		return err, rsp
	}

	rsp.Data = &api.TalkInfo{
		ID:         conversationInfo.ID,
		Type:       req.TalkType,
		ReceiverId: req.ContactId,
		IsTop:      0,
		IsDisturb:  0,
		IsOnline:   0,
		IsRobot:    0,
		Name:       conversationName,
		Avatar:     conversationAvatar,
		RemarkName: conversationName,
		UnreadNum:  0,
		MsgText:    "",
		UpdatedAt:  utils.FormatTimeStr(nowTime),
		CreatedAt:  utils.FormatTimeStr(nowTime),
	}

	log.InfoContextf(ctx, "create conversation req:%+v, rsp:%+v", req, *rsp)
	return nil, rsp
}
