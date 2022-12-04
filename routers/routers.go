package routers

import (
	"PInfo-server/routers/auth"
	"PInfo-server/routers/chat"
	"PInfo-server/routers/contact"
	"PInfo-server/routers/emoticon"
	"PInfo-server/routers/group"
	"PInfo-server/routers/im"
	"PInfo-server/routers/note"
	"PInfo-server/routers/sms"
	"PInfo-server/routers/upload"
	"PInfo-server/routers/users"

	"github.com/gin-gonic/gin"
)

type Option func(*gin.Engine)

var options []Option

func init() {
	Register(auth.Routers)
	Register(chat.Routers)
	Register(contact.Routers)
	Register(emoticon.Routers)
	Register(group.Routers)
	Register(im.Routers)
	Register(note.Routers)
	Register(sms.Routers)
	Register(upload.Routers)
	Register(users.Routers)
}

// Register 注册路由配置
func Register(opts ...Option) {
	options = append(options, opts...)
}

func Routes() []Option {
	return options
}
