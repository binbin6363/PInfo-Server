package routers

import "github.com/gin-gonic/gin"

type Option func(*gin.Engine)

var options []Option

// Register 注册路由配置
func Register(opts ...Option) {
	options = append(options, opts...)
}

func Routes() []Option {
	return options
}
