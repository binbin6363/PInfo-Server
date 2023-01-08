package service

import (
	"PInfo-server/utils"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"time"

	"PInfo-server/api"
	"PInfo-server/config"
	"PInfo-server/log"
	"PInfo-server/model"
)

func (s *Service) sendSingleTextMessage(ctx context.Context, req *api.SendTextMsgReq) (error, *api.SendTextMsgRsp) {
	err, myFriendInfo := s.dao.GetContactDetailInfo(ctx, req.Uid, req.ReceiverId)
	if err == nil && (myFriendInfo.Status != int(model.ContactFriend)) {
		log.Errorf("i(%d) am not your(%d) friend, send message failed", req.Uid, req.ReceiverId)
		return errors.New("我不是对方的好友"), nil
	}
	err, yourFriendInfo := s.dao.GetContactDetailInfo(ctx, req.ReceiverId, req.Uid)
	if err == nil && (yourFriendInfo.Status != int(model.ContactFriend)) {
		log.Errorf("peer(%d) is not your(%d) friend, send message failed", req.ReceiverId, req.Uid)
		return errors.New("对方不是你的好友"), nil
	}

	msg := &model.SingleMessages{
		Uid:         req.Uid,
		MsgID:       s.dao.GenMsgID(), // 本应该由中心服务生成，此处暂且放在本机生成
		ClientMsgID: req.ClientMsgId,
		SenderID:    req.Uid,
		ReceiverID:  req.ReceiverId,
		MsgType:     1,
		Content:     req.Text,
		MsgStatus:   0,
		CreateTime:  time.Now().Unix(),
		UpdateTime:  time.Now().Unix(),
	}

	// 给发送者插入消息
	msg.Uid = req.Uid
	if err := s.dao.AddOneSingleMessage(ctx, msg); err != nil {
		log.Infof("save msg for sender failed! info:%+v", msg)
		return err, nil
	}

	con := &model.Conversations{
		ID:                 0,
		Uid:                req.Uid,
		ContactID:          req.ReceiverId,
		ConversationType:   1,
		ConversationName:   myFriendInfo.FriendRemark,
		ConversationStatus: 1,
		Unread:             0,
		MsgDigest:          msg.Content,
		Sequence:           msg.MsgID,
		CreateTime:         msg.CreateTime,
		UpdateTime:         msg.UpdateTime,
	}
	if err := s.dao.UpdateConversationSingleMsg(ctx, con); err != nil {
		log.Infof("save msg for Conversation failed! info:%+v", con)
		return err, nil
	}

	msg.Uid = req.ReceiverId
	msg.ID = 0
	if err := s.dao.AddOneSingleMessage(ctx, msg); err != nil {
		// 此步骤只许成功，不能失败。失败要进离线队列
		log.Infof("save msg for receiver failed! info:%+v", msg)
		//return err
	}

	con.ID = 0
	con.ConversationName = yourFriendInfo.FriendRemark
	con.Uid = req.ReceiverId
	con.ContactID = req.Uid
	if err := s.dao.UpdateConversationSingleMsg(ctx, con); err != nil {
		log.Infof("save msg for Conversation failed! info:%+v", con)
		return err, nil
	}

	// 前端设计不合理，消息竟然需要携带用户信息
	_, userInfo := s.dao.GetUserInfoByUid(ctx, req.Uid)

	// Sequence 这个字段是否需要，尚待考虑。可以不用该字段
	rsp := &api.SendTextMsgRsp{
		Content: api.SendTextMsgContent{
			Data: api.SendTextMsgData{
				Id:         msg.MsgID,
				Sequence:   1,
				TalkType:   req.TalkType,
				MsgType:    1,
				UserId:     req.Uid,
				ReceiverId: req.ReceiverId,
				Content:    req.Text,
				Nickname:   userInfo.NickName,
				Avatar:     userInfo.Avatar,
				IsRevoke:   0,
				IsRead:     0,
				IsMark:     0,
				CreatedAt:  utils.FormatTimeStr(msg.CreateTime),
			},
			ReceiverId: req.ReceiverId,
			SenderId:   req.Uid,
			TalkType:   1,
		},
	}

	notice := &api.SendTextMsgEvtNotice{
		Event:   "event_talk",
		Content: rsp.Content,
	}
	// 消息通知对方，走下游websocket。下游决定推送的设备
	bytesData, _ := json.Marshal(notice)
	url := fmt.Sprintf("http://%s/notice/message/text", config.AppConfig().ConnInfo.Addr)
	resp, err := http.Post(url, "application/json; charset=utf-8", bytes.NewReader(bytesData))
	defer resp.Body.Close()
	if err != nil {
		log.Infof("notify conn failed, err:%+v", err)
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Infof("notify conn success, req:%+v, rsp:%s", req, string(body))
	}

	// 更新会话列表

	return nil, rsp
}

