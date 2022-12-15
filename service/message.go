package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"PInfo-server/api"
	"PInfo-server/config"
	"PInfo-server/model"
)

// AddOneMessage .
func (s *Service) AddOneMessage(ctx context.Context, req *api.SendTextMsgReq) (error, *api.SendTextMsgRsp) {
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
	if err := s.dao.AddOneMessage(ctx, msg); err != nil {
		log.Printf("save msg for sender failed! info:%+v\n", msg)
		return err, nil
	}

	con := &model.Conversations{
		ID:                 0,
		Uid:                msg.Uid,
		ContactID:          req.ReceiverId,
		ConversationType:   1,
		ConversationName:   "",
		ConversationStatus: 1,
		Unread:             0,
		MsgDigest:          msg.Content,
		Sequence:           msg.MsgID,
		CreateTime:         msg.CreateTime,
		UpdateTime:         msg.UpdateTime,
	}
	if err := s.dao.UpdateConversationMsg(ctx, con); err != nil {
		log.Printf("save msg for Conversation failed! info:%+v\n", con)
		return err, nil
	}

	// 给接收者插入消息
	msg.Uid = req.ReceiverId
	msg.ID = 0
	if err := s.dao.AddOneMessage(ctx, msg); err != nil {
		// 此步骤只许成功，不能失败。失败要进离线队列
		log.Printf("save msg for receiver failed! info:%+v\n", msg)
		//return err
	}

	con.ID = 0
	con.Uid = req.ReceiverId
	con.ContactID = req.Uid
	if err := s.dao.UpdateConversationMsg(ctx, con); err != nil {
		log.Printf("save msg for Conversation failed! info:%+v\n", con)
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
				CreatedAt:  time.Unix(msg.CreateTime, 0).Format("2006-01-02 15:04:05"),
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
		log.Printf("notify conn failed, err:%+v\n", err)
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Printf("notify conn success, req:%+v, rsp:%s\n", req, string(body))
	}

	// 更新会话列表

	return nil, rsp
}

// QueryMessage 从大往小拉取历史消息
func (s *Service) QueryMessage(ctx context.Context, req *api.MsgRecordsReq) (err error, rsp *api.MsgRecordsRsp) {
	rsp = &api.MsgRecordsRsp{}
	rsp.Limit = req.Limit
	if req.TalkType == 1 {
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
				CreatedAt:  time.Unix(msgList[idx].CreateTime, 0).Format("2006-01-02 15:04:05"),
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
	}

	return nil, rsp
}
