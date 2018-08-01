package constatns

import "time"

const (
	WechatMsgSendAPI        = "https://qyapi.weixin.qq.com/cgi-bin/message/send"
	WechatUploadMediaAPI    = "https://qyapi.weixin.qq.com/cgi-bin/media/upload"
	WechatTokenGetAPI       = "https://qyapi.weixin.qq.com/cgi-bin/gettoken"
	WechatUserGetAPI        = "https://qyapi.weixin.qq.com/cgi-bin/user/get"
	WehcatPartyMemberGetAPI = "https://qyapi.weixin.qq.com/cgi-bin/user/list"
	MobSMSCallAPI           = "https://webapi.sms.mob.com/sms/voice"

	WechatMsgTypeText   = "text"
	WechatMsgTypeImage  = "image"
	WechatMsgTypeMPNews = "news"
)

var (
	WechatAccessTokenExpireTime time.Duration
	HttpRequestTimeOut          time.Duration
)

func init() {
	WechatAccessTokenExpireTime, _ = time.ParseDuration("7000s")
	HttpRequestTimeOut, _ = time.ParseDuration("5s")
}
