package wechatapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	wechatapiget "../../../dto/wechatapiget"
	wechatapipush "../../../dto/wechatapipush"
	settings "../../../settings"
	constatns "../../constatns"
	lrucache "../lrucache"
)

//WechatAPI wechat api proxy
type WechatAPI struct {
	cache       *lrucache.CacheManager
	accessToken string
	corpID      string
	corpSecret  string
	agentID     string
	updateLoack *sync.RWMutex
	client      http.Client
	updateTime  time.Time
}

func New(corpID, corpSecret, agentID string) *WechatAPI {
	updateLoack := new(sync.RWMutex)
	cache := lrucache.GetCacheManager()
	client := http.Client{}
	api := WechatAPI{
		cache:       cache,
		corpID:      corpID,
		corpSecret:  corpSecret,
		agentID:     agentID,
		updateLoack: updateLoack,
		client:      client,
		updateTime:  time.Time{},
		accessToken: "",
	}
	api.updateAccessToken()
	return &api

}

func (api *WechatAPI) updateAccessToken() error {
	api.updateLoack.Lock()
	defer api.updateLoack.Unlock()
	req, err := http.NewRequest("GET", constatns.WechatTokenGetAPI, nil)
	if err != nil {
		return err
	}
	urlQuery := req.URL.Query()
	urlQuery.Add("corpid", api.corpID)
	urlQuery.Add("corpsecret", api.corpSecret)
	req.URL.RawQuery = urlQuery.Encode()
	res, err := api.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	jsonBody, _ := ioutil.ReadAll(res.Body)
	accessToken := wechatapiget.Token{}
	err = json.Unmarshal(jsonBody, &accessToken)
	if err != nil {
		return err
	}
	api.accessToken = accessToken.AccessToken
	api.updateTime = time.Now()
	return nil
}

func (api *WechatAPI) getAccessToken() (string, error) {

	api.updateLoack.RLock()
	defer api.updateLoack.RUnlock()
	if api.updateTime.Add(constatns.WechatAccessTokenExpireTime).After(time.Now()) {
		api.updateLoack.RUnlock()
		err := api.updateAccessToken()
		api.updateLoack.RLock()
		if err != nil {
			return "", err
		}

	}
	return api.accessToken, nil
}

func (api *WechatAPI) UploadImage(imagePath string) (mediaid string, err error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return
	}
	defer file.Close()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("image", filepath.Base(imagePath))
	if err != nil {
		return
	}
	_, err = io.Copy(part, file)
	err = writer.Close()
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", constatns.WechatUploadMediaAPI, body)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	urlQuery := req.URL.Query()
	token, err := api.getAccessToken()
	if err != nil {
		return
	}
	urlQuery.Add("access_token", token)
	urlQuery.Add("type", "image")

	req.URL.RawQuery = urlQuery.Encode()
	res, err := api.client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	jsonbody, _ := ioutil.ReadAll(res.Body)
	media := wechatapiget.MediaUpload{}
	err = json.Unmarshal(jsonbody, &media)
	if err != nil {
		return
	}
	if media.MediaID == "" {
		fmt.Println(media.ErrCode, media.ErrMgs)
	}
	mediaid = media.MediaID
	return
}

func (api *WechatAPI) SendMsg(msg *wechatapipush.WechatAPIMsg) (err error) {
	msg.SetAgentID(api.agentID)
	jsonBody := new(bytes.Buffer)
	err = json.NewEncoder(jsonBody).Encode(msg)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", constatns.WechatMsgSendAPI, jsonBody)
	if err != nil {
		return
	}
	urlQuery := req.URL.Query()
	token, err := api.getAccessToken()
	if err != nil {
		return
	}
	urlQuery.Add("access_token", token)
	req.URL.RawQuery = urlQuery.Encode()
	res, err := api.client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	jsonbody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	var reply wechatapiget.SendMsgReply
	err = json.Unmarshal(jsonbody, &reply)
	if err != nil {
		return err
	}
	settings.GetLogger(nil).Println(reply)
	if reply.ErrCode != 0 {
		return fmt.Errorf("Send Msg Fail:%s", reply.ErrMsg)
	}
	return
}

func (api *WechatAPI) getUser(partyID string) (userInfo map[string]*json.RawMessage, err error) {
	req, err := http.NewRequest("GET", constatns.WehcatPartyMemberGetAPI, nil)
	if err != nil {
		return
	}
	urlQuery := req.URL.Query()
	token, err := api.getAccessToken()
	if err != nil {
		return
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	urlQuery.Add("access_token", token)
	urlQuery.Add("department_id", partyID)
	urlQuery.Add("fetch_child", "0")
	req.URL.RawQuery = urlQuery.Encode()
	res, err := api.client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	jsonbody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(jsonbody, &userInfo)
	return
}

func (api *WechatAPI) GetUserList(partyID string) (*wechatapiget.UserInfoList, error) {
	userInfoList := wechatapiget.NewUserInfoList(partyID)
	isGot := api.cache.Get(&userInfoList)
	if isGot {
		return &userInfoList, nil
	}
	userinfo, err := api.getUser(partyID)
	if err != nil {
		return nil, err
	}
	var userList []wechatapiget.UserInfo
	err = json.Unmarshal(*userinfo["userlist"], &userList)
	if err != nil {
		return nil, err
	}
	userInfoList.SetValue(userList)
	err = api.cache.Set(&userInfoList)
	return &userInfoList, err
}

func (api *WechatAPI) GetPhoneList(partyID string) ([]string, error) {
	userinfoList, err := api.GetUserList(partyID)
	if err != nil {
		return nil, err
	}
	return userinfoList.GetPhoneList(), err
}

func (api *WechatAPI) GetEmailList(partyID string) ([]string, error) {
	userinfoList, err := api.GetUserList(partyID)
	if err != nil {
		return nil, err
	}
	return userinfoList.GetEmailList(), err
}
