package dao

import (
	"context"
	"errors"
	"gorm.io/gorm/clause"

	"PInfo-server/log"
	"PInfo-server/model"
	"PInfo-server/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Dao is Data Access Object
type Dao struct {
	commDB *gorm.DB
	sf     *utils.Snowflake
}

// New creates Dao instance
// dsn eg: "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
func New(dsn string, dataCenterId, WorkerId int64) *Dao {
	d := &Dao{}

	cli, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("dao: New db gorm client error(%v)", err)
	}
	d.commDB = cli

	s, err := utils.NewSnowflake(dataCenterId, WorkerId)
	if err != nil {
		log.Fatalf("dao: NewSnowflake error(%v), dataCenterId:%d, WorkerId:%d", err, dataCenterId, WorkerId)
	}
	log.Infof("dao: NewSnowflake dataCenterId:%d, WorkerId:%d", dataCenterId, WorkerId)
	d.sf = s

	return d
}

func (d *Dao) db(ctx context.Context) *gorm.DB {
	return d.commDB
}

// GenMsgID 生成消息ID。雪花算法，保证递增
func (d *Dao) GenMsgID() int64 {
	return d.sf.NextVal()
}

// GenGroupID 生成群ID。复用雪花
func (d *Dao) GenGroupID() int64 {
	return d.sf.NextVal()
}

// GetUserInfoByUserName 获取用户信息
func (d *Dao) GetUserInfoByUserName(ctx context.Context, username string) (error, *model.UserInfo) {
	r := d.db(ctx)
	if username == "" {
		log.Error("username is empty, invalid")
		return errors.New("username is invalid"), nil
	}

	userInfo := &model.UserInfo{}
	if err := r.Debug().Where("username=?", username).Limit(1).Find(&userInfo).Error; err != nil {
		log.Infof("GetUserInfoByUserName read db error(%v) username(%s)", err, username)
		return err, nil
	}

	log.Infof("GetUserInfoByUserName read db ok username(%s), info:%+v", username, userInfo)
	return nil, userInfo
}

// GetUserInfoByUid 获取用户信息
func (d *Dao) GetUserInfoByUid(ctx context.Context, uid int64) (error, *model.UserInfo) {
	r := d.db(ctx)
	if uid == 0 {
		log.Error("uid is invalid")
		return errors.New("uid is invalid"), nil
	}

	userInfo := &model.UserInfo{}
	if err := r.Debug().Where("uid=?", uid).Limit(1).Find(&userInfo).Error; err != nil {
		log.Infof("GetUserInfoByUid read db error(%v) uid(%d)", err, uid)
		return err, nil
	}

	log.Infof("GetUserInfoByUid read db ok uid(%d)", uid)
	return nil, userInfo
}

// SetUserInfo 设置用户信息
func (d *Dao) SetUserInfo(ctx context.Context, userInfo *model.UserInfo) error {
	r := d.db(ctx)
	if userInfo.Uid == 0 {
		log.Error("uid invalid")
		return errors.New("uid invalid")
	}

	r = r.Clauses(clause.OnConflict{
		// key列
		Columns: []clause.Column{{Name: "username"}},
		// 需要更新的列
		DoUpdates: clause.AssignmentColumns([]string{"passhash", "nickname",
			"phone", "email", "avatar", "gender", "user_tag", "motto", "update_time"}),
	}).Create(userInfo)

	if err := r.Error; err != nil {
		log.Infof("SetUserInfo update db error(%v) user info:%+v", err, userInfo)
		return err
	}

	log.Infof("SetUserInfo update db ok user info:%+v", userInfo)
	return nil
}

