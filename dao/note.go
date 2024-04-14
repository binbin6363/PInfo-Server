package dao

import (
	"PInfo-server/log"
	"PInfo-server/model"
	"context"
	"errors"
	"time"

	"gorm.io/gorm/clause"
)

const (
	MaxClassNum = 100
)

type FindType int

const (
	Recently FindType = 1
	Marked   FindType = 2
	Classfy  FindType = 3
	Tag      FindType = 4
	Recyle   FindType = 5
	KeyWord  FindType = 6
)

// ArticleEdit 新增/更新
func (d *Dao) ArticleEdit(ctx context.Context, articleInfo *model.Articles) error {
	r := d.db(ctx)
	if articleInfo.Uid == 0 {
		log.ErrorContextf(ctx, "uid invalid")
		return errors.New("uid invalid")
	}

	articleInfo.UpdateTime = time.Now().Unix()
	if articleInfo.ID == 0 {
		articleInfo.CreateTime = articleInfo.UpdateTime
	}
	var err error
	if articleInfo.ID == 0 {
		log.InfoContextf(ctx, "create article, uid: %d, title: %s", articleInfo.Uid, articleInfo.Title)
		err = r.Create(articleInfo).Error
	} else {
		log.InfoContextf(ctx, "update article, uid: %d, id: %d", articleInfo.Uid, articleInfo.ID)
		err = r.Updates(articleInfo).Error
	}
	if err != nil {
		log.ErrorContextf(ctx, "ArticleEdit err: %v, uid: %d", err, articleInfo.Uid)
	} else {
		log.InfoContextf(ctx, "ArticleEdit ok, uid: %d, id: %d", articleInfo.Uid, articleInfo.ID)
	}

	return nil
}

// ArticleList 拉取文章列表
// todo：先把page当成ID查
func (d *Dao) ArticleList(ctx context.Context, page, findType int, uid, cid int64, kw string) ([]*model.Articles, error) {
	r := d.db(ctx)
	if page == 0 {
		log.ErrorContextf(ctx, "page invalid")
		return nil, errors.New("page invalid")
	}

	r = r.Where("uid=? and id>?", uid, page)

	switch FindType(findType) {
	case Recently:
		r = r.Where("update_time>?", time.Now().Add(-24*time.Hour).Unix()) // 24小时内编辑的当作最近编辑
	case Marked:
	// r = r.Where("")
	// do nothing
	case Classfy:
		r = r.Where("class_id=?", cid)
	case Tag:
	// r = r.Where("class_id=?", cid)
	// todo: 待实现
	case Recyle:
	// todo: 待实现
	case KeyWord:
		if len(kw) > 0 {
			r = r.Where("title like ?", "%"+kw+"%")
		}
	default:
		log.ErrorContextf(ctx, "unknown search type: %d", findType)
		return nil, errors.New("unknown search type")
	}

	// 分页，取第index页的count条数据。倒序
	limit := 10
	r = r.Order(clause.OrderByColumn{Column: clause.Column{Name: "id"}, Desc: true})
	r = r.Limit(limit)

	articleList := make([]*model.Articles, 0)
	if err := r.Select([]string{"id", "class_id", "title", "update_time"}).Find(&articleList).Error; err != nil {
		log.InfoContextf(ctx, "ArticleList read db error(%v)", err)
		return nil, err
	}

	log.InfoContextf(ctx, "ArticleList ok, size:%d", len(articleList))
	return articleList, nil
}

// ArticleDetail 查询文章详情
func (d *Dao) ArticleDetail(ctx context.Context, uid, articleId int64) (*model.Articles, error) {
	r := d.db(ctx)
	if uid == 0 || articleId == 0 {
		log.ErrorContextf(ctx, "uid|articleId invalid")
		return nil, errors.New("uid|articleId invalid")
	}

	info := &model.Articles{}
	err := r.Where("id=? and uid=?", articleId, uid).Find(info).Error
	if err != nil {
		log.ErrorContextf(ctx, "ArticleDetail err: %v", err)
		return nil, err
	}
	log.InfoContextf(ctx, "ArticleDetail ok id: %d", articleId)
	return info, nil
}

// ClassList 查询分类列表
func (d *Dao) ClassList(ctx context.Context, uid int64) ([]*model.Classes, error) {
	r := d.db(ctx)
	r = r.Where("uid=?", uid)

	// 倒序
	limit := MaxClassNum
	r = r.Order(clause.OrderByColumn{Column: clause.Column{Name: "id"}, Desc: true})
	r = r.Limit(limit)

	cList := make([]*model.Classes, 0)
	if err := r.Select([]string{"id", "uid", "flag", "name", "update_time"}).Find(&cList).Error; err != nil {
		log.InfoContextf(ctx, "ClassList read db error(%v)", err)
		return nil, err
	}

	log.InfoContextf(ctx, "ClassList ok, size:%d", len(cList))
	return cList, nil
}

