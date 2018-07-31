package wechatapipush

import (
	"strings"

	"../../common/constatns"
)

type mpnewArticle struct {
	Title   string `json:"title"`
	MediaID string `json:"thumb_media_id"`
	Content string `json:"content"`
	Digest  string `json:"digest"`
}

type textAttr struct {
	Content string `json:"content"`
}

type imageAttr struct {
	MediaID string `json:"media_id"`
}

type newAttr struct {
	Articles []mpnewArticle `json:"articles"`
}
type WechatAPIMsg struct {
	ToUser  string    `json:"touser,omitempty"`
	ToParty string    `json:"toparty,omitempty"`
	Totag   string    `json:"totag,omitempty"`
	Msgtype string    `json:"msgtype"`
	Agentid string    `json:"agentid"`
	Text    textAttr  `json:"text,omitempty"`
	Image   imageAttr `json:"image,omitempty"`
	News    newAttr   `json:"news,omitempty"`
}

func NewTextMsg(content string) WechatAPIMsg {
	return WechatAPIMsg{
		Msgtype: constatns.WechatMsgTypeText,
		Text:    textAttr{content},
	}
}

func NewImageMsg(mediaID string) WechatAPIMsg {
	return WechatAPIMsg{
		Msgtype: constatns.WechatMsgTypeImage,
		Image:   imageAttr{mediaID},
	}
}

func (msg *WechatAPIMsg) SetToUser(userIDList []string) {
	touser := strings.Join(userIDList, "|")
	msg.ToUser = touser
}
func (msg *WechatAPIMsg) SetToParty(partyIDList []string) {
	toparty := strings.Join(partyIDList, "|")
	msg.ToParty = toparty
}
func (msg *WechatAPIMsg) SetToTag(tagList []string) {
	totag := strings.Join(tagList, "|")
	msg.Totag = totag
}

func (msg *WechatAPIMsg) SetAgentID(agentID string) {
	msg.Agentid = agentID
}
