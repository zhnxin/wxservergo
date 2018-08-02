package web

import (
	"fmt"
	"strings"

	"../common/utils/wxbizmsgcrypt"
	"../service"
	"../settings"
	"github.com/kataras/iris"
)

type textSendDto struct {
	Subject string `json:"subject"`
	Content string `json:"content"`
}

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

func (c *DemoController) Get() error {
	msgSignature, timestamp, nonce := c.getQueryArg()
	echostr := c.Ctx.URLParam("echostr")
	if msgSignature == "" || timestamp == "" || nonce == "" || echostr == "" {
		return fmt.Errorf("lack of query arguments")
	}
	res, err := c.Service.URLVerify(msgSignature, timestamp, nonce, echostr)
	if err != nil {
		return err
	}
	_, err = c.Ctx.Write(res)
	return err
}
func (c *DemoController) Post() error {
	msgSignature, timestamp, nonce := c.getQueryArg()
	if msgSignature == "" || timestamp == "" || nonce == "" {
		return fmt.Errorf("lack of query arguments")
	}
	reciveMsg := wxbizmsgcrypt.ReviceMsg{}
	err := c.Ctx.ReadXML(&reciveMsg)
	if err != nil {
		return err
	}
	msg, err := c.Service.DecryptMsg(reciveMsg, msgSignature, timestamp, nonce)
	if err != nil {
		return err
	}
	replyMsg, err := c.Service.MessageHandler(msg)
	if err != nil {
		return err
	}
	encryptBytes, err := c.Service.EncryptMsg(replyMsg, nonce)
	if err != nil {
		return err
	}
	_, err = c.Ctx.Write(encryptBytes)
	return err
}

func (c *DemoController) GetEmail() error {
	partyID := c.Ctx.URLParam("party_id")
	if partyID == "" {
		return fmt.Errorf("parameter not found: party_id")
	}
	emailList, err := c.Service.GetEmailList(partyID)
	if err != nil {
		return err
	}
	_, err = c.Ctx.WriteString(strings.Join(emailList, " "))
	return err
}

func (c *DemoController) PostEmail() error {
	jsonData := &textSendDto{}
	c.Ctx.ReadJSON(jsonData)
	PartyID := c.Ctx.URLParam("party_id")
	toUser, err := c.Service.GetEmailList(PartyID)
	if err != nil {
		return err
	}
	settings.GetLogger(nil).Printf("send email to party:%s->%v", PartyID, toUser)
	c.Service.SendEmail(toUser, jsonData.Subject, jsonData.Content)
	return nil
}

func (c *DemoController) GetPhone() error {
	partyID := c.Ctx.URLParam("party_id")
	if partyID == "" {
		return fmt.Errorf("parameter not found: party_id")
	}
	phoneList, err := c.Service.GetPhoneList(partyID)
	if err != nil {
		return err
	}
	_, err = c.Ctx.WriteString(strings.Join(phoneList, " "))
	return err
}
