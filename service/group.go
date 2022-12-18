package service

import (
	"PInfo-server/api"
	"PInfo-server/model"
	"context"
	"github.com/spf13/cast"
	"log"
	"strings"
	"time"
)

func (s *Service) GetGroupMembers(ctx context.Context, req *api.GroupMembersReq) (err error, rsp *api.CommRsp) {
	// 如果ID为0，则获取我的联系人全集
	rsp = &api.CommRsp{
		Code:    0,
		Message: "",
		Data:    nil,
	}
	inviteGroupRsp := make([]*api.GroupMemberInfo, 0)
	if req.GroupId == 0 {
		err, contacts := s.dao.GetContactList(ctx, req.Uid, 2)
		if err != nil {
			return err, rsp
		}

		for _, con := range contacts {
			mem := &api.GroupMemberInfo{
				GroupId:      con.ContactID,
				Nickname:     con.Nickname,
				Gender:       con.Gender,
				Motto:        con.Motto,
				Avatar:       con.Avatar,
				FriendRemark: con.FriendRemark,
				IsOnline:     1,
			}
			inviteGroupRsp = append(inviteGroupRsp, mem)
			log.Printf("==== group mem:%+v\n", mem)
		}

		rsp.Data = inviteGroupRsp
		return nil, rsp

	}
	// 否则，获取该群内成员列表

	return nil, nil
}

func (s *Service) CreateGroup(ctx context.Context, req *api.CreateGroupReq) (err error, rsp *api.CommRsp) {
	rsp = &api.CommRsp{
		Code:    0,
		Message: "",
		Data:    nil,
	}

	if req.GroupName == "" {
		req.GroupName = "group " + req.UserName
	}

	createGroupRsp := &api.CreateGroupRsp{}
	createGroupRsp.GroupId = s.dao.GenGroupID()

	// 群信息写入db
	nowTime := time.Now().Unix()
	groupInfo := &model.Groups{
		GroupID:       createGroupRsp.GroupId,
		GroupName:     req.GroupName,
		GroupStatus:   1,
		GroupAvatar:   "",
		GroupTag:      "",
		GroupAnnounce: "",
		Sequence:      nowTime,
		CreateTime:    nowTime,
		UpdateTime:    nowTime,
	}
	if err = s.dao.SetGroupInfo(ctx, groupInfo); err != nil {
		rsp.Code = 5000
		rsp.Message = "set group info failed"
		return err, rsp
	}

	// 群成员写入db
	groupMembers := make([]*model.GroupMembers, 0)
	idList := strings.Split(req.Ids, ",")
	for _, id := range idList {
		uid := cast.ToInt64(id)
		groupMember := &model.GroupMembers{
			GroupID:    createGroupRsp.GroupId,
			Uid:        uid,
			UserRole:   1,
			RemarkName: "",
			Sequence:   nowTime,
			CreateTime: nowTime,
			UpdateTime: nowTime,
		}
		groupMembers = append(groupMembers, groupMember)
	}

	if err = s.dao.BatchAddGroupMember(ctx, groupMembers); err != nil {
		rsp.Code = 5000
		rsp.Message = "set group member failed"
		return err, rsp
	}

	// 会话
	con := &model.Conversations{
		Uid:                req.Uid,
		ContactID:          createGroupRsp.GroupId, // 群ID
		ConversationType:   2,
		ConversationName:   groupInfo.GroupName,
		ConversationStatus: 1,
		Unread:             1, // 更新时基于+1操作
		MsgDigest:          "",
		Sequence:           nowTime,
		CreateTime:         nowTime,
		UpdateTime:         nowTime,
	}

	if err := s.dao.UpdateConversationGroupMsg(ctx, []*model.Conversations{con}); err != nil {
		log.Printf("save msg for Conversation failed! group id:%d, group name:%s\n", createGroupRsp.GroupId, groupInfo.GroupName)
		return err, nil
	}

	rsp.Data = createGroupRsp

	log.Printf("[INFO] create group ok, req:%+v, rsp:%+v\n", req, rsp)
	return nil, rsp
}
