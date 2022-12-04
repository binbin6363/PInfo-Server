package main

import (
	"PInfo-server/routers"
	"PInfo-server/routers/auth"
	"PInfo-server/routers/chat"
	"PInfo-server/routers/contact"
	"PInfo-server/routers/emoticon"
	"PInfo-server/routers/group"
	"PInfo-server/routers/note"
	"PInfo-server/routers/sms"
	"PInfo-server/routers/users"
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
	routers.Register(auth.Routers, chat.Routers, users.Routers, group.Routers, note.Routers,
		contact.Routers, sms.Routers, emoticon.Routers)

	// 初始化路由
	r := routers.Init()

	service.Init()

	if err := r.Run(config.AppConfig().ServerInfo.Listen); err != nil {
		log.Fatalf("startup service failed, err:%v\n\n", err)
	}
}
