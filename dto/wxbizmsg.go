package dto

import "encoding/xml"

type WXBizMsg struct {
	XML          xml.Name `xml:"xml"`
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
	Itiem       xml.Name `xml:"item"`
	Title       string
	Description string
	Url         string
	PicUrl      string
}
type WechatReplyMsg struct {
	XML          xml.Name `xml:"xml"`
	ToUserName   string
	FromUserName string
	CreateTime   string
	MsgType      string
	Content      string               `xml:"Content,omitempty"`
	MediaID      string               `xml:"Image>MediaId,omitempty"`
	ArticleCount int                  `xml:"ArticleCount,omitempty"`
	Articles     []MpnewsReplyArtical `xml:"Articles,omitempty"`
}