// GetContactList 获取我的好友列表信息
// SELECT contacts.uid as uid, contacts.contact_id as contact_id, user_infos.nickname as nickname, user_infos.gender as gender, user_infos.motto as motto, user_infos.avatar as avatar, contacts.remark_name as remark_name, contacts.status as status FROM `contacts` left join user_infos on contacts.uid=user_infos.uid where contacts.uid=10000\G;
func (d *Dao) GetContactList(ctx context.Context, uid int64, status int) (error, []*model.UserContact) {
	r := d.db(ctx)
	if uid == 0 {
		log.Error("uid is invalid")
		return errors.New("uid is invalid"), nil
	}

	r = r.Table(model.Contacts{}.TableName()).Select("contacts.uid as uid, contacts.contact_id as contact_id, "+
		"user_infos.nickname as nickname, user_infos.gender as gender, user_infos.motto as motto, "+
		"user_infos.avatar as avatar, contacts.remark_name as remark_name, contacts.status as status, contacts.create_time as create_time").
		Joins("left join user_infos on contacts.contact_id=user_infos.uid where contacts.uid=? and contacts.status=?",
			uid, status)

	userContacts := make([]*model.UserContact, 0)
	if err := r.Debug().Scan(&userContacts).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Infof("record not exist, uid:%d", uid)
		} else {
			log.Errorf("GetContactList read db error(%v) uid(%d)", err, uid)
		}
		return err, nil
	}

	log.Infof("GetContactList read db ok uid(%d)", uid)
	return nil, userContacts
}

// AddOneSingleMessage 添加单人消息
func (d *Dao) AddOneSingleMessage(ctx context.Context, msg *model.SingleMessages) error {
	r := d.db(ctx)
	if msg.Uid == 0 || msg.MsgID == 0 {
		log.Error("uid|msgid invalid")
		return errors.New("uid|msgid invalid")
	}

	if err := r.Debug().Create(msg).Error; err != nil {
		log.Infof("AddOneSingleMessage insert db error(%v) msg:%+v", err, msg)
		return err
	}

	log.Infof("AddOneSingleMessage insert db ok msg:%+v", msg)
	return nil
}

// AddOneGroupMessage 添加群消息
func (d *Dao) AddOneGroupMessage(ctx context.Context, msg *model.GroupMessages) error {
	r := d.db(ctx)
	if msg.GroupID == 0 || msg.MsgID == 0 {
		log.Error("group id|msgid invalid")
		return errors.New("group id|msgid invalid")
	}

	if err := r.Debug().Create(msg).Error; err != nil {
		log.Infof("AddOneGroupMessage insert db error(%v) msg:%+v", err, msg)
		return err
	}

	log.Infof("AddOneGroupMessage insert db ok msg:%+v", msg)
	return nil
}

// UpdateConversationSingleMsg 针对单人消息更新会话列表及未读数
func (d *Dao) UpdateConversationSingleMsg(ctx context.Context, con *model.Conversations) error {
	r := d.db(ctx)
	if con.Uid == 0 || con.ContactID == 0 {
		log.Error("uid invalid")
		return errors.New("uid invalid")
	}

	r = r.Clauses(clause.OnConflict{
		// key列
		Columns: []clause.Column{{Name: "uid"}, {Name: "contact_id"}},
		// 需要更新的列。页面上仅支持这四列的手动修改。其他列的修改，都应该直接走server_list.csv更新（通用的）
		DoUpdates: clause.AssignmentColumns([]string{"sequence", "msg_digest", "update_time"}),
	}).Debug().Create(con)

	// 未读数加1
	r.Debug().UpdateColumn("unread", gorm.Expr("unread + ?", 1))

	log.Infof("UpdateConversationSingleMsg update db ok conversations info:%+v", con)
	return nil
}

// UpdateConversationGroupMsg 针对群消息更新会话列表及未读数
func (d *Dao) UpdateConversationGroupMsg(ctx context.Context, con []*model.Conversations) error {
	r := d.db(ctx)
	if len(con) == 0 {
		log.Error("group id invalid")
		return errors.New("group id invalid")
	}

	// 更新该群下的所有用户的会话信息
	r.Clauses(clause.OnConflict{
		// key列
		Columns: []clause.Column{{Name: "contact_id"}},
		// 需要更新的列。页面上仅支持这四列的手动修改。其他列的修改，都应该直接走server_list.csv更新（通用的）
		DoUpdates: clause.AssignmentColumns([]string{"sequence", "msg_digest", "update_time"}),
	}).Debug().Create(con)

	// 群成员未读数加1
	r.Table(model.Conversations{}.TableName()).Debug().Where("contact_id=?", con[0].ContactID).UpdateColumn("unread", gorm.Expr("unread + ?", 1))

	log.Infof("UpdateConversationGroupMsg update db ok conversations info:%+v", con)
	return nil
}

