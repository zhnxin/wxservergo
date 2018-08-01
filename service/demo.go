package service

import (
	"../common/utils/actionplugin"
	"../common/utils/wxbizmsgcrypt"
	"../dto"
	"../dto/wechatapiget"
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
}

type demoService struct {
	baseService *BaseService
	actionP     *actionplugin.ActionPlugin
}

func NewDemoService(token, sEncodingAESKey, sCorpID, sCorpSecret, agentID string) (DemoService, error) {
	baseService, err := NewBaseService(token, sEncodingAESKey, sCorpID, sCorpSecret, agentID)
	if err != nil {
		return nil, err
	}
	actionP := actionplugin.New("demo")
	return &demoService{
		baseService: baseService,
		actionP:     actionP,
	}, nil
}

func (s *demoService) MessageHandler(msg *dto.WXBizMsg) (*dto.WechatReplyMsg, error) {
	fn, err := s.actionP.GetHandler(msg)
	if err != nil {
		return nil, err
	}
	return fn()
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
