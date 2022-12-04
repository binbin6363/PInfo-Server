package dao

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"PInfo-server/config"
	"PInfo-server/log"
	"PInfo-server/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
		log.InfoContextf(ctx, "user not exist")
		return err, false
	}

	if err != nil {
		log.InfoContextf(ctx, "query user name failed, err:%+v", err)
		return err, false
	}

	if res.UserName == "" {
		log.InfoContextf(ctx, "user %s not exist", username)
		return nil, false
	}

	log.InfoContextf(ctx, "user exist, name:%s, create time:%d", res.UserName, res.CreateTime)
	return nil, true
}

// AllocNewUserID 获取新用户ID
func (d *Dao) AllocNewUserID(ctx context.Context) (uid int64, err error) {
	r := d.db(ctx)
	type Result struct {
		Uid int64 `gorm:"column:uid"`
	}

	tx := r.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err = tx.Error; err != nil {
		return 0, err
	}

	res := &Result{}
	err = r.Model(&model.UserInfo{}).Order(clause.OrderByColumn{
		Column: clause.Column{
			Name: "uid",
		},
		Desc: true,
	}).Select([]string{"uid"}).First(&res).Error
	if err != nil {
		log.InfoContextf(ctx, "select max uid failed, err:%+v", err)
		tx.Rollback()
		return 0, err
	}

	uid = res.Uid + 1
	userInfo := &model.UserInfo{}
	userInfo.Uid = uid
	userInfo.UserName = fmt.Sprintf("%d", uid)
	userInfo.CreateTime = time.Now().Unix()
	err = r.Model(&model.UserInfo{}).Create(userInfo).Error
	if err != nil {
		log.InfoContextf(ctx, "alloc user id failed, err:%+v", err)
		tx.Rollback()
		return 0, err
	}

	err = tx.Commit().Error
	if err == nil {
		log.InfoContextf(ctx, "alloc user id ok, id:%d", uid)
	} else {
		uid = 0
		log.ErrorContextf(ctx, "alloc user id failed, err:%v", err)
	}

	return uid, err
}

// GetUserInfoByUserName 获取用户信息
func (d *Dao) GetUserInfoByUserName(ctx context.Context, username string) (error, *model.UserInfo) {
	r := d.db(ctx)
	if len(username) == 0 {
		log.Error("username is empty, invalid")
		return errors.New("username is invalid"), nil
	}

	userInfo := &model.UserInfo{}
	if err := r.Where("username=?", username).Limit(1).Find(&userInfo).Error; err != nil {
		log.InfoContextf(ctx, "GetUserInfoByUserName read db error(%v) username(%s)", err, username)
		return err, nil
	}
	userInfo.Avatar = d.makeFullAvatar(ctx, userInfo.Avatar)

	log.InfoContextf(ctx, "GetUserInfoByUserName read db ok username(%s), info:%+v", username, userInfo)
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
		log.InfoContextf(ctx, "GetUserInfoByUid read db error(%v) uid(%d)", err, uid)
		return err, nil
	}

	// 根据key生成URL
	userInfo.Avatar = d.makeFullAvatar(ctx, userInfo.Avatar)

	log.InfoContextf(ctx, "GetUserInfoByUid read db ok uid(%d)", uid)
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
		DoUpdates: clause.AssignmentColumns([]string{"passhash", "username", "nickname",
			"phone", "email", "avatar", "gender", "user_tag", "motto", "update_time"}),
	}).Create(userInfo)

	if err := r.Error; err != nil {
		log.InfoContextf(ctx, "SetUserInfo update db error(%v) user info:%+v", err, userInfo)
		return err
	}

	log.InfoContextf(ctx, "SetUserInfo update db ok user info:%+v", userInfo)
	return nil
}

// parseShortAvatar 从完整URL解析出 bucket + key 的路径
func (d *Dao) parseShortAvatar(ctx context.Context, avatarUrl string) string {
	shortAvatar := avatarUrl
	if len(avatarUrl) != 0 {
		if ava, err := d.ParseUrlKey(ctx, avatarUrl); err == nil {
			shortAvatar = ava
		} else {
			log.ErrorContextf(ctx, "parse url failed, use ori: %s", avatarUrl)
		}
	}
	return shortAvatar
}

// makeFullAvatar 构造全路径并获取前面
func (d *Dao) makeFullAvatar(ctx context.Context, avatar string) string {

	// 根据key生成URL
	if len(avatar) > 0 {
		if _, key, ok := strings.Cut(avatar, config.AppConfig().CosInfo.AvatarBucket); !ok {
			log.ErrorContextf(ctx, "parse key fail")
		} else {
			log.InfoContextf(ctx, "after cut, key: %s", key)
			str, e := d.GetPresignUrl(ctx, config.AppConfig().CosInfo.AvatarBucket,
				key, time.Duration(config.AppConfig().CosInfo.Expire))
			if e == nil {
				return str
			} else {
				log.ErrorContextf(ctx, "get avatar url err: %v", e)
			}
		}

	}

	log.ErrorContextf(ctx, "makeFullAvatar fail, use ori:%s", avatar)
	return avatar
}
