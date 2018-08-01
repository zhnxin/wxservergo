package main

import (
	"log"

	"github.com/kataras/iris"

	"./web"
	"github.com/kataras/iris/mvc"

	"./service"
	"./settings"
)

var (
	// CONFIG config obj for server
	CONFIG *settings.Config
	Logger *log.Logger
)

func init() {
	config, err := settings.GetConfig("config.toml")
	if err != nil {
		panic(err)
	}
	CONFIG = &config
	Logger = settings.GetLogger(CONFIG)
}

func main() {
	demoService, err := service.NewDemoService(CONFIG.Wechat.Token, CONFIG.Wechat.EncodingAESKey,
		CONFIG.Wechat.CorpID, CONFIG.Wechat.Secret, CONFIG.Wechat.AgentID)
	if err != nil {
		Logger.Fatal(err.Error())
	}
	app := iris.New()
	mvc.New(app.Party("/")).Handle(&web.DemoController{Service: demoService})
	app.Run(iris.Addr(CONFIG.Server), iris.WithCharset("UTF-8"), iris.WithoutVersionChecker)
}
