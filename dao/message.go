package dao

import (
	"PInfo-server/log"
	"PInfo-server/model"
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// AddOneSingleMessage 添加单人消息
func (d *Dao) AddOneSingleMessage(ctx context.Context, msg *model.SingleMessages) error {
	r := d.db(ctx)
	if msg.Uid == 0 || msg.MsgID == 0 {
		log.Error("uid|msgid invalid")
		return errors.New("uid|msgid invalid")
	}

	if err := r.Create(msg).Error; err != nil {
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

	if err := r.Create(msg).Error; err != nil {
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
	}).Create(con)

	// 未读数加1
	r.UpdateColumn("unread", gorm.Expr("unread + ?", 1))

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
	}).Create(con)

	// 群成员未读数加1
	r.Table(model.Conversations{}.TableName()).Where("contact_id=?", con[0].ContactID).UpdateColumn("unread", gorm.Expr("unread + ?", 1))

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
		minMsgId = 0x7fffffffffffffff
	}
	if uid == peerUid {
		r = r.Where("uid=? and (senderid=? and receiverid=?) and msgid<?", uid, peerUid, peerUid, minMsgId)
	} else {
		r = r.Where("uid=? and (senderid=? or receiverid=?) and msgid<?", uid, peerUid, peerUid, minMsgId)
	}
	// 分页，取第index页的count条数据。倒序
	r = r.Order(clause.OrderByColumn{Column: clause.Column{Name: "msgid"}, Desc: true})
	r = r.Limit(limit)

	msgList := make([]*model.SingleMessages, 0)
	if err := r.Find(&msgList).Error; err != nil {
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
	if err := r.Find(&msgList).Error; err != nil {
		log.Infof("QueryGroupMessage read db error(%v)", err)
		return err, nil
	}

	log.Infof("QueryGroupMessage ok, msg size:%d", len(msgList))
	return nil, msgList
}
