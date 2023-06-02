package dao

import (
	"PInfo-server/log"
	"PInfo-server/utils"
	"context"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Dao is Data Access Object
type Dao struct {
	commDB *gorm.DB
	sf     *utils.Snowflake
}

// New creates Dao instance
// dsn eg: "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
func New(dsn string, dataCenterId, WorkerId int64) *Dao {
	d := &Dao{}

	cli, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("dao: New db gorm client error(%v)", err)
	}
	d.commDB = cli

	s, err := utils.NewSnowflake(dataCenterId, WorkerId)
	if err != nil {
		log.Fatalf("dao: NewSnowflake error(%v), dataCenterId:%d, WorkerId:%d", err, dataCenterId, WorkerId)
	}
	log.Infof("dao: NewSnowflake dataCenterId:%d, WorkerId:%d", dataCenterId, WorkerId)
	d.sf = s

	return d
}

func (d *Dao) db(ctx context.Context) *gorm.DB {
	return d.commDB.Debug()
}

// GenMsgID 生成消息ID。雪花算法，保证递增
func (d *Dao) GenMsgID() int64 {
	return d.sf.NextVal()
}
