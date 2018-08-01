package service

import (
	"fmt"

	wechatapi "../common/utils/wechatapi"
	wxbizmsgcrypt "../common/utils/wxbizmsgcrypt"
	"../dto/wechatapiget"
)

type WechatService interface {
	URLVerify(string, string, string, string) (string, error)
	MessageHandler() (string, error)
}

type BaseService struct {
	msgCrypt  *wxbizmsgcrypt.WXBizMsgCrypt
	wechatAPI *wechatapi.WechatAPI
}

func (s *BaseService) URLVerify(msgSignature, timestamp, nonce, echostr string) ([]byte, error) {
	errCode, msg := s.msgCrypt.VerifyURL(msgSignature, timestamp, nonce, echostr)
	if errCode != wxbizmsgcrypt.WXBizMsgCrypt_OK {
		return nil, fmt.Errorf("VerifyURL fail, error code:%d", errCode)
	}
	return msg, nil

}

func (s *BaseService) GetUser(partyID string) ([]wechatapiget.UserInfo, error) {
	userinfoList, err := s.wechatAPI.GetUserList(partyID)
	if err != nil {
		return nil, err
	}
	return userinfoList.GetValue(), nil
}

func (s *BaseService) GetEmailList(partyID string) (emailList []string, err error) {
	emailList, err = s.wechatAPI.GetEmailList(partyID)
	return
}

func (s *BaseService) GetPhoneList(partyID string) (phoneList []string, err error) {
	phoneList, err = s.wechatAPI.GetPhoneList(partyID)
	return
}
