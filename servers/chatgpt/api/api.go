package api

import (
	"github.com/LSDXXX/libs/api"
	"github.com/LSDXXX/libs/pkg/container"
	"github.com/LSDXXX/libs/pkg/util"
	"github.com/LSDXXX/servers/chatgpt/api/handlers/auth"
	"github.com/LSDXXX/servers/chatgpt/api/handlers/chat"
	"github.com/LSDXXX/servers/chatgpt/api/handlers/conversation"
	"github.com/LSDXXX/servers/chatgpt/bot"
	"github.com/LSDXXX/servers/chatgpt/config"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

type staticFileHandler struct {
	path string
}

func (s *staticFileHandler) Use(e *gin.Engine) {
	e.Use(static.Serve("/", static.LocalFile(s.path, false)))
	// e.Static("/dist", s.path)
}

func Init() {

	conf := config.ServerConfig()

	chatbot := bot.NewChatbot(conf.Logic.Email, conf.Logic.Password, conf.Logic.Proxy)
	util.PanicWhenError(container.Singleton(func() *bot.Chatbot {
		return chatbot
	}))
	auth, err := auth.NewAuthHandler()
	if err != nil {
		panic(err)
	}
	api.RegisterHttpRouter(auth)
	api.RegisterHttpRouter(&staticFileHandler{
		path: "/root/code/server-bak/frontend/dist",
	})

	conversation.Register(auth.Middleware())
	chat.Register(auth.Middleware())
}
