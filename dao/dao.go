package dao

import (
	"PInfo-server/model"
	"context"
	"errors"
	"github.com/GUAIK-ORG/go-snowflake/snowflake"
	"gorm.io/gorm/clause"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Dao is Data Access Object
type Dao struct {
	commDB *gorm.DB
	sf     *snowflake.Snowflake
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

	s, err := snowflake.NewSnowflake(dataCenterId, WorkerId)
	if err != nil {
		log.Fatalf("dao: NewSnowflake error(%v), dataCenterId:%d, WorkerId:%d", err, dataCenterId, WorkerId)
	}
	log.Printf("dao: NewSnowflake dataCenterId:%d, WorkerId:%d\n", dataCenterId, WorkerId)
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

// GetUserInfoByUserName 获取用户信息
func (d *Dao) GetUserInfoByUserName(ctx context.Context, username string) (error, *model.UserInfo) {
	r := d.db(ctx)
	if username == "" {
		log.Println("username is empty, invalid")
		return errors.New("username is invalid"), nil
	}

	userInfo := &model.UserInfo{}
	if err := r.Debug().Where("username=?", username).Limit(1).Find(&userInfo).Error; err != nil {
		log.Printf("GetUserInfoByUserName read db error(%v) username(%s)\n", err, username)
		return err, nil
	}

	log.Printf("GetUserInfoByUserName read db ok username(%s), info:%+v\n", username, userInfo)
	return nil, userInfo
}

// GetUserInfoByUid 获取用户信息
func (d *Dao) GetUserInfoByUid(ctx context.Context, uid int64) (error, *model.UserInfo) {
	r := d.db(ctx)
	if uid == 0 {
		log.Println("uid is invalid")
		return errors.New("uid is invalid"), nil
	}

	userInfo := &model.UserInfo{}
	if err := r.Debug().Where("uid=?", uid).Limit(1).Find(&userInfo).Error; err != nil {
		log.Printf("GetUserInfoByUid read db error(%v) uid(%d)\n", err, uid)
		return err, nil
	}

	log.Printf("GetUserInfoByUid read db ok uid(%d)\n", uid)
	return nil, userInfo
}

// SetUserInfo 设置用户信息
func (d *Dao) SetUserInfo(ctx context.Context, userInfo *model.UserInfo) error {
	r := d.db(ctx)
	if userInfo.Uid == 0 {
		log.Println("uid invalid")
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
		log.Printf("SetUserInfo update db error(%v) user info:%+v\n", err, userInfo)
		return err
	}

	log.Printf("SetUserInfo update db ok user info:%+v\n", userInfo)
	return nil
}

// GetContactList 获取好友列表信息
// select user_infos.uid as uid, user_infos.nickname as nickname, user_infos.phone as phone,user_infos.email
// as email,user_infos.avatar as avatar,user_infos.gender as gender, contacts.contact_id as contact_id,
// contacts.remark_name as remark_name from user_infos left join contacts on contacts.uid=user_infos.uid
// where user_infos.uid=10000;
func (d *Dao) GetContactList(ctx context.Context, uid int64) (error, []*model.UserContact) {
	r := d.db(ctx)
	if uid == 0 {
		log.Println("uid is invalid")
		return errors.New("uid is invalid"), nil
	}

	r.Table(model.UserInfo{}.TableName()).Select("user_infos.uid as uid, user_infos.nickname as nickname, "+
		"user_infos.gender as gender, user_infos.motto as motto, user_infos.avatar as avatar, "+
		"contacts.remark_name as remark_name").
		Joins("left join contacts on contacts.uid=user_infos.uid").Where("user_infos.uid=?", uid)

	userContacts := make([]*model.UserContact, 0)
	if err := r.Debug().Find(&userContacts).Error; err != nil {
		log.Printf("GetContactList read db error(%v) uid(%d)\n", err, uid)
		return err, nil
	}

	log.Printf("GetContactList read db ok uid(%d)", uid)
	return nil, userContacts
}

func (d *Dao) AddOneMessage(ctx context.Context, msg *model.SingleMessages) error {
	r := d.db(ctx)
	if msg.Uid == 0 || msg.MsgID == 0 {
		log.Println("uid|msgid invalid")
		return errors.New("uid|msgid invalid")
	}

	if err := r.Debug().Create(msg).Error; err != nil {
		log.Printf("AddOneMessage insert db error(%v) msg:%+v\n", err, msg)
		return err
	}

	log.Printf("AddOneMessage insert db ok msg:%+v\n", msg)
	return nil
}

// UpdateConversationMsg 更新会话列表及未读数
func (d *Dao) UpdateConversationMsg(ctx context.Context, con *model.Conversations) error {
	r := d.db(ctx)
	if con.Uid == 0 {
		log.Println("uid invalid")
		return errors.New("uid invalid")
	}

	r = r.Clauses(clause.OnConflict{
		// key列
		Columns: []clause.Column{{Name: "uid"}, {Name: "contact_id"}},
		// 需要更新的列。页面上仅支持这四列的手动修改。其他列的修改，都应该直接走server_list.csv更新（通用的）
		DoUpdates: clause.AssignmentColumns([]string{"msgid", "msg_digest", "update_time"}),
	}).Create(con)

	// 未读数加1
	r.UpdateColumn("unread", gorm.Expr("unread + ?", 1))

	log.Printf("UpdateConversationMsg update db ok conversations info:%+v\n", con)
	return nil
}

// QuerySingleMessage 拉取单人历史消息
func (d *Dao) QuerySingleMessage(ctx context.Context, uid, peerUid, minMsgId int64, limit int) (error, []*model.SingleMessages) {
	r := d.db(ctx)
	if uid == 0 || peerUid == 0 {
		log.Println("uid|msgid invalid")
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
		log.Printf("QuerySingleMessage read db error(%v)\n", err)
		return err, nil
	}

	log.Printf("QuerySingleMessage ok, msg size:%d", len(msgList))
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
		log.Printf("user not exist\n")
		return err, false
	}

	if err != nil {
		log.Printf("query user name failed, err:%+v\n", err)
		return err, false
	}

	if res.UserName == "" {
		log.Printf("user %s not exist\n", username)
		return nil, false
	}

	log.Printf("user exist, name:%s, create time:%d", res.UserName, res.CreateTime)
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
		log.Printf("alloc user id failed, err:%+v", err)
		return err, 0
	}

	uid = res.Uid + 1
	log.Printf("alloc user id ok, id:%d", uid)
	return nil, uid
}
