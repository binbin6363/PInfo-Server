package dao

import (
	"context"
	"errors"

	"PInfo-server/log"
	"PInfo-server/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// GenGroupID 生成群ID。复用雪花
func (d *Dao) GenGroupID() int64 {
	return d.sf.NextVal()
}

func (d *Dao) SetGroupInfo(ctx context.Context, groupInfo *model.Groups) error {
	r := d.db(ctx)
	if groupInfo.GroupID == 0 {
		log.Error("group id invalid")
		return errors.New("group id invalid")
	}

	r = r.Clauses(clause.OnConflict{
		// key列
		Columns: []clause.Column{{Name: "group_id"}},
		// 需要更新的列
		DoUpdates: clause.AssignmentColumns([]string{"group_status", "group_name",
			"group_avatar", "group_tag", "group_announce", "sequence", "update_time"}),
	}).Create(groupInfo)

	if err := r.Error; err != nil {
		log.Infof("SetGroupInfo update db error(%v) user info:%+v", err, groupInfo)
		return err
	}

	log.Infof("SetGroupInfo update db ok user info:%+v", groupInfo)
	return nil
}

func (d *Dao) GetGroupInfo(ctx context.Context, GroupId int64) (error, *model.Groups) {
	r := d.db(ctx)
	if GroupId == 0 {
		log.Error("contact id is invalid")
		return errors.New("contact id is invalid"), nil
	}

	groupInfo := &model.Groups{}
	if err := r.Where("group_id=?", GroupId).First(&groupInfo).Error; err != nil {
		log.Infof("GetGroupInfo read db error(%v) group id(%d)", err, GroupId)
		return err, nil
	}

	log.Infof("GetGroupInfo read db ok group id(%d), info(%+v)", GroupId, groupInfo)
	return nil, groupInfo
}

func (d *Dao) BatchAddGroupMember(ctx context.Context, groupMembers []*model.GroupMembers) error {
	r := d.db(ctx)
	if len(groupMembers) == 0 {
		log.Error("group member empty")
		return errors.New("group member empty")
	}

	r = r.Clauses(clause.OnConflict{
		// key列
		Columns: []clause.Column{{Name: "group_id"}, {Name: "uid"}},
		// 需要更新的列
		DoUpdates: clause.AssignmentColumns([]string{"user_role", "remark_name",
			"sequence", "update_time"}),
	}).Create(groupMembers)

	if err := r.Error; err != nil {
		log.Infof("BatchAddGroupMember update db error(%v) user info:%+v", err, groupMembers)
		return err
	}

	log.Infof("BatchAddGroupMember update db ok, add user size:%d", len(groupMembers))
	return nil
}

func (d *Dao) GetGroupMemberList(ctx context.Context, groupId int64) (error, []*model.GroupMembers) {
	r := d.db(ctx)
	if groupId == 0 {
		log.Error("group id is invalid")
		return errors.New("group id is invalid"), nil
	}

	groupMembers := make([]*model.GroupMembers, 0)
	if err := r.Table(model.GroupMembers{}.TableName()).Where("group_id=?", groupId).Scan(&groupMembers).Error; err != nil {
		log.Infof("GetGroupMemberList read db error(%v) groupId(%d)", err, groupId)
		return err, nil
	}

	log.Infof("GetGroupInfo read db ok groupId(%d), members size:%d", groupId, len(groupMembers))
	return nil, groupMembers
}

func (d *Dao) GetGroupMemberInfoList(ctx context.Context, groupId, sequence int64) (error, []*model.GroupMemberInfoList) {
	r := d.db(ctx)
	if groupId == 0 {
		log.Error("group id is invalid")
		return errors.New("group id is invalid"), nil
	}

	r = r.Table(model.GroupMembers{}.TableName()).Select("group_members.uid as uid, group_members.user_role as user_role, "+
		"user_infos.nickname as nickname, user_infos.gender as gender, user_infos.motto as motto,user_infos.avatar as avatar,  "+
		"group_members.remark_name as remark_name, group_members.sequence as sequence, group_members.create_time as create_time").
		Joins("left join user_infos on group_members.uid=user_infos.uid where group_members.group_id=? and sequence>?",
			groupId, sequence)

	groupUserInfos := make([]*model.GroupMemberInfoList, 0)
	if err := r.Scan(&groupUserInfos).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Infof("record not exist, group id:%d", groupId)
		} else {
			log.Errorf("GetGroupMemberInfoList read db error(%v) group id(%d)", err, groupId)
		}
		return err, nil
	}

	log.Infof("GetGroupMemberInfoList read db ok uid(%d)", groupId)
	return nil, groupUserInfos
}

func (d *Dao) GetGroupMemberInfo(ctx context.Context, groupId, uid int64) (error, *model.GroupMembers) {
	r := d.db(ctx)
	if groupId == 0 {
		log.Error("group id is invalid")
		return errors.New("group id is invalid"), nil
	}

	r = r.Where("group_id=? and uid=?", groupId, uid)

	groupUserInfo := &model.GroupMembers{}
	if err := r.First(&groupUserInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Infof("record not exist, group id:%d, uid:%d", groupId, uid)
		} else {
			log.Errorf("GetGroupMemberInfo read db error(%v) group id(%d)", err, groupId)
		}
		return err, nil
	}

	log.Infof("GetGroupMemberInfo read db ok uid(%d)", groupId)
	return nil, groupUserInfo
}

