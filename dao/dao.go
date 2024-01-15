package dao

import (
	"context"

	"PInfo-server/config"
	"PInfo-server/log"
	"PInfo-server/utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Dao is Data Access Object
type Dao struct {
	commDB *gorm.DB
	sf     *utils.Snowflake
	sess   *session.Session
}

// New creates Dao instance
// dsn eg: "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
func New(dbInfo *config.DBInfo, svrInfo *config.ServerInfo, cosInfo *config.CosInfo) *Dao {
	d := &Dao{}

	cli, err := gorm.Open(mysql.Open(dbInfo.Dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("dao: New db gorm client error(%v)", err)
	}
	d.commDB = cli

	s, err := utils.NewSnowflake(svrInfo.DataCenterId, svrInfo.WorkerId)
	if err != nil {
		log.Fatalf("dao: NewSnowflake error(%v), dataCenterId:%d, WorkerId:%d",
			err, svrInfo.DataCenterId, svrInfo.WorkerId)
	}
	log.Infof("dao: NewSnowflake dataCenterId:%d, WorkerId:%d", svrInfo.DataCenterId, svrInfo.WorkerId)
	d.sf = s

	d.sess, _ = session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials(cosInfo.SecretID, cosInfo.SecretKey, ""),
		Endpoint:         aws.String(cosInfo.Domain),
		Region:           aws.String(cosInfo.Region),
		DisableSSL:       aws.Bool(cosInfo.DisableSSL),
		S3ForcePathStyle: aws.Bool(cosInfo.ForcePathStyle),
	})

	return d
}

func (d *Dao) db(ctx context.Context) *gorm.DB {
	return d.commDB.Debug()
}

// GenMsgID 生成消息ID。雪花算法，保证递增
func (d *Dao) GenMsgID() int64 {
	return d.sf.NextVal()
}
