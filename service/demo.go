package service

import (
	"../dto"
	"../dto/wechatapiget"
)

type DemoService interface {
	URLVerify(string, string, string, string) (string, error)
	MessageHandler(*dto.WXBizMsg) (*dto.WechatReplyMsg, error)
	ReloadPlugin()
	GetUser(string) ([]wechatapiget.UserInfo, error)
	GetEmailList(string) (emailList []string, err error)
	GetPhoneList(string) (phoneList []string, err error)
}

type demoService struct {
	baseService *BaseService
	actionP     *Ac
}
