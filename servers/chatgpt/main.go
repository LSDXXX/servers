package main

import (
	"context"
	"flag"

	commonapi "github.com/LSDXXX/libs/api"
	"github.com/LSDXXX/libs/app"
	"github.com/LSDXXX/libs/config"
	"github.com/LSDXXX/libs/infra"
	"github.com/LSDXXX/libs/pkg/container"
	"github.com/LSDXXX/libs/pkg/log"
	"github.com/LSDXXX/servers/chatgpt/api"
	serverconfig "github.com/LSDXXX/servers/chatgpt/config"
	"github.com/eatmoreapple/openwechat"
	"github.com/sirupsen/logrus"
)

func wechatBot() {
	bot := openwechat.DefaultBot(openwechat.Desktop)
	bot.MessageHandler = func(msg *openwechat.Message) {
		logrus.Debugf("from use: %s", msg.FromUserName)
		if msg.IsText() {
			logrus.Debugf("content: %+v", msg)
		}
	}
	if err := bot.Login(); err != nil {
		panic(err)
	}
	bot.Block()
}

var configPath string

func init() {
	// flag.StringVar(&configPath, "conf", "/data/weiling/conf/logic-engine-worker/main.yaml", "config path")
	flag.StringVar(&configPath, "conf", "./main.yaml", "config path")
}

func main() {
	flag.Parse()
	var conf serverconfig.Config

	err := config.ReadConfig(configPath, &conf)
	if err != nil {
		panic(err)
	}
	serverconfig.SetServerConfig(&conf)
	container.Singleton(func() *serverconfig.Config {
		return &conf
	})
	container.Singleton(func() *config.Config {
		return &conf.Common
	})

	log.InitGlobalLog(&conf.Common.Log)

	err = infra.Init(
		infra.WithDB(),
		// infra.WithRedis(),
	)
	if err != nil {
		panic(err)
	}
	logrus.Info("infra init success")

	err = commonapi.Init("worker")
	if err != nil {
		panic(err)
	}
	logrus.Info("common api build success")

	api.Init()

	if err = app.Start(context.Background()); err != nil {
		panic(err)
	}
}
