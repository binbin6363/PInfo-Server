package utils

import (
	"PInfo-server/api"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

// EncryptPassword 对密码加密
func EncryptPassword(password string) (error, string) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return err, string(bytes)
}

// CheckPasswordHash 密码校验
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GetUid 从上下文获取uid
func GetUid(c *gin.Context) (uid int64) {
	if aUid, ok := c.Get("uid"); ok {
		uid = cast.ToInt64(aUid)
		return uid
	}
	return 0
}

// GetUserName 从上下文获取username
func GetUserName(c *gin.Context) (err error, username string) {
	if aUserName, ok := c.Get("username"); ok {
		username = cast.ToString(aUserName)
		return nil, username
	}
	return errors.New("gin get username failed"), ""
}

// SendJsonRsp 回复json消息
func SendJsonRsp(c *gin.Context, rsp *api.CommRsp) {
	if rsp == nil {
		log.Printf("[INFO] handle ok, send empty rsp\n")
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "ok",
			"data":    nil,
		})
	} else if rsp.Code == 0 || rsp.Code == 200 {
		rsp.Code = 200
		log.Printf("[INFO] handle ok, send rsp:%+v\n", rsp)
		c.JSON(http.StatusOK, rsp)
	} else {
		log.Printf("[INFO] handle error, code:%d, message:%s\n", rsp.Code, rsp.Message)
		c.JSON(http.StatusOK, rsp)
	}
}
