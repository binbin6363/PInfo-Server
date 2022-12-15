package service

import (
	"PInfo-server/api"
	"context"
	"log"
	"time"
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
		log.Printf("search user by user name failed, err:%v, user:%s\n", err, req.UserName)
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
			IsOnline:   0,
			IsRobot:    0,
			Name:       con.UserName,
			Avatar:     con.Avatar,
			RemarkName: con.RemarkName,
			UnreadNum:  con.Unread,
			MsgText:    con.MsgDigest,
			UpdatedAt:  time.Unix(con.UpdateTime, 0).Format("2006-01-02 15:04:05"),
		}
		listRsp.TalkList = append(listRsp.TalkList, talkInfo)
	}

	rsp.Data = listRsp

	log.Printf("[INFO] get conversation req:%+v, rsp:%+v", req, rsp)
	return nil, rsp
}