func (s *Service) sendGroupTextMessage(ctx context.Context, req *api.SendTextMsgReq) (error, *api.SendTextMsgRsp) {
	groupId := req.ReceiverId
	msg := &model.GroupMessages{
		GroupID:     groupId,
		MsgID:       s.dao.GenMsgID(), // 本应该由中心服务生成，此处暂且放在本机生成
		ClientMsgID: req.ClientMsgId,  // 客户端发送的消息做去重
		SenderID:    req.Uid,
		MsgType:     1,
		Content:     req.Text,
		MsgStatus:   0,
		CreateTime:  time.Now().Unix(),
		UpdateTime:  time.Now().Unix(),
	}

	// 群消息入库
	if err := s.dao.AddOneGroupMessage(ctx, msg); err != nil {
		log.Infof("save group msg for sender failed! info:%+v", msg)
		return err, nil
	}

	// 获取群头像
	err, groupInfo := s.dao.GetGroupInfo(ctx, groupId)
	if err == gorm.ErrRecordNotFound {
		log.Errorf("group not exist, group id:%d", groupId)
		return err, nil
	}

	// 获取群成员ID列表
	_, groupMembers := s.dao.GetGroupMemberList(ctx, groupId)
	conversationList := make([]*model.Conversations, 0, len(groupMembers))
	for _, groupMember := range groupMembers {
		con := &model.Conversations{
			Uid:                groupMember.Uid,
			ContactID:          groupId, // 群ID
			ConversationType:   2,
			ConversationName:   groupInfo.GroupName,
			ConversationStatus: 1,
			Unread:             1, // 更新时基于+1操作
			MsgDigest:          msg.Content,
			Sequence:           msg.MsgID,
			CreateTime:         msg.CreateTime,
			UpdateTime:         msg.UpdateTime,
		}
		conversationList = append(conversationList, con)
	}
	if err := s.dao.UpdateConversationGroupMsg(ctx, conversationList); err != nil {
		log.Infof("save msg for Conversation failed! group id:%d, group name:%s", groupId, groupInfo.GroupName)
		return err, nil
	}
	// 前端设计不合理，消息竟然需要携带用户信息
	_, userInfo := s.dao.GetUserInfoByUid(ctx, req.Uid)

	// Sequence 这个字段是否需要，尚待考虑。可以不用该字段
	rsp := &api.SendTextMsgRsp{
		Content: api.SendTextMsgContent{
			Data: api.SendTextMsgData{
				Id:         msg.MsgID,
				Sequence:   1,
				TalkType:   req.TalkType,
				MsgType:    1,
				UserId:     req.Uid,
				ReceiverId: groupId,
				Content:    req.Text,
				Nickname:   userInfo.NickName,
				Avatar:     userInfo.Avatar,
				IsRevoke:   0,
				IsRead:     0,
				IsMark:     0,
				CreatedAt:  utils.FormatTimeStr(msg.CreateTime),
			},
			ReceiverId: req.ReceiverId,
			SenderId:   req.Uid,
			TalkType:   2,
		},
	}

	notice := &api.SendTextMsgEvtNotice{
		Event:   "event_talk",
		Content: rsp.Content,
	}
	// 消息通知对方，走下游websocket。下游决定推送的设备
	bytesData, _ := json.Marshal(notice)
	url := fmt.Sprintf("http://%s/notice/message/text", config.AppConfig().ConnInfo.Addr)
	resp, err := http.Post(url, "application/json; charset=utf-8", bytes.NewReader(bytesData))
	defer resp.Body.Close()
	if err != nil {
		log.Infof("notify conn failed, err:%+v", err)
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Infof("notify conn success, req:%+v, rsp:%s", req, string(body))
	}

	// 更新会话列表

	return nil, rsp
}

// SendTextMessage .
func (s *Service) SendTextMessage(ctx context.Context, req *api.SendTextMsgReq) (error, *api.SendTextMsgRsp) {
	if req.TalkType == 1 {
		return s.sendSingleTextMessage(ctx, req)
	} else if req.TalkType == 2 {
		return s.sendGroupTextMessage(ctx, req)
	} else {
		log.Infof("unsupported talk type:%d", req.TalkType)
		return errors.New("unsupported talk type"), nil
	}
}

