package service

import (
	"PInfo-server/api"
	"PInfo-server/log"
	"PInfo-server/model"
	"PInfo-server/utils"
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	DefaultClassName = "默认分类"
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
	if err := s.dao.ArticleEdit(ctx, article); err != nil {
		log.ErrorContextf(ctx, "ArticleEdit err: %v, uid: %d", err, req.Uid)
		return nil, err
	}
	log.InfoContextf(ctx, "ArticleEdit ok")
	return nil, nil
}

// ArticleList 拉取文章列表
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

// ArticleDetail 拉取文章详情
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

// ClassList 拉取分类列表
func (s *Service) ClassList(ctx context.Context, req *api.ClassListReq) (*api.ClassListRsp, error) {
	rsp := &api.ClassListRsp{}
	data, err := s.dao.ClassList(ctx, utils.GetUid(ctx))
	if err != nil {
		log.ErrorContextf(ctx, "ClassList err: %v, uid: %d", err, utils.GetUid(ctx))
		return nil, err
	}
	log.InfoContextf(ctx, "ClassList ok")
	for _, d := range data {
		rsp.ClassItems = append(rsp.ClassItems, api.ClassItem{
			Id:        d.ID,
			ClassName: d.Name,
			Count:     0,
			IsDefault: d.Flag == 0,
			UpdatedAt: utils.FormatTimeStr(d.UpdateTime),
		})
	}
	// 添加默认分类
	rsp.ClassItems = append([]api.ClassItem{{
		Id:        0,
		ClassName: DefaultClassName,
		Count:     0,
		IsDefault: true,
		UpdatedAt: utils.FormatTimeStr(time.Now().Unix()),
	}}, rsp.ClassItems...)
	return rsp, nil
}

// ClassEdit 新增/更新分类
func (s *Service) ClassEdit(ctx context.Context, req *api.ClassEditReq) (*api.ClassEditRsp, error) {
	uid := utils.GetUid(ctx)
	cla := &model.Classes{
		ID:   req.ClassId,
		Uid:  uid,
		Flag: 1,
		Name: req.ClassName,
	}

	if err := s.dao.ClassEdit(ctx, cla); err != nil {
		log.ErrorContextf(ctx, "ClassEdit err: %v, uid: %d", err, uid)
		return nil, err
	}
	log.InfoContextf(ctx, "ClassEdit ok")
	return &api.ClassEditRsp{ClassId: req.ClassId}, nil
}

// ClassDelete 删除分类
func (s *Service) ClassDelete(ctx context.Context, req *api.ClassDeleteReq) (*api.ClassDeleteRsp, error) {
	uid := utils.GetUid(ctx)
	cla := &model.Classes{
		ID:  req.ClassId,
		Uid: uid,
	}

	if err := s.dao.ClassDelete(ctx, cla); err != nil {
		log.ErrorContextf(ctx, "ClassDelete err: %v, uid: %d", err, uid)
		return nil, err
	}
	log.InfoContextf(ctx, "ClassDelete ok")
	return &api.ClassDeleteRsp{ClassId: req.ClassId}, nil
}

func (s *Service) ClassSort(c *gin.Context, req *api.ClassSortReq) (*api.ClassSortRsp, error) {
	return nil, errors.New("no implemented")
}

func (s *Service) ClassMerge(c *gin.Context, req *api.ClassMergeReq) (*api.ClassMergeRsp, error) {
	return nil, errors.New("no implemented")
}

// ========== tag ==========

// TagList 拉取tag列表
func (s *Service) TagList(ctx context.Context, req *api.TagListReq) (*api.TagListRsp, error) {
	rsp := &api.TagListRsp{}
	data, err := s.dao.TagList(ctx, utils.GetUid(ctx))
	if err != nil {
		log.ErrorContextf(ctx, "TagList err: %v, uid: %d", err, utils.GetUid(ctx))
		return nil, err
	}
	log.InfoContextf(ctx, "TagList ok")
	for _, d := range data {
		rsp.TagItems = append(rsp.TagItems, api.TagItem{
			Id:        d.ID,
			TagName:   d.Name,
			Count:     0,
			IsDefault: d.Flag == 0,
			UpdatedAt: utils.FormatTimeStr(d.UpdateTime),
		})
	}
	return rsp, nil
}

// TagEdit 新增/更新分类
func (s *Service) TagEdit(ctx context.Context, req *api.TagEditReq) (*api.TagEditRsp, error) {
	uid := utils.GetUid(ctx)
	cla := &model.Tags{
		ID:   req.TagId,
		Uid:  uid,
		Flag: 1,
		Name: req.TagName,
	}

	if err := s.dao.TagEdit(ctx, cla); err != nil {
		log.ErrorContextf(ctx, "TagEdit err: %v, uid: %d", err, uid)
		return nil, err
	}
	log.InfoContextf(ctx, "TagEdit ok")
	return &api.TagEditRsp{TagId: req.TagId}, nil
}

// TagDelete 删除分类
func (s *Service) TagDelete(ctx context.Context, req *api.TagDeleteReq) (*api.TagDeleteRsp, error) {
	uid := utils.GetUid(ctx)
	cla := &model.Tags{
		ID:  req.TagId,
		Uid: uid,
	}

	if err := s.dao.TagDelete(ctx, cla); err != nil {
		log.ErrorContextf(ctx, "TagDelete err: %v, uid: %d", err, uid)
		return nil, err
	}
	log.InfoContextf(ctx, "TagDelete ok")
	return &api.TagDeleteRsp{TagId: req.TagId}, nil
}

func (s *Service) ArticleMoveClass(ctx context.Context, req *api.ArticleMoveClassReq) (*api.ArticleMoveClassRsp, error) {
	uid := utils.GetUid(ctx)
	cla := &model.Articles{
		ID:      req.ArticleId,
		Uid:     uid,
		ClassId: req.ClassId,
	}

	if err := s.dao.ArticleMove(ctx, cla); err != nil {
		log.ErrorContextf(ctx, "ArticleMoveClass err: %v, uid: %d", err, uid)
		return nil, err
	}
	log.InfoContextf(ctx, "ArticleMoveClass ok")
	return &api.ArticleMoveClassRsp{ArticleId: req.ArticleId}, nil
}
