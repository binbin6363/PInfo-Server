package service

import (
	"PInfo-server/api"
	"PInfo-server/log"
	"PInfo-server/model"
	"PInfo-server/utils"
	"context"
	"time"
)

// EditArticle 新增/更新文章
func (s *Service) EditArticle(ctx context.Context, req *api.EditArticleReq) (*api.EditArticleRsp, error) {
	article := &model.Articles{
		ID:        req.ArticleId,
		Uid:       req.Uid,
		ClassId:   req.ClassId,
		Title:     req.Title,
		Content:   req.Content,
		MdContent: req.MdContent,
	}
	if err := s.dao.EditArticle(ctx, article); err != nil {
		log.Errorf("EditArticle err: %v, uid: %d", err, req.Uid)
		return nil, err
	}
	log.Infof("EditArticle ok")
	return nil, nil
}

func (s *Service) ClassList(ctx context.Context, req *api.ClassListReq) (*api.ClassListRsp, error) {
	rsp := &api.ClassListRsp{}

	rsp.ClassItems = append(rsp.ClassItems, api.ClassItem{
		Id:        1,
		ClassName: "默认分组",
		Count:     0,
		IsDefault: true,
		UpdatedAt: utils.FormatTimeStr(time.Now().Unix()),
	})
	log.Infof("done ClassList, rsp: %v", rsp)
	return rsp, nil
}

func (s *Service) ArticleList(ctx context.Context, req *api.ListArticleReq) (*api.ListArticleRsp, error) {
	rsp := &api.ListArticleRsp{}

	result, err := s.dao.ArticleList(ctx, req.Page, req.FindType, req.Uid, req.Cid, req.Keyword)
	if err != nil {
		log.Errorf("ArticleList err: %v", err)
		return nil, err
	}

	for idx := range result {
		item := result[idx]
		rsp.Items = append(rsp.Items, api.ArticleInfo{
			Id:        item.ID,
			Title:     item.Title,
			ClassId:   item.ClassId,
			Status:    1,
			Image:     "",
			Abstract:  "",
			UpdatedAt: utils.FormatTimeStr(item.UpdateTime),
		})
	}
	log.Infof("done ClassList, rsp: %v", rsp)
	return rsp, nil
}
