package service

import (
	"encoding/xml"
	"fmt"

	wechatapi "../common/utils/wechatapi"
	wxbizmsgcrypt "../common/utils/wxbizmsgcrypt"
	"../dto"
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

func NewBaseService(token, sEncodingAESKey, sCorpID, sCorpSecret, agentID string) (*BaseService, error) {
	statuscode, msgCrypt := wxbizmsgcrypt.GenerateWXBizMsgCrypt(token, sEncodingAESKey, sCorpID)
	if statuscode != wxbizmsgcrypt.WXBizMsgCrypt_OK {
		return nil, fmt.Errorf("fail to wxbizmsgcrypt")
	}
	wechatAPI := wechatapi.New(sCorpID, sCorpSecret, agentID)
	return &BaseService{
		msgCrypt:  msgCrypt,
		wechatAPI: wechatAPI,
	}, nil
}

func (s *BaseService) EncryptMsg(msg *dto.WechatReplyMsg, nonce string) ([]byte, error) {
	statusCode, encodeBytes := s.msgCrypt.EncryptMsg(msg.String(), nonce)
	if statusCode != wxbizmsgcrypt.WXBizMsgCrypt_OK {
		return nil, fmt.Errorf("EncryptMsg Error:%d", statusCode)
	}
	return encodeBytes, nil
}

func (s *BaseService) DecryptMsg(msg wxbizmsgcrypt.ReviceMsg, sMsgSignature, sTimeStamp, sNonce string) (*dto.WXBizMsg, error) {
	statusCode, decrptMsg := s.msgCrypt.DecryptMsg(msg, sMsgSignature, sTimeStamp, sNonce)
	if statusCode != wxbizmsgcrypt.WXBizMsgCrypt_OK {
		return nil, fmt.Errorf("DecryptMsg Fail :%d", statusCode)
	}
	decryptedMsg := dto.WXBizMsg{}
	err := xml.Unmarshal(decrptMsg, &decryptedMsg)
	return &decryptedMsg, err
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
