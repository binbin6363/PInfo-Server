package dao

import (
	"PInfo-server/config"
	"PInfo-server/log"
	"PInfo-server/model"
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

func (d *Dao) CheckUserExist(ctx context.Context, username string) (err error, exist bool) {
	r := d.db(ctx)
	type Result struct {
		UserName   string `gorm:"column:username"`
		CreateTime int64  `gorm:"column:create_time"`
	}

	res := &Result{}
	err = r.Table("user_infos").Select([]string{"username", "create_time"}).
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
	err = r.Table("user_infos").Order(clause.OrderByColumn{
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

// GetUserInfoByUserName 获取用户信息
func (d *Dao) GetUserInfoByUserName(ctx context.Context, username string) (error, *model.UserInfo) {
	r := d.db(ctx)
	if username == "" {
		log.Error("username is empty, invalid")
		return errors.New("username is invalid"), nil
	}

	userInfo := &model.UserInfo{}
	if err := r.Where("username=?", username).Limit(1).Find(&userInfo).Error; err != nil {
		log.Infof("GetUserInfoByUserName read db error(%v) username(%s)", err, username)
		return err, nil
	}
	userInfo.Avatar = d.getFullAvatar(ctx, userInfo.Avatar)

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
	if err := r.Where("uid=?", uid).Limit(1).Find(&userInfo).Error; err != nil {
		log.Infof("GetUserInfoByUid read db error(%v) uid(%d)", err, uid)
		return err, nil
	}

	// 根据key生成URL
	userInfo.Avatar = d.getFullAvatar(ctx, userInfo.Avatar)

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

	// 拆解URL，只存储key
	userInfo.Avatar = d.parseShortAvatar(ctx, userInfo.Avatar)

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

func (d *Dao) parseShortAvatar(ctx context.Context, avatarUrl string) string {
	shortAvatar := avatarUrl
	if len(avatarUrl) != 0 {
		if ava, err := d.ParseUrlKey(ctx, avatarUrl); err == nil {
			shortAvatar = ava
		} else {
			log.Errorf("parse url failed, use ori: %s", avatarUrl)
		}
	}
	return shortAvatar
}

func (d *Dao) getFullAvatar(ctx context.Context, avatar string) string {
	fullAvatar := avatar
	// 根据key生成URL
	if len(avatar) > 0 {
		str, e := d.GetPresignUrl(ctx, config.AppConfig().CosInfo.AvatarBucket,
			avatar, time.Duration(config.AppConfig().CosInfo.Expire))
		if e == nil {
			fullAvatar = str
		} else {
			log.Errorf("get avatar url err: %v", e)
		}
	}
	return fullAvatar
}