func (s *Service) querySingleMessage(ctx context.Context, req *api.MsgRecordsReq) (err error, rsp *api.MsgRecordsRsp) {
	rsp = &api.MsgRecordsRsp{}
	rsp.Limit = req.Limit
	err, msgList := s.dao.QuerySingleMessage(ctx, req.Uid, req.PeerId, req.MinMsgId, req.Limit)
	if err != nil {
		return err, nil
	}

	var minId int64 = 9223372036854775807

	_, selfInfo := s.dao.GetUserInfoByUid(ctx, req.Uid)
	_, peerInfo := s.dao.GetUserInfoByUid(ctx, req.PeerId)
	for idx := range msgList {
		msgRow := api.MessageRow{
			Id:         msgList[idx].MsgID,
			Sequence:   1,
			TalkType:   req.TalkType,
			MsgType:    msgList[idx].MsgType,
			UserId:     msgList[idx].SenderID,
			ReceiverId: msgList[idx].ReceiverID,
			IsRevoke:   msgList[idx].MsgStatus,
			IsMark:     1,
			IsRead:     1,
			Content:    msgList[idx].Content,
			CreatedAt:  utils.FormatTimeStr(msgList[idx].CreateTime),
		}
		if msgList[idx].SenderID == req.Uid {
			msgRow.Avatar = selfInfo.Avatar
			msgRow.Nickname = selfInfo.NickName
		} else if msgList[idx].SenderID == req.PeerId {
			msgRow.Avatar = peerInfo.Avatar
			msgRow.Nickname = peerInfo.NickName
		}
		rsp.Rows = append(rsp.Rows, msgRow)
		if msgRow.Id < minId {
			minId = msgRow.Id
		}
	}
	rsp.MaxRecordId = minId
	log.Infof("get single message, uid:%d, peer id:%d, min msgid:%d, size:%d",
		req.Uid, req.PeerId, minId, len(rsp.Rows))
	return err, rsp
}

func (s *Service) queryGroupMessage(ctx context.Context, req *api.MsgRecordsReq) (err error, rsp *api.MsgRecordsRsp) {
	rsp = &api.MsgRecordsRsp{}
	rsp.Limit = req.Limit
	err, msgList := s.dao.QueryGroupMessage(ctx, req.PeerId, req.MinMsgId, req.Limit)
	if err != nil {
		return err, nil
	}

	var minId int64 = 9223372036854775807

	infoMap := make(map[int64]*model.UserInfo, 0)
	for idx := range msgList {
		msgRow := api.MessageRow{
			Id:         msgList[idx].MsgID,
			Sequence:   1,
			TalkType:   req.TalkType,
			MsgType:    msgList[idx].MsgType,
			UserId:     msgList[idx].SenderID,
			ReceiverId: msgList[idx].GroupID,
			IsRevoke:   msgList[idx].MsgStatus,
			IsMark:     1,
			IsRead:     1,
			Content:    msgList[idx].Content,
			CreatedAt:  utils.FormatTimeStr(msgList[idx].CreateTime),
		}

		if _, ok := infoMap[msgList[idx].SenderID]; !ok {
			err, info := s.dao.GetUserInfoByUid(ctx, msgList[idx].SenderID)
			if err != nil {
				log.Errorf("get group user info failed, uid:%d, err:%+v", msgList[idx].SenderID, err)
				continue
			}
			infoMap[msgList[idx].SenderID] = info
		}

		info := infoMap[msgList[idx].SenderID]
		log.Infof("show uid:%d, info:%+v", msgList[idx].SenderID, info)
		msgRow.Avatar = info.Avatar
		msgRow.Nickname = info.NickName

		rsp.Rows = append(rsp.Rows, msgRow)
		if msgRow.Id < minId {
			minId = msgRow.Id
		}
	}
	rsp.MaxRecordId = minId
	log.Infof("get group message, uid:%d, group id:%d, min msgid:%d, size:%d",
		req.Uid, req.PeerId, minId, len(rsp.Rows))
	return err, rsp
}

// QueryMessage 从大往小拉取历史消息
func (s *Service) QueryMessage(ctx context.Context, req *api.MsgRecordsReq) (err error, rsp *api.MsgRecordsRsp) {

	if req.TalkType == 1 {
		return s.querySingleMessage(ctx, req)
	} else if req.TalkType == 2 {
		return s.queryGroupMessage(ctx, req)
	}

	return errors.New("unknown talk type"), nil
}