func (d *Dao) SetGroupMemberInfo(ctx context.Context, groupUserInfo *model.GroupMembers) error {
	r := d.db(ctx)
	if groupUserInfo.GroupID == 0 || groupUserInfo.Uid == 0 {
		log.Error("group id or uid invalid")
		return errors.New("group id or uid invalid")
	}

	r = r.Clauses(clause.OnConflict{
		// key列
		Columns: []clause.Column{{Name: "id"}, {Name: "group_id"}, {Name: "uid"}},
		// 需要更新的列
		DoUpdates: clause.AssignmentColumns([]string{"user_role", "remark_name", "sequence", "update_time"}),
	}).Create(groupUserInfo)

	if err := r.Error; err != nil {
		log.Infof("SetGroupMemberInfo update db error(%v) user info:%+v", err, groupUserInfo)
		return err
	}

	log.Infof("SetGroupMemberInfo update db ok user info:%+v", groupUserInfo)
	return nil

}

func (d *Dao) GetGroupList(ctx context.Context, uid int64) (err error, groupList []*model.GroupInfoList) {
	r := d.db(ctx)
	if uid == 0 {
		log.Error("uid is invalid")
		return errors.New("uid is invalid"), nil
	}
	type Result struct {
		GroupId int64 `gorm:"column:group_id"`
	}
	res := make([]*Result, 0)

	// 拉取用户的群ID列表
	err = r.Table("group_members").Select([]string{"group_id"}).
		Where("uid=?", uid).Scan(&res).Error
	if err == gorm.ErrRecordNotFound {
		log.Infof("user not exist")
		return err, nil
	}

	// 从群信息中获取
	var groupIds []int64
	for idx := range res {
		groupIds = append(groupIds, res[idx].GroupId)
	}
	groupList = make([]*model.GroupInfoList, 0)
	err = r.Table("groups").Where("group_id in ?", groupIds).Scan(&groupList).Error
	if err != nil {
		log.Infof("query group list failed, err:%+v", err)
		return err, nil
	}

	type Leader struct {
		GroupId int64 `gorm:"column:group_id"`
		Uid     int64 `gorm:"column:uid"`
	}
	leaders := make([]*Leader, 0)
	err = r.Table("group_members").Select([]string{"group_id", "uid"}).
		Where("group_id in ? and user_role=2", groupIds).Scan(&leaders).Error
	if err == gorm.ErrRecordNotFound {
		log.Infof("user not exist")
		//return err, nil
	}

	groupLeader := make(map[int64]int64)
	for _, ld := range leaders {
		groupLeader[ld.GroupId] = ld.Uid
	}

	for idx := range groupList {
		groupList[idx].Leader = groupLeader[groupList[idx].GroupID]
	}

	return nil, groupList
}

func (d *Dao) GetGroupDetailInfo(ctx context.Context, groupId, uid int64) (error, *model.GroupDetailInfo) {
	r := d.db(ctx)
	if groupId == 0 {
		log.Error("group id is invalid")
		return errors.New("group id is invalid"), nil
	}

	groupDetailInfo := &model.GroupDetailInfo{}
	err := r.Table(model.Groups{}.TableName()).Select("groups.group_id as group_id, groups.group_name as group_name, "+
		"groups.group_avatar as group_avatar, groups.group_announce as group_announce, groups.create_time as create_time, "+
		"group_members.remark_name as remark_name, group_members.user_role as user_role, group_members.disturb as disturb").
		Joins("left join group_members on group_members.group_id=groups.group_id where group_members.group_id=? and group_members.uid=?",
			groupId, uid).Scan(&groupDetailInfo).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Infof("record not exist, group id:%d", groupId)
		} else {
			log.Errorf("GetGroupDetailInfo read db error(%v) group id(%d)", err, groupId)
		}
		return err, nil
	}

	type Leader struct {
		GroupId    int64  `gorm:"column:group_id"`
		Uid        int64  `gorm:"column:uid"`
		RemarkName string `gorm:"column:remark_name"`
	}
	leaders := &Leader{}
	err = r.Table(model.GroupMembers{}.TableName()).Select([]string{"group_id", "uid", "remark_name"}).
		Where("group_id = ? and user_role=2", groupId).First(&leaders).Error
	if err == gorm.ErrRecordNotFound {
		log.Infof("not found manager in group id:%d", groupId)
	} else if err == nil {
		groupDetailInfo.ManagerName = leaders.RemarkName
		if leaders.Uid == uid {
			groupDetailInfo.IsManager = true
		} else {
			groupDetailInfo.IsManager = false
		}
	}

	log.Infof("GetGroupDetailInfo read db ok group id(%d), uid(%d)", groupId, uid)
	return nil, groupDetailInfo
}
