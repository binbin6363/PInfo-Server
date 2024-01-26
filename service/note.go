package service

import (
	"PInfo-server/api"
	"PInfo-server/log"
	"PInfo-server/model"
	"PInfo-server/utils"
	"context"
	"time"
)

// ArticleEdit 新增/更新文章
func (s *Service) ArticleEdit(ctx context.Context, req *api.ArticleEditReq) (*api.ArticleEditRsp, error) {
	article := &model.Articles{
		ID:        req.ArticleId,
		Uid:       req.Uid,
		ClassId:   req.ClassId,
		Title:     req.Title,
		MdContent: req.MdContent,
	}
	if err := s.dao.EditArticle(ctx, article); err != nil {
		log.ErrorContextf(ctx, "ArticleEdit err: %v, uid: %d", err, req.Uid)
		return nil, err
	}
	log.InfoContextf(ctx, "ArticleEdit ok")
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
	log.InfoContextf(ctx, "done ClassList, rsp: %v", rsp)
	return rsp, nil
}

func (s *Service) ArticleList(ctx context.Context, req *api.ArticleListReq) (*api.ArticleListRsp, error) {
	rsp := &api.ArticleListRsp{}

	result, err := s.dao.ArticleList(ctx, req.Page, req.FindType, req.Uid, req.Cid, req.Keyword)
	if err != nil {
		log.ErrorContextf(ctx, "ArticleList err: %v", err)
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
	log.InfoContextf(ctx, "done ClassList, rsp: %v", rsp)
	return rsp, nil
}

func (s *Service) ArticleDetail(ctx context.Context, req *api.ArticleDetailReq) (*api.ArticleDetailRsp, error) {

	result, err := s.dao.ArticleDetail(ctx, req.Uid, req.ArticleId)
	if err != nil {
		log.ErrorContextf(ctx, "ArticleDetail err: %v", err)
		return nil, err
	}

	rsp := &api.ArticleDetailRsp{
		Id:         result.ID,
		Title:      result.Title,
		Classify:   "",
		Abstract:   "",
		Image:      "",
		ClassId:    result.ClassId,
		IsAsterisk: 1,
		MdContent:  &result.MdContent,
		UpdatedAt:  utils.FormatTimeStr(result.UpdateTime),
	}

	log.InfoContextf(ctx, "done ClassList, rsp: %v", rsp)
	return rsp, nil
}
