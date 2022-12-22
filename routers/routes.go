package routers

import (
	"PInfo-server/config"
	"net/http"
	"strings"

	"PInfo-server/log"
	"PInfo-server/service"
	"github.com/gin-gonic/gin"
	gindump "github.com/tpkeeper/gin-dump"
)

type Option func(*gin.Engine)

var options []Option

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.GetHeader("Origin") //请求头部
		if origin != "" {
			// 可将将* 替换为指定的域名
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// JWTAuthMiddleware 基于JWT的认证中间件
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 首次的登录不校验token
		if strings.Contains(c.Request.URL.Path, "/auth/login") ||
			strings.Contains(c.Request.URL.Path, "/auth/register") {
			return
		}
		// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
		// 这里假设Token放在Header的Authorization中，并使用Bearer开头
		// 这里的具体实现方式要依据你的实际业务情况决定
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusOK, gin.H{
				"code": 4000,
				"msg":  "没有认证信息",
			})
			c.Abort()
			return
		}
		// 按空格分割
		parts := strings.Fields(authHeader)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusOK, gin.H{
				"code": 4001,
				"msg":  "认证信息鉴权失败",
			})
			c.Abort()
			return
		}
		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		mc, err := service.DefaultService.ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 4002,
				"msg":  "无效的认证信息",
			})
			c.Abort()
			return
		}
		// 将当前请求的username信息保存到请求的上下文c上
		c.Set("uid", mc.Id)
		c.Set("username", mc.Audience)
		c.Next() // 后续的处理函数可以用过c.Get("username")来获取当前请求的用户信息
	}
}

// Register 注册路由配置
func Register(opts ...Option) {
	options = append(options, opts...)
}

// Init 初始化
func Init() *gin.Engine {
	r := gin.Default()
	r.Use(Cors())
	r.Use(JWTAuthMiddleware())
	if config.AppConfig().ServerInfo.DebugReqRsp {
		r.Use(gindump.DumpWithOptions(true, true, true, false, false, func(dumpStr string) {
			log.Info(dumpStr)
		}))
	}
	for _, opt := range options {
		opt(r)
	}
	return r
}
