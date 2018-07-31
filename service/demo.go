package service

import "../dto"

type DemoService interface {
	URLVerify(string, string, string, string) (string, error)
	MessageHandler() (string, error)
	GetUser(string) ([]dto.UserInfo, error)
	GetEmailList(string) (emailList []string, err error)
	GetPhoneList(string) (phoneList []string, err error)
}
