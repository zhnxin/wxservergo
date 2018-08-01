package web

import (
	"fmt"

	"../common/utils/wxbizmsgcrypt"
	"../service"
	"github.com/kataras/iris"
)

type DemoController struct {
	Service service.DemoService
	Ctx     iris.Context
}

func (c *DemoController) getQueryArg() (msgSignature, timestamp, nonce string) {
	msgSignature = c.Ctx.URLParam("msg_signature")
	timestamp = c.Ctx.URLParam("timestamp")
	nonce = c.Ctx.URLParam("nonce")
	return
}

func (c *DemoController) Get() (string, error) {
	msgSignature, timestamp, nonce := c.getQueryArg()
	echostr := c.Ctx.URLParam("echostr")
	if msgSignature == "" || timestamp == "" || nonce == "" || echostr == "" {
		return "", fmt.Errorf("lack of query arguments")
	}
	res, err := c.Service.URLVerify(msgSignature, timestamp, nonce, echostr)
	return string(res), err
}
func (c *DemoController) Post() (string, error) {
	msgSignature, timestamp, nonce := c.getQueryArg()
	if msgSignature == "" || timestamp == "" || nonce == "" {
		return "", fmt.Errorf("lack of query arguments")
	}
	reciveMsg := wxbizmsgcrypt.ReviceMsg{}
	err := c.Ctx.ReadXML(&reciveMsg)
	if err != nil {
		return "", err
	}
	msg, err := c.Service.DecryptMsg(reciveMsg, msgSignature, timestamp, nonce)
	if err != nil {
		return "", err
	}
	replyMsg, err := c.Service.MessageHandler(msg)
	if err != nil {
		return "", err
	}
	encryptBytes, err := c.Service.EncryptMsg(replyMsg, nonce)
	return string(encryptBytes), err
}
