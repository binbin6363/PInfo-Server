package service

import (
	"PInfo-server/api"
	"PInfo-server/log"
	"PInfo-server/model"
	"context"
	"errors"
	"github.com/spf13/cast"
	"strings"
	"time"
)

// GetGroupList 获取群列表
func (s *Service) GetGroupList(ctx context.Context, req *api.GroupListReq) (err error, rsp *api.CommRsp) {
	// 如果ID为0，则获取我的联系人全集
	rsp = &api.CommRsp{
		Code:    0,
		Message: "",
		Data:    nil,
	}

	groupListRsp := &api.GroupListRsp{}
	err, groupList := s.dao.GetGroupList(ctx, req.Uid)
	if err != nil {
		log.Errorf("GetGroupList failed, err:%v", err)
		return err, rsp
	}

	for _, group := range groupList {
		groupRspInfo := &api.GroupInfo{
			Id:        group.GroupID,
			GroupName: group.GroupName,
			Avatar:    group.GroupAvatar,
			Profile:   group.GroupAnnounce, // 	群简介
			Leader:    group.Leader,
			IsDisturb: 0, // 免打扰
		}
		groupListRsp.GroupInfoList = append(groupListRsp.GroupInfoList, groupRspInfo)
		log.Infof("==== group info:%+v", groupRspInfo)
	}

	rsp.Data = groupListRsp
	return nil, rsp

}

func (s *Service) InviteGroupMember(ctx context.Context, req *api.GroupMembersReq) (err error, rsp *api.CommRsp) {
	rsp = &api.CommRsp{
		Code:    0,
		Message: "",
		Data:    nil,
	}
	inviteGroupRsp := make([]*api.GroupMemberInfo, 0)
	err, contacts := s.dao.GetContactList(ctx, req.Uid, int(model.ContactFriend))
	if err != nil {
		return err, rsp
	}

	for _, con := range contacts {
		mem := &api.GroupMemberInfo{
			Id:           con.ContactID,
			Nickname:     con.Nickname,
			Gender:       con.Gender,
			Motto:        con.Motto,
			Avatar:       con.Avatar,
			FriendRemark: con.FriendRemark,
			IsOnline:     1,
		}
		inviteGroupRsp = append(inviteGroupRsp, mem)
		log.Infof("==== group mem:%+v", mem)
	}

	rsp.Data = inviteGroupRsp
	return nil, rsp
}

func (s *Service) GetGroupMembers(ctx context.Context, req *api.GroupMembersReq) (err error, rsp *api.CommRsp) {
	// 如果ID为0，则获取我的联系人全集
	rsp = &api.CommRsp{
		Code:    0,
		Message: "",
		Data:    nil,
	}
	inviteGroupRsp := make([]*api.GroupMemberInfo, 0)
	err, memberList := s.dao.GetGroupMemberInfoList(ctx, req.GroupId, 0)
	if err != nil {
		return err, nil
	}

	inGroup := false
	for _, member := range memberList {
		if member.Uid == req.Uid {
			inGroup = true
		}
		mem := &api.GroupMemberInfo{
			Id:           member.Uid,
			Uid:          member.Uid,
			Nickname:     member.NickName,
			Gender:       member.Gender,
			Motto:        member.Motto,
			Avatar:       member.Avatar,
			FriendRemark: member.FriendRemark,
			IsOnline:     1,
		}
		if mem.FriendRemark == "" {
			mem.FriendRemark = mem.Nickname
		}
		inviteGroupRsp = append(inviteGroupRsp, mem)
		log.Infof("==== group mem:%+v", mem)

	}

	if !inGroup {
		return errors.New("user not in group"), nil
	}

	log.Infof("get group members ok, group id:%d, member size:%d", req.GroupId, len(inviteGroupRsp))
	rsp.Data = inviteGroupRsp
	return nil, rsp
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
	// 先加自己
	groupMembers = append(groupMembers, &model.GroupMembers{
		GroupID:    createGroupRsp.GroupId,
		Uid:        req.Uid,
		UserRole:   2, // 创建者是管理员
		RemarkName: "",
		Sequence:   nowTime,
		CreateTime: nowTime,
		UpdateTime: nowTime,
	})
	log.Infof("create group, admin:%v", groupMembers[0])

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
		log.Infof("create group, add group user:%v", groupMember)
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
		log.Infof("save msg for Conversation failed! group id:%d, group name:%s", createGroupRsp.GroupId, groupInfo.GroupName)
		return err, nil
	}

	rsp.Data = createGroupRsp

	log.Infof("create group ok, req:%+v, rsp:%+v", req, createGroupRsp)
	return nil, rsp
}

