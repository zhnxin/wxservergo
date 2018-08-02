package main

import (
	"log"

	"github.com/kataras/iris"

	"./common/utils/email"
	"./web"
	"github.com/kataras/iris/mvc"

	"./service"
	"./settings"
)

var (
	// CONFIG config obj for server
	CONFIG      *settings.Config
	Logger      *log.Logger
	EmailClient *email.Client
)

func init() {
	var err error
	CONFIG, err = settings.GetConfig("config.toml")
	if err != nil {
		panic(err)
	}
	Logger = settings.GetLogger(CONFIG)
	EmailClient = email.GenerateEmailClientSingleton(
		CONFIG.Email.User, CONFIG.Email.Password, CONFIG.Email.NickName,
		CONFIG.Email.Host, CONFIG.Email.Port, CONFIG.Email.IsSSL)
}

func main() {
	demoService, err := service.NewDemoService(CONFIG.Wechat.Token, CONFIG.Wechat.EncodingAESKey,
		CONFIG.Wechat.CorpID, CONFIG.Wechat.Secret, CONFIG.Wechat.AgentID, EmailClient)
	if err != nil {
		Logger.Fatal(err.Error())
	}
	app := iris.New()
	app.Use(func(ctx iris.Context) {
		ctx.Application().Logger().Infof("request for path:%s %s", ctx.Method(), ctx.Path())
		ctx.Next()
	})
	mvc.New(app.Party("/")).Handle(&web.DemoController{Service: demoService})
	app.Run(iris.Addr(CONFIG.Server), iris.WithCharset("UTF-8"), iris.WithoutVersionChecker)
}
