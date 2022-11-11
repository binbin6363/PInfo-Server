package main

import (
	"PInfo-server/routers"
	"PInfo-server/routers/login"
	"PInfo-server/service"
	"flag"
	"log"

	"PInfo-server/config"

	"github.com/gin-gonic/gin"
)

func main() {
	confFile := flag.String("f", "../etc/conf.yaml", "配置文件路径")

	flag.Parse()
	config.Init(*confFile)
	gin.SetMode(gin.DebugMode)
	// 加载多个APP的路由配置
	routers.Register(login.Routers)
	// 初始化路由
	r := routers.Init()

	service.Init()

	if err := r.Run(config.AppConfig().ServerInfo.Listen); err != nil {
		log.Fatalf("startup service failed, err:%v\n\n", err)
	}
}
