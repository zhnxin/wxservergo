package main

import (
	"fmt"

	"github.com/buaazp/fasthttprouter"

	"github.com/valyala/fasthttp"

	wxbizmsgcrypt "./common/utils/wxbizmsgcrypt"
	"./settings"
)

var (
	// CONFIG config obj for server
	CONFIG *settings.Config
	// WECAHT manager
	WECAHT *wxbizmsgcrypt.WXBizMsgCrypt
)

func init() {
	config, err := settings.GetConfig("config.toml")
	if err != nil {
		panic(err)
	}
	CONFIG = &config
	errCode, wxcrypt := wxbizmsgcrypt.GenerateWXBizMsgCrypt(CONFIG.Wechat.Token, CONFIG.Wechat.EncodingAESKey, CONFIG.Wechat.CorpID)
	if errCode != wxbizmsgcrypt.WXBizMsgCrypt_OK {
		panic(errCode)
	}
	WECAHT = &wxcrypt
}

func getVerifyArgs(ctx *fasthttp.RequestCtx) (msgSignature, timestamp, nonce string) {
	msgSignature = string(ctx.QueryArgs().Peek("msg_signature"))
	timestamp = string(ctx.QueryArgs().Peek("timestamp"))
	nonce = string(ctx.QueryArgs().Peek("nonce"))
	return
}

func urlVirifyHandler(ctx *fasthttp.RequestCtx) {
	msgSignature, timestamp, nonce := getVerifyArgs(ctx)
	echostr := string(ctx.QueryArgs().Peek("echostr"))
	fmt.Println("msgSignature: ", msgSignature)
	fmt.Println("timestamp: ", timestamp)
	fmt.Println("nonce: ", nonce)
	fmt.Println("echostr: ", echostr)
	errCode, msg := WECAHT.VerifyURL(msgSignature, timestamp, nonce, echostr)
	if errCode != wxbizmsgcrypt.WXBizMsgCrypt_OK {
		fmt.Fprintf(ctx, "errcode:%d", errCode)
	} else {
		fmt.Fprint(ctx, msg)
	}
}

func eventReplyHandler(ctx *fasthttp.RequestCtx) {

}

func main() {
	router := fasthttprouter.New()
	router.GET("/public/api", urlVirifyHandler)
	if err := fasthttp.ListenAndServe(CONFIG.Server, router.Handler); err != nil {
		fmt.Println("start fasthttp fail:", err.Error())
	}
}
