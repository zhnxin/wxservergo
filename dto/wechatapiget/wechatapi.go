package wechatapiget

import (
	"encoding/json"
	"fmt"
)

//Token wechat access token api respone
type Token struct {
	ErrCode     int    `json:"errcode"`
	ErrMgs      string `json:"errmsg"`
	AccessToken string `json:"access_token,omitempty"`
	ExpiresIN   int    `json:"expires_in,omitempty"`
}

type MediaUpload struct {
	ErrCode   int    `json:"errcode"`
	ErrMgs    string `json:"errmsg"`
	Type      string `json:"type"`
	MediaID   string `json:"media_id"`
	CreatedAt string `json:"created_at"`
}

type UserInfo struct {
	UserID     string `json:"userid"`
	Name       string `json:"name"`
	Department []int  `json:"department"`
	Mobile     string `json:"mobile"`
	Email      string `json:"email"`
}

type UserInfoList struct {
	value []UserInfo
	key   string
}

func NewUserInfoList(key string) UserInfoList {
	return UserInfoList{
		value: nil,
		key:   fmt.Sprintf("userinfolist-%s", key),
	}
}

func (userlist *UserInfoList) CacheKey() []byte {
	return []byte(userlist.key)
}

func (userlist *UserInfoList) SetValue(value []UserInfo) {
	userlist.value = value
}

func (userlist *UserInfoList) GetValue() []UserInfo {
	return userlist.value
}

func (userlist *UserInfoList) JSON() (data []byte, err error) {
	data, err = json.Marshal(userlist.value)
	if err != nil {
		err = fmt.Errorf("UserInfoList:Json:%s,%v", string(userlist.key), err)
	}
	return
}

func (userlist *UserInfoList) LoadJSON(data []byte) error {
	var value []UserInfo
	err := json.Unmarshal(data, &value)
	if err == nil {
		userlist.value = value
	}
	return err
}
func (userlist *UserInfoList) CacheExpireTime() int {
	return 3600 * 24
}

func (userlist *UserInfoList) GetPhoneList() []string {
	if userlist.value == nil {
		return nil
	}
	value := make([]string, 0, len(userlist.value))
	for _, info := range userlist.value {
		value = append(value, info.Mobile)
	}
	return value
}

func (userlist *UserInfoList) GetEmailList() []string {
	if userlist.value == nil {
		return nil
	}
	value := make([]string, 0, len(userlist.value))
	for _, info := range userlist.value {
		value = append(value, info.Email)
	}
	return value
}
func (userlist *UserInfoList) GetUserIDList() []string {
	if userlist.value == nil {
		return nil
	}
	value := make([]string, 0, len(userlist.value))
	for _, info := range userlist.value {
		value = append(value, info.UserID)
	}
	return value
}

type SendMsgReply struct {
	ErrCode      int    `json:"errcode"`
	ErrMsg       string `json:"errmsg"`
	Invaliduser  string `json:"invaliduser,omitempty"`
	Invalidparty string `json:"invalidparty,omitempty"`
	Invalidtag   string `json:"invalidtag,omitempty"`
}