// ClassEdit 编辑或新增分类
func (d *Dao) ClassEdit(ctx context.Context, cla *model.Classes) error {
	r := d.db(ctx)
	if cla.Uid == 0 {
		log.ErrorContextf(ctx, "uid invalid")
		return errors.New("uid invalid")
	}
	r = r.Where("uid=?", cla.Uid)

	cla.UpdateTime = time.Now().Unix()
	if cla.ID == 0 {
		cla.CreateTime = cla.UpdateTime // 首次则添加创建时间
	}

	var err error
	if cla.ID == 0 {
		log.InfoContextf(ctx, "create class, uid: %d, title: %s", cla.Uid, cla.Name)
		err = r.Create(cla).Error
	} else {
		log.InfoContextf(ctx, "update class, uid: %d, id: %d", cla.Uid, cla.ID)
		err = r.Updates(cla).Error
	}
	if err != nil {
		log.ErrorContextf(ctx, "ClassEdit err: %v, uid: %d", err, cla.Uid)
	} else {
		log.InfoContextf(ctx, "ClassEdit ok, uid: %d, id: %d", cla.Uid, cla.ID)
	}

	return err
}

// ClassDelete 删除分类
func (d *Dao) ClassDelete(ctx context.Context, cla *model.Classes) error {
	r := d.db(ctx)
	if cla.Uid == 0 {
		log.ErrorContextf(ctx, "uid invalid")
		return errors.New("uid invalid")
	}
	r = r.Where("uid=?", cla.Uid)

	if cla.ID == 0 {
		return errors.New("delete class but id not set")
	}

	err := r.Delete(cla).Limit(1).Error
	if err != nil {
		log.ErrorContextf(ctx, "ClassDelete err: %v, uid: %d", err, cla.Uid)
	} else {
		log.InfoContextf(ctx, "ClassDelete ok, uid: %d, id: %d", cla.Uid, cla.ID)
	}

	return err
}

// ==== tag

// TagList 查询分类列表
func (d *Dao) TagList(ctx context.Context, uid int64) ([]*model.Tags, error) {
	r := d.db(ctx)
	r = r.Where("uid=?", uid)

	// 倒序
	limit := MaxClassNum
	r = r.Order(clause.OrderByColumn{Column: clause.Column{Name: "id"}, Desc: true})
	r = r.Limit(limit)

	tList := make([]*model.Tags, 0)
	if err := r.Select([]string{"id", "uid", "flag", "name", "update_time"}).Find(&tList).Error; err != nil {
		log.InfoContextf(ctx, "TagList read db error(%v)", err)
		return nil, err
	}

	log.InfoContextf(ctx, "TagList ok, size:%d", len(tList))
	return tList, nil
}

// TagEdit 编辑或新增分类
func (d *Dao) TagEdit(ctx context.Context, tag *model.Tags) error {
	r := d.db(ctx)
	if tag.Uid == 0 {
		log.ErrorContextf(ctx, "uid invalid")
		return errors.New("uid invalid")
	}
	r = r.Where("uid=?", tag.Uid)

	tag.UpdateTime = time.Now().Unix()
	if tag.ID == 0 {
		tag.CreateTime = tag.UpdateTime // 首次则添加创建时间
	}

	var err error
	if tag.ID == 0 {
		log.InfoContextf(ctx, "create tag, uid: %d, title: %s", tag.Uid, tag.Name)
		err = r.Create(tag).Error
	} else {
		log.InfoContextf(ctx, "update tag, uid: %d, id: %d", tag.Uid, tag.ID)
		err = r.Updates(tag).Error
	}
	if err != nil {
		log.ErrorContextf(ctx, "TagEdit err: %v, uid: %d", err, tag.Uid)
	} else {
		log.InfoContextf(ctx, "TagEdit ok, uid: %d, id: %d", tag.Uid, tag.ID)
	}

	return err
}

// TagDelete 删除分类
func (d *Dao) TagDelete(ctx context.Context, tag *model.Tags) error {
	r := d.db(ctx)
	if tag.Uid == 0 {
		log.ErrorContextf(ctx, "uid invalid")
		return errors.New("uid invalid")
	}
	r = r.Where("uid=?", tag.Uid)

	if tag.ID == 0 {
		return errors.New("delete tag but id not set")
	}

	err := r.Delete(tag).Limit(1).Error
	if err != nil {
		log.ErrorContextf(ctx, "TagDelete err: %v, uid: %d", err, tag.Uid)
	} else {
		log.InfoContextf(ctx, "TagDelete ok, uid: %d, id: %d", tag.Uid, tag.ID)
	}

	return err
}

// ArticleMove 修改文章分类
func (d *Dao) ArticleMove(ctx context.Context, article *model.Articles) error {
	r := d.db(ctx)
	if article.Uid == 0 {
		log.ErrorContextf(ctx, "uid invalid")
		return errors.New("uid invalid")
	}
	if article.ID == 0 {
		log.ErrorContextf(ctx, "article id invalid")
		return errors.New("article id invalid")
	}
	r = r.Where("id=? AND uid=?", article.ID, article.Uid)
	article.UpdateTime = time.Now().Unix()

	var err error
	log.InfoContextf(ctx, "ArticleMove, uid: %d, article id: %d, to class id: %d",
		article.Uid, article.ID)
	err = r.Select("class_id", "update_time").Updates(article).Error

	if err != nil {
		log.ErrorContextf(ctx, "ArticleMove err: %v, uid: %d", err, article.Uid)
	} else {
		log.InfoContextf(ctx, "ArticleMove ok, uid: %d, id: %d", article.Uid, article.ID)
	}

	return err
}
