package service

import (
	"fmt"

	"../common/utils/actionplugin"
	"../common/utils/email"
	"../common/utils/wxbizmsgcrypt"
	"../dto"
	"../dto/wechatapiget"
	"../dto/wechatapipush"
	"../settings"
)

type DemoService interface {
	URLVerify(string, string, string, string) ([]byte, error)
	MessageHandler(*dto.WXBizMsg) (*dto.WechatReplyMsg, error)
	EncryptMsg(*dto.WechatReplyMsg, string) ([]byte, error)
	DecryptMsg(wxbizmsgcrypt.ReviceMsg, string, string, string) (*dto.WXBizMsg, error)
	ReloadPlugin()
	GetUser(string) ([]wechatapiget.UserInfo, error)
	GetEmailList(string) (emailList []string, err error)
	GetPhoneList(string) (phoneList []string, err error)
	SendEmail([]string, string, string) error
	SendWechatText([]string, string) error
}

type demoService struct {
	baseService *baseService
	actionP     *actionplugin.ActionPlugin
	emailClient *email.Client
}

func NewDemoService(token, sEncodingAESKey, sCorpID, sCorpSecret, agentID string, emailClient *email.Client) (DemoService, error) {
	baseService, err := NewBaseService(token, sEncodingAESKey, sCorpID, sCorpSecret, agentID)
	if err != nil {
		return nil, err
	}
	actionP := actionplugin.New("demo")
	return &demoService{
		baseService: baseService,
		actionP:     actionP,
		emailClient: emailClient,
	}, nil
}

func (s *demoService) MessageHandler(msg *dto.WXBizMsg) (*dto.WechatReplyMsg, error) {
	fn, err := s.actionP.GetHandler(msg)
	if err != nil {
		return nil, err
	}
	replyMsg, err := fn(msg)
	if err != nil {
		return nil, err
	}
	return replyMsg, nil
}
func (s *demoService) ReloadPlugin() {
	s.actionP.Load()
}
func (s *demoService) GetUser(partyID string) ([]wechatapiget.UserInfo, error) {
	return s.baseService.GetUser(partyID)
}
func (s *demoService) URLVerify(msgSignature, timestamp, nonce, echostr string) ([]byte, error) {
	return s.baseService.URLVerify(msgSignature, timestamp, nonce, echostr)
}

func (s *demoService) GetEmailList(partyID string) ([]string, error) {
	return s.baseService.GetEmailList(partyID)
}
func (s *demoService) GetPhoneList(partyID string) ([]string, error) {
	return s.baseService.GetPhoneList(partyID)
}

func (s *demoService) EncryptMsg(msg *dto.WechatReplyMsg, nonce string) ([]byte, error) {
	return s.baseService.EncryptMsg(msg, nonce)
}

func (s *demoService) DecryptMsg(msg wxbizmsgcrypt.ReviceMsg, sMsgSignature, sTimeStamp, sNonce string) (*dto.WXBizMsg, error) {
	return s.baseService.DecryptMsg(msg, sMsgSignature, sTimeStamp, sNonce)
}

func (s *demoService) SendEmail(toUser []string, subject string, content string) error {
	if s.emailClient == nil {
		return fmt.Errorf("not email client")
	}
	go func() {
		if err := s.emailClient.SendEmail(toUser, subject, content); err != nil {
			settings.GetLogger(nil).Println(err.Error())
		}
	}()
	return nil
}

func (s *demoService) SendWechatText(toParty []string, content string) error {
	go func() {
		msg := wechatapipush.NewTextMsg(content)
		msg.SetToParty(toParty)
		if err := s.baseService.wechatAPI.SendMsg(msg); err != nil {
			settings.GetLogger(nil).Println(err.Error())
		}
	}()
	return nil

}