// QuerySingleMessage 拉取单人历史消息
func (d *Dao) QuerySingleMessage(ctx context.Context, uid, peerUid, minMsgId int64, limit int) (error, []*model.SingleMessages) {
	r := d.db(ctx)
	if uid == 0 || peerUid == 0 {
		log.Error("uid|msgid invalid")
		return errors.New("uid|msgid invalid"), nil
	}

	// 如果id为0，从最大开始拉
	if minMsgId == 0 {
		r = r.Where("uid=? and (senderid=? or receiverid=?)", uid, peerUid, peerUid)
	} else {
		r = r.Where("uid=? and (senderid=? or receiverid=?) and msgid<?", uid, peerUid, peerUid, minMsgId)
	}
	// 分页，取第index页的count条数据。倒序
	r = r.Order(clause.OrderByColumn{Column: clause.Column{Name: "msgid"}, Desc: true})
	r = r.Limit(limit)

	msgList := make([]*model.SingleMessages, 0)
	if err := r.Debug().Find(&msgList).Error; err != nil {
		log.Infof("QuerySingleMessage read db error(%v)", err)
		return err, nil
	}

	log.Infof("QuerySingleMessage ok, msg size:%d", len(msgList))
	return nil, msgList
}

// QueryGroupMessage 拉取群历史消息
func (d *Dao) QueryGroupMessage(ctx context.Context, groupId, minMsgId int64, limit int) (error, []*model.GroupMessages) {
	r := d.db(ctx)
	if groupId == 0 {
		log.Error("groupId invalid")
		return errors.New("groupId invalid"), nil
	}

	// 如果id为0，从最大开始拉
	if minMsgId == 0 {
		r = r.Where("groupid=?", groupId)
	} else {
		r = r.Where("groupid=? and msgid<?", groupId, minMsgId)
	}
	// 分页，取第index页的count条数据。倒序
	r = r.Order(clause.OrderByColumn{Column: clause.Column{Name: "msgid"}, Desc: true})
	r = r.Limit(limit)

	msgList := make([]*model.GroupMessages, 0)
	if err := r.Debug().Find(&msgList).Error; err != nil {
		log.Infof("QueryGroupMessage read db error(%v)", err)
		return err, nil
	}

	log.Infof("QueryGroupMessage ok, msg size:%d", len(msgList))
	return nil, msgList
}

func (d *Dao) CheckUserExist(ctx context.Context, username string) (err error, exist bool) {
	r := d.db(ctx)
	type Result struct {
		UserName   string `gorm:"column:username"`
		CreateTime int64  `gorm:"column:create_time"`
	}

	res := &Result{}
	err = r.Debug().Table("user_infos").Select([]string{"username", "create_time"}).
		Where("username=?", username).First(&res).Error
	if err == gorm.ErrRecordNotFound {
		log.Infof("user not exist")
		return err, false
	}

	if err != nil {
		log.Infof("query user name failed, err:%+v", err)
		return err, false
	}

	if res.UserName == "" {
		log.Infof("user %s not exist", username)
		return nil, false
	}

	log.Infof("user exist, name:%s, create time:%d", res.UserName, res.CreateTime)
	return nil, true
}

// AllocNewUserID 获取新用户ID
func (d *Dao) AllocNewUserID(ctx context.Context) (err error, uid int64) {
	r := d.db(ctx)
	type Result struct {
		Uid int64 `gorm:"column:uid"`
	}

	res := &Result{}
	err = r.Debug().Table("user_infos").Order(clause.OrderByColumn{
		Column: clause.Column{
			Name: "uid",
		},
		Desc: true,
	}).Select([]string{"uid"}).Scan(&res).Error

	if err != nil {
		log.Infof("alloc user id failed, err:%+v", err)
		return err, 0
	}

	uid = res.Uid + 1
	log.Infof("alloc user id ok, id:%d", uid)
	return nil, uid
}

