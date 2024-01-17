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
