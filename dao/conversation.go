package dao

import (
	"PInfo-server/log"
	"PInfo-server/model"
	"context"
	"errors"
	"gorm.io/gorm/clause"
)

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
	if err := r.Scan(&conDetails).Error; err != nil {
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
	if err := r.Where("uid=? and contact_id=?", uid, contactId).First(&conversationInfo).Error; err != nil {
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

	r = r.Model(&model.Conversations{}).Where("contact_id = ?", conversationInfo.ContactID).
		UpdateColumn("conversation_name", conversationInfo.ConversationName)
	if err := r.Error; err != nil {
		log.Infof("BatchSetGroupConversationName update db error(%v) user info:%+v", err, conversationInfo)
		return err
	}

	log.Infof("BatchSetGroupConversationName update db ok user info:%+v", conversationInfo)
	return nil
}