// GetContactDetailInfo 获取uid好友contactId的详细信息
func (d *Dao) GetContactDetailInfo(ctx context.Context, uid, contactId int64) (error, *model.UserContact) {
	r := d.db(ctx)
	if uid == 0 {
		log.Error("contact id is invalid")
		return errors.New("contact id is invalid"), nil
	}

	r = r.Table(model.Contacts{}.TableName()).Select("contacts.uid as uid, contacts.contact_id as contact_id, "+
		"user_infos.nickname as nickname, user_infos.gender as gender, user_infos.motto as motto, "+
		"user_infos.avatar as avatar, contacts.remark_name as remark_name, contacts.status as status").
		Joins("left join user_infos on contacts.contact_id=user_infos.uid where contacts.uid=? and contacts.contact_id=?",
			uid, contactId)

	userContacts := &model.UserContact{}
	if err := r.Debug().First(&userContacts).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Infof("record not exist, uid:%d", uid)
		} else {
			log.Errorf("GetContactList read db error(%v) uid(%d)", err, uid)
		}
		return err, nil
	}

	log.Infof("GetContactDetailInfo read db ok uid(%d)", uid)
	return nil, userContacts
}

// GetContactInfo 获取uid好友contactId的基础信息
func (d *Dao) GetContactInfo(ctx context.Context, uid, contactId int64) (error, *model.Contacts) {
	r := d.db(ctx)
	if uid == 0 || contactId == 0 {
		log.Error("contact id is invalid")
		return errors.New("contact id is invalid"), nil
	}

	contactInfo := &model.Contacts{}
	if err := r.Debug().Where("uid=? and contact_id=?", uid, contactId).First(&contactInfo).Error; err != nil {
		log.Infof("GetContactInfo read db error(%v) uid(%d)", err, uid)
		return err, nil
	}

	log.Infof("GetContactInfo read db ok uid(%d), contactInfo:%+v", uid, *contactInfo)
	return nil, contactInfo
}

func (d *Dao) SetContactInfo(ctx context.Context, contact *model.Contacts) error {
	r := d.db(ctx)
	if contact.Uid == 0 || contact.ContactID == 0 {
		log.Error("uid invalid")
		return errors.New("uid invalid")
	}

	r = r.Clauses(clause.OnConflict{
		// key列
		Columns: []clause.Column{{Name: "uid"}, {Name: "contact_id"}},
		// 需要更新的列
		DoUpdates: clause.AssignmentColumns([]string{"remark_name", "status",
			"sequence", "update_time"}),
	}).Create(contact)

	if err := r.Error; err != nil {
		log.Infof("SetContactInfo update db error(%v) user info:%+v", err, contact)
		return err
	}

	log.Infof("SetContactInfo update db ok user info:%+v", contact)
	return nil
}

// GetConversationList 差量获取会话列表
func (d *Dao) GetConversationList(ctx context.Context, uid, sequence int64) (error, []*model.ConversationDetails) {
	r := d.db(ctx)
	if uid == 0 {
		log.Error("uid is invalid")
		return errors.New("uid is invalid"), nil
	}

	// todo: 看起来是有bug。查到多余的数据
	r = r.Table(model.Conversations{}.TableName()).Select("conversations.id as id, conversations.uid as uid, "+
		"conversations.contact_id as contact_id, conversations.conversation_type as conversation_type, "+
		"conversations.unread as unread, conversations.msg_digest as msg_digest, conversations.sequence as sequence, "+
		"conversations.update_time as update_time, user_infos.username as username,conversations.conversation_name as conversation_name, "+
		"contacts.remark_name as remark_name, user_infos.motto as motto, user_infos.avatar as avatar").
		Joins("left join contacts on conversations.contact_id=contacts.contact_id and conversations.uid=contacts.uid").
		Joins("left join user_infos on conversations.contact_id=user_infos.uid where conversations.uid=? and conversations.sequence>?", uid, sequence)

	conDetails := make([]*model.ConversationDetails, 0)
	if err := r.Debug().Scan(&conDetails).Error; err != nil {
		log.Infof("GetConversationList read db error(%v) uid(%d)", err, uid)
		return err, nil
	}

	log.Infof("GetConversationList read db ok uid(%d), sequence(%d)", uid, sequence)
	return nil, conDetails
}

