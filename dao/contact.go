package dao

import (
	"context"
	"errors"

	"PInfo-server/log"
	"PInfo-server/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

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
	if err := r.Scan(&userContacts).Error; err != nil {
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
	if err := r.First(&userContacts).Error; err != nil {
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
	if err := r.Where("uid=? and contact_id=?", uid, contactId).First(&contactInfo).Error; err != nil {
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
