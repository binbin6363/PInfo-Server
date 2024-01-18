package dao

import (
	"PInfo-server/log"
	"PInfo-server/model"
	"context"
	"errors"
	"time"

	"gorm.io/gorm/clause"
)

// EditArticle 新增/更新
func (d *Dao) EditArticle(ctx context.Context, articleInfo *model.Articles) error {
	r := d.db(ctx)
	if articleInfo.Uid == 0 {
		log.Error("uid invalid")
		return errors.New("uid invalid")
	}

	articleInfo.UpdateTime = time.Now().Unix()
	if articleInfo.CreateTime == 0 {
		articleInfo.CreateTime = articleInfo.UpdateTime
	}
	var err error
	if articleInfo.ID == 0 {
		log.Infof("create article, uid: %d, title: %s", articleInfo.Uid, articleInfo.Title)
		err = r.Create(articleInfo).Error
	} else {
		log.Infof("update article, uid: %d, id: %d", articleInfo.Uid, articleInfo.ID)
		err = r.Updates(articleInfo).Error
	}
	if err != nil {
		log.Error("EditArticle err: %v, uid: %d", err, articleInfo.Uid)
	} else {
		log.Infof("EditArticle ok, uid: %d, id: %d", articleInfo.Uid, articleInfo.ID)
	}

	return nil
}

// ArticleList 拉取文章列表
// todo：先把page当成ID查
func (d *Dao) ArticleList(ctx context.Context, page, findType int, uid, cid int64, kw string) ([]*model.Articles, error) {
	r := d.db(ctx)
	if page == 0 {
		log.Error("page invalid")
		return nil, errors.New("page invalid")
	}

	r = r.Where("uid=? and id>?", uid, page)
	if len(kw) > 0 {
		r = r.Where("title like ?", kw)
	}
	// 分页，取第index页的count条数据。倒序
	limit := 10
	r = r.Order(clause.OrderByColumn{Column: clause.Column{Name: "id"}, Desc: true})
	r = r.Limit(limit)

	articleList := make([]*model.Articles, 0)
	if err := r.Select([]string{"id", "class_id", "title", "update_time"}).Find(&articleList).Error; err != nil {
		log.Infof("ArticleList read db error(%v)", err)
		return nil, err
	}

	log.Infof("ArticleList ok, size:%d", len(articleList))
	return articleList, nil
}

func (d *Dao) ArticleDetail(ctx context.Context, uid, articleId int64) (*model.Articles, error) {
	r := d.db(ctx)
	if uid == 0 || articleId == 0 {
		log.Error("uid|articleId invalid")
		return nil, errors.New("uid|articleId invalid")
	}

	info := &model.Articles{}
	err := r.Where("id=? and uid=?", articleId, uid).Find(info).Error
	if err != nil {
		log.Error("ArticleDetail err: %v", err)
		return nil, err
	}
	log.Infof("ArticleDetail ok id: %d", articleId)
	return info, nil
}
