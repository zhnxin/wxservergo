package web

import (
	"fmt"
	"strings"

	"../common/utils/wxbizmsgcrypt"
	"../service"
	"../settings"
	"encoding/json"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
)

type textSendDto struct {
	Subject string `json:"subject,omitempty"`
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

func (c *DemoController) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle("GET", "/public", "Get")
	b.Handle("POST", "/public", "Post")
}

func (c *DemoController) Get() error {
	msgSignature, timestamp, nonce := c.getQueryArg()
	echoStr := c.Ctx.URLParam("echostr")
	if msgSignature == "" || timestamp == "" || nonce == "" || echoStr == "" {
		return fmt.Errorf("lack of query arguments")
	}
	res, err := c.Service.URLVerify(msgSignature, timestamp, nonce, echoStr)
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
	settings.GetLogger(nil).Printf("DemoControler:Post:ReviceMsg:%v",msg)
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
	err := c.Ctx.ReadJSON(jsonData)
	if err != nil {
		return fmt.Errorf("DemoController:PostEmail:%v", err)
	}
	PartyID := c.Ctx.URLParam("party_id")
	toUser, err := c.Service.GetEmailList(PartyID)
	if err != nil {
		return err
	}
	settings.GetLogger(nil).Printf("send email to party:%s->%v\n", PartyID, toUser)
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

func (c *DemoController) GetText() error {
	partyID := c.Ctx.URLParam("party_id")
	if partyID == "" {
		return fmt.Errorf("parameter not found: party_id")
	}
	userinfo, err := c.Service.GetUser(partyID)
	if err != nil {
		return fmt.Errorf("DemoController:GetText:%v", err)
	}
	data, err := json.Marshal(userinfo)
	if err != nil {
		return fmt.Errorf("DemoController:GetText:%v", err)
	}
	_, err = c.Ctx.Write(data)
	return err
}

func (c *DemoController) PostText() error {
	jsonData := &textSendDto{}
	err := c.Ctx.ReadJSON(jsonData)
	if err != nil {
		return fmt.Errorf("DemoController:PostText:%v", err)
	}
	partyID := c.Ctx.URLParam("party_id")
	if partyID == "" {
		return fmt.Errorf("parameter not found: party_id")
	}
	settings.GetLogger(nil).Println("Democontroller:send msg to ",partyID)
	c.Service.SendWechatText([]string{partyID}, jsonData.Content)
	return nil

}

func (c *DemoController) GetPlugin() {
	settings.GetLogger(nil).Println("DemoController: reflesh ")
	c.Service.ReloadPlugin()
}
