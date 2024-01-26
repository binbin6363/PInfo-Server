package service

import (
	"PInfo-server/config"
	"PInfo-server/dao"
	"PInfo-server/log"
	"bytes"
	"context"
	"io"
	"net/http"
	"time"
)

var DefaultService *Service

// Service is service logic object
type Service struct {
	dao *dao.Dao
}

// New creates service instance
func New() *Service {
	srv := Service{
		dao: dao.New(config.AppConfig().DBInfo,
			config.AppConfig().ServerInfo,
			config.AppConfig().CosInfo),
	}
	return &srv
}

// Init .
func Init() {
	DefaultService = New()
}

// HttpPost timeout毫秒超时时间
func (s *Service) HttpPost(ctx context.Context, url string, data []byte, timeout int) ([]byte, error) {
	client := http.Client{
		Timeout: time.Millisecond * time.Duration(timeout), // 设置超时时间，单位毫秒
	}

	// 创建 POST 请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		log.ErrorContextf(ctx, "http.NewRequest err: %v", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	// 发送请求
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		log.ErrorContextf(ctx, "post err: %v", err)
		return nil, err
	} else {
		log.Infof("post ok, url: %s", url)
	}

	buffer := bytes.Buffer{}
	cnt, err := io.Copy(&buffer, resp.Body)
	if err != nil {
		log.ErrorContextf(ctx, "post copy err: %v", err)
	} else {
		log.Infof("recv post rsp data len: %d", cnt)
	}

	return buffer.Bytes(), err
}