func (s *Service) InviteGroup(ctx context.Context, req *api.InviteGroupReq) (err error, rsp *api.CommRsp) {

	// 群成员写入db
	groupMembers := make([]*model.GroupMembers, 0)
	nowTime := time.Now().Unix()
	idList := strings.Split(req.Ids, ",")
	for _, id := range idList {
		uid := cast.ToInt64(id)
		groupMember := &model.GroupMembers{
			GroupID:    req.GroupId,
			Uid:        uid,
			UserRole:   1,
			RemarkName: "",
			Sequence:   nowTime,
			CreateTime: nowTime,
			UpdateTime: nowTime,
		}
		groupMembers = append(groupMembers, groupMember)
		log.Infof("invite group, add group user:%v", groupMember)
	}

	if err = s.dao.BatchAddGroupMember(ctx, groupMembers); err != nil {
		rsp.Code = 5000
		rsp.Message = "set group member failed"
		return err, rsp
	}

	log.Infof("InviteGroup ok, req:%+v, rsp:%+v", req, rsp)
	return nil, rsp
}

func (s *Service) RemarkNameInGroup(ctx context.Context, req *api.RemarkNameInGroupReq) (err error, rsp *api.CommRsp) {
	rsp = &api.CommRsp{
		Code:    0,
		Message: "",
		Data:    nil,
	}

	if req.RemarkName == "" {
		req.RemarkName = req.UserName
	}

	// 获取我在内的信息
	err, memberInfo := s.dao.GetGroupMemberInfo(ctx, req.GroupId, req.Uid)
	if err != nil {
		return err, nil
	}

	memberInfo.RemarkName = req.RemarkName
	memberInfo.Sequence = time.Now().Unix()
	memberInfo.UpdateTime = memberInfo.Sequence

	if err = s.dao.SetGroupMemberInfo(ctx, memberInfo); err != nil {
		rsp.Code = 5000
		rsp.Message = "RemarkNameInGroup failed"
		return err, rsp
	}

	log.Infof("RemarkNameInGroup ok, req:%+v, rsp:%+v", req, rsp)
	return nil, rsp
}

func (s *Service) GetGroupDetail(ctx context.Context, req *api.GroupDetailReq) (err error, rsp *api.CommRsp) {
	rsp = &api.CommRsp{
		Code:    0,
		Message: "",
		Data:    nil,
	}

	// 获取我在内的信息
	err, groupDetailInfo := s.dao.GetGroupDetailInfo(ctx, req.GroupId, req.Uid)
	if err != nil {
		rsp.Code = 5000
		rsp.Message = "GetGroupDetailInfo failed"
		return err, rsp
	}

	groupDetailRsp := &api.GroupDetailRsp{
		Avatar:          groupDetailInfo.GroupAvatar,
		CreatedAt:       groupDetailInfo.CreatedAt,
		GroupId:         groupDetailInfo.GroupId,
		GroupName:       groupDetailInfo.GroupName,
		IsDisturb:       groupDetailInfo.IsDisturb,
		IsManager:       groupDetailInfo.IsManager,
		ManagerNickname: groupDetailInfo.ManagerName,
		Notice:          groupDetailInfo.Notice,
		Profile:         "",
		VisitCard:       groupDetailInfo.VisitCard,
	}
	rsp.Data = groupDetailRsp

	log.Infof("GetGroupDetail ok, req:%+v, rsp:%+v", req, groupDetailRsp)
	return nil, rsp
}

func (s *Service) SetGroupInfo(ctx context.Context, req *api.SetGroupInfoReq) (err error, rsp *api.CommRsp) {
	rsp = &api.CommRsp{
		Code:    0,
		Message: "",
		Data:    nil,
	}

	// 获取我在内的信息
	err, groupInfo := s.dao.GetGroupInfo(ctx, req.GroupId)
	if err != nil {
		return err, nil
	}

	needSetCon := false
	nowTime := time.Now().Unix()
	if groupInfo.GroupAvatar != req.GroupAvatar && req.GroupAvatar != "" {
		groupInfo.GroupAvatar = req.GroupAvatar
	}
	if groupInfo.GroupAnnounce != req.GroupProfile && req.GroupProfile != "" {
		groupInfo.GroupAnnounce = req.GroupProfile
	}
	if groupInfo.GroupName != req.GroupName && req.GroupName != "" {
		groupInfo.GroupName = req.GroupName
		// 需要更新会话名
		needSetCon = true
	}
	groupInfo.GroupAnnounce = req.GroupProfile
	groupInfo.GroupName = req.GroupName
	groupInfo.Sequence = nowTime
	groupInfo.UpdateTime = nowTime

	if err = s.dao.SetGroupInfo(ctx, groupInfo); err != nil {
		rsp.Code = 5000
		rsp.Message = "SetGroupInfo failed"
		return err, rsp
	}

	// 修改群名，同步修改会话名
	if needSetCon {
		conversationInfo := &model.Conversations{
			ContactID:        groupInfo.GroupID,
			ConversationName: groupInfo.GroupName,
			Sequence:         nowTime,
			UpdateTime:       nowTime,
		}
		if err := s.dao.BatchSetGroupConversationName(ctx, conversationInfo); err != nil {
			log.Error("set group conversation name failed, group id:%d, err:%v", groupInfo.GroupID, err)
		}
	}

	log.Infof("SetGroupInfo ok, req:%+v, rsp:%+v", req, rsp)
	return nil, rsp
}
