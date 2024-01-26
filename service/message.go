package service

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"time"

	"gorm.io/gorm"

	"PInfo-server/api"
	"PInfo-server/config"
	"PInfo-server/log"
	"PInfo-server/model"
	"PInfo-server/utils"
)

func (s *Service) sendSingleTextMessage(ctx context.Context, req *api.SendTextMsgReq) (error, *api.SendTextMsgRsp) {
	err, myFriendInfo := s.dao.GetContactDetailInfo(ctx, req.Uid, req.ReceiverId)
	if err != nil {
		if gorm.ErrRecordNotFound == err {
			log.Errorf("对方不是你的好友")
			return errors.New("对方不是你的好友"), nil
		}
		log.Errorf("GetContactDetailInfo err: %v", err)
		return errors.New("服务器内部异常"), nil
	}
	if myFriendInfo.Status != int(model.ContactFriend) {
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
		MsgType:     model.MsgTypeText,
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
		Content: api.SendMsgContent{
			Data: api.SendMsgData{
				Id:         msg.MsgID,
				Sequence:   1,
				TalkType:   req.TalkType,
				MsgType:    model.MsgTypeText,
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
	resp, err := s.HttpPost(ctx, url, bytesData, 500)
	if err != nil {
		log.Infof("notify conn failed, err:%+v", err)
	} else {
		log.Infof("notify conn success, req:%+v, rsp:%s", req, string(resp))
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
		MsgType:     model.MsgTypeText,
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
		Content: api.SendMsgContent{
			Data: api.SendMsgData{
				Id:         msg.MsgID,
				Sequence:   1,
				TalkType:   req.TalkType,
				MsgType:    model.MsgTypeText,
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
	if resp != nil {
		defer resp.Body.Close()
	}
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

// makeFullUrl 根据存储中的key组装原始URL
func (s *Service) makeFullUrl(ctx context.Context, msgType int, key string) string {
	if len(key) == 0 {
		return key
	}
	d := config.AppConfig().CosInfo.Domain
	b := config.AppConfig().CosInfo.MediaBucket
	switch msgType {
	case model.MsgTypeImg, model.MsgTypeFile, model.MsgTypeAudio, model.MsgTypeVideo:
		b = config.AppConfig().CosInfo.MediaBucket
	default:
		log.Errorf("unknown msg type: %d for url", msgType)
		return key
	}
	u, e := s.dao.GetPresignUrl(ctx, b, key, time.Duration(config.AppConfig().CosInfo.Expire))
	if e == nil {
		return u
	}
	return d + key
}

func (s *Service) querySingleMessage(ctx context.Context, req *api.MsgRecordsReq) (err error, rsp *api.MsgRecordsRsp) {
	rsp = &api.MsgRecordsRsp{}
	rsp.Limit = req.Limit
	err, msgList := s.dao.QuerySingleMessage(ctx, req.Uid, req.PeerId, req.MinMsgId, req.Limit, req.MsgType)
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
		if msgList[idx].MediaInfo != "" &&
			(msgRow.MsgType == model.MsgTypeFile ||
				msgRow.MsgType == model.MsgTypeImg ||
				msgRow.MsgType == model.MsgTypeAudio ||
				msgRow.MsgType == model.MsgTypeVideo) {
			if e := json.Unmarshal([]byte(msgList[idx].MediaInfo), &msgRow.FileItem); e == nil {
				msgRow.FileItem.Url = s.makeFullUrl(ctx, msgRow.MsgType, msgRow.FileItem.Url)
			} else {
				log.Errorf("Unmarshal media info err: %v, media: %s", e, msgList[idx].MediaInfo)
				continue
			}
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
	err, msgList := s.dao.QueryGroupMessage(ctx, req.PeerId, req.MinMsgId, req.Limit, req.MsgType)
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

// makeImgKey 构建图片文件路径，格式：img/year/md5。返回key,md5
func (s *Service) makeImgKey(name string, reader io.Reader) (string, string) {
	suffix := filepath.Ext(name)
	hash := md5.New()
	_, _ = io.Copy(hash, reader)
	b := hash.Sum(nil)
	md5Str := base64.StdEncoding.EncodeToString(b)
	hexStr := hex.EncodeToString(b)
	key := fmt.Sprintf("img/%d/%s_512x512%s", time.Now().Year(), hexStr, suffix)
	log.Infof("upload image key: %s", key)
	return key, md5Str
}

func (s *Service) SendImageMessage(ctx context.Context, req *api.SendImageMsgReq) (rsp *api.SendImageMsgRsp, err error) {
	_, userInfo := s.dao.GetUserInfoByUid(ctx, req.Uid)
	rsp = &api.SendImageMsgRsp{
		Content: api.SendMsgContent{
			ReceiverId: req.ReceiverId,
			SenderId:   req.Uid,
			TalkType:   req.TalkType,
			Data: api.SendMsgData{
				Sequence:    1,
				TalkType:    req.TalkType,
				MsgType:     model.MsgTypeImg,
				UserId:      req.Uid,
				ReceiverId:  req.ReceiverId,
				Nickname:    userInfo.NickName,
				Avatar:      userInfo.Avatar,
				IsRevoke:    0,
				IsRead:      0,
				IsMark:      0,
				CreatedAt:   utils.FormatTimeStr(time.Now().Unix()),
				FileContent: &api.FileContent{},
			},
		},
	}

	msg := &model.SingleMessages{
		Uid:         req.Uid,
		MsgID:       s.dao.GenMsgID(), // 本应该由中心服务生成，此处暂且放在本机生成
		ClientMsgID: req.ClientMsgId,
		SenderID:    req.Uid,
		ReceiverID:  req.ReceiverId,
		MsgType:     model.MsgTypeImg,
		MsgStatus:   0,
		CreateTime:  time.Now().Unix(),
		UpdateTime:  time.Now().Unix(),
	}
	// 上传图片
	bucket := config.AppConfig().CosInfo.MediaBucket
	expireHour := config.AppConfig().CosInfo.Expire
	for _, fileHeaders := range req.Form.File {
		for _, file := range fileHeaders {
			inFile, err := file.Open()
			if err != nil {
				log.Errorf("open infile failed, path:%s, err:%v", file.Filename, err)
				continue
			}
			defer inFile.Close()
			// 构建图片路径
			key, md5Str := s.makeImgKey(file.Filename, inFile)
			// 上传图片，先重置文件指针
			if _, err = inFile.Seek(0, 0); err != nil {
				log.Errorf("seek file err: %v", err)
				break
			}
			err = s.dao.UploadFile(ctx, bucket, key, inFile, "image/jpeg", md5Str)
			if err != nil {
				log.Errorf("UploadFile failed, path:%s, err:%v", key, err)
				break
			} else {
				log.Infof("UploadFile ok, name:%s, size:%d", file.Filename, file.Size)
			}

			// 上传成功，写db
			msg.Content = file.Filename
			fileItem := api.FileItem{
				Url: key,
			}
			if b, e := json.Marshal(fileItem); e == nil {
				msg.MediaInfo = string(b)
			} else {
				log.Errorf("Marshal media info err: %v", e)
			}
			// 给发送者插入消息
			msg.Uid = req.Uid
			if err = s.dao.AddOneSingleMessage(ctx, msg); err != nil {
				log.Infof("save img msg for sender failed! info:%+v", msg)
				break
			}
			rsp.Content.Data.Id = msg.ID
			rsp.Content.Data.Content = file.Filename
			// 给接收者插入消息
			msg.Uid = req.ReceiverId
			msg.ID = 0
			if err = s.dao.AddOneSingleMessage(ctx, msg); err != nil {
				log.Infof("save img msg for receiver failed! info:%+v", msg)
				break
			}

			if p, e := s.dao.GetPresignUrl(ctx, bucket, key, time.Duration(expireHour)); e == nil {
				log.Infof("get presign ok, key:%s, url:%s", key, p)
				rsp.Content.Data.FileContent.Name = file.Filename
				rsp.Content.Data.FileContent.Url = p
				// 仅支持单文件发送
				break
			} else {
				log.Errorf("get presign fail, key:%s, err: %v", key, e)
			}
			break
		}
	}

	return rsp, err
}