func (d *Dao) GetConversation(ctx context.Context, uid, contactId int64) (error, *model.Conversations) {
	r := d.db(ctx)
	if uid == 0 || contactId == 0 {
		log.Error("contact id is invalid")
		return errors.New("contact id is invalid"), nil
	}

	conversationInfo := &model.Conversations{}
	if err := r.Debug().Where("uid=? and contact_id=?", uid, contactId).First(&conversationInfo).Error; err != nil {
		log.Infof("GetConversation read db error(%v) uid(%d)", err, uid)
		return err, nil
	}

	log.Infof("GetConversation read db ok uid(%d), info(%+v)", uid, conversationInfo)
	return nil, conversationInfo
}

func (d *Dao) SetConversation(ctx context.Context, conversationInfo *model.Conversations) error {
	r := d.db(ctx)
	if conversationInfo.Uid == 0 || conversationInfo.ContactID == 0 {
		log.Error("uid invalid")
		return errors.New("uid invalid")
	}

	r = r.Clauses(clause.OnConflict{
		// key列
		Columns: []clause.Column{{Name: "uid"}, {Name: "contact_id"}, {Name: "conversation_type"}},
		// 需要更新的列
		DoUpdates: clause.AssignmentColumns([]string{"conversation_name", "conversation_status",
			"unread", "msg_digest", "sequence", "update_time"}),
	}).Create(conversationInfo)

	if err := r.Error; err != nil {
		log.Infof("SetConversation update db error(%v) user info:%+v", err, conversationInfo)
		return err
	}

	log.Infof("SetConversation update db ok user info:%+v", conversationInfo)
	return nil
}

func (d *Dao) BatchSetGroupConversationName(ctx context.Context, conversationInfo *model.Conversations) error {
	r := d.db(ctx)
	if conversationInfo.ContactID == 0 || conversationInfo.ConversationName == "" {
		log.Error("groupId|conversationName invalid")
		return errors.New("groupId|conversationName invalid")
	}

	r = r.Debug().Model(&model.Conversations{}).Where("contact_id = ?", conversationInfo.ContactID).
		UpdateColumn("conversation_name", conversationInfo.ConversationName)
	if err := r.Error; err != nil {
		log.Infof("BatchSetGroupConversationName update db error(%v) user info:%+v", err, conversationInfo)
		return err
	}

	log.Infof("BatchSetGroupConversationName update db ok user info:%+v", conversationInfo)
	return nil
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
	if err := r.Debug().Where("group_id=?", GroupId).First(&groupInfo).Error; err != nil {
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
	if err := r.Table(model.GroupMembers{}.TableName()).Debug().Where("group_id=?", groupId).Scan(&groupMembers).Error; err != nil {
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
	if err := r.Debug().Scan(&groupUserInfos).Error; err != nil {
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
	if err := r.Debug().First(&groupUserInfo).Error; err != nil {
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
	err = r.Debug().Table("group_members").Select([]string{"group_id"}).
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
	err = r.Debug().Table("groups").Where("group_id in ?", groupIds).Scan(&groupList).Error
	if err != nil {
		log.Infof("query group list failed, err:%+v", err)
		return err, nil
	}

	type Leader struct {
		GroupId int64 `gorm:"column:group_id"`
		Uid     int64 `gorm:"column:uid"`
	}
	leaders := make([]*Leader, 0)
	err = r.Debug().Table("group_members").Select([]string{"group_id", "uid"}).
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
			groupId, uid).Debug().Scan(&groupDetailInfo).Error
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
	err = r.Debug().Table(model.GroupMembers{}.TableName()).Select([]string{"group_id", "uid", "remark_name"}).
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
