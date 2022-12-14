package service

import (
	"PInfo-server/config"
	"PInfo-server/dao"
)

var DefaultService *Service

// Service is service logic object
type Service struct {
	dao *dao.Dao
}

// New creates service instance
func New() *Service {
	srv := Service{
		dao: dao.New(config.AppConfig().DBInfo.Dsn,
			config.AppConfig().ServerInfo.DataCenterId,
			config.AppConfig().ServerInfo.WorkerId),
	}
	return &srv
}

// Init .
func Init() {
	DefaultService = New()
}
