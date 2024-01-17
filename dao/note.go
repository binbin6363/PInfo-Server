package dao

import (
	"PInfo-server/log"
	"PInfo-server/model"
	"context"
	"errors"

	"gorm.io/gorm/clause"
)

// EditArticle 新增/更新
func (d *Dao) EditArticle(ctx context.Context, articleInfo *model.Articles) error {
	r := d.db(ctx)
	if articleInfo.Uid == 0 {
		log.Error("uid invalid")
		return errors.New("uid invalid")
	}

	r = r.Clauses(clause.OnConflict{
		// key列
		Columns: []clause.Column{{Name: "id"}},
		// 需要更新的列。页面上仅支持这四列的手动修改。其他列的修改，都应该直接走server_list.csv更新（通用的）
		DoUpdates: clause.AssignmentColumns([]string{"content", "md_content", "update_time"}),
	}).Create(articleInfo)

	log.Infof("EditArticle update db ok id: %d", articleInfo.ID)
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
