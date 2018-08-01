package dto

import "encoding/xml"

type WXBizMsg struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string
	FromUserName string
	CreateTime   string
	MsgType      string
	AgentID      string
	//click
	Event    string `xml:"Event,omitempty"`
	EventKey string `xml:"EventKey,omitempty"`
	//text
	Content string `xml:"Content,omitempty"`
	MsgId   string `xml:"MsgId,omitempty"`
	//picture
	MediaId string `xml:"MediaId,omitempty"`
}
type MpnewsReplyArtical struct {
	XMLName     xml.Name `xml:"item"`
	Title       string
	Description string
	Url         string
	PicUrl      string
}

type WechatReplyMsgImage struct {
	XMLName xml.Name `xml:"Image"`
	MediaID string   `xml:"MediaId,omitempty"`
}
type WechatReplyMsg struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string
	FromUserName string
	CreateTime   string
	MsgType      string
	Content      string `xml:"Content,omitempty"`
	Image        WechatReplyMsgImage
	ArticleCount int                  `xml:"ArticleCount,omitempty"`
	Articles     []MpnewsReplyArtical `xml:"Articles,omitempty"`
}

func (msg WechatReplyMsg) String() string {
	data, _ := xml.Marshal(msg)
	return string(data)
}
