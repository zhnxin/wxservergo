package main

import (
	"fmt"
	"math/rand"
	"time"

	"../common/utils/actionplugin"
	"../dto"
)

var ClickEventFuncMap map[string]actionplugin.HandlerFunc

func GetHandler(msg *dto.WXBizMsg) (actionplugin.HandlerFunc, error) {
	switch msg.MsgType {
	case "event":
		{
			if msg.Event == "click" {
				fn, ok := ClickEventFuncMap[msg.EventKey]
				if ok {
					return fn, nil
				}
				return nil, fmt.Errorf("click eventKey not found")
			}
		}
	default:
		{
			return nil, fmt.Errorf("not match")
		}

	}
	return nil, fmt.Errorf("nothing match")
}

func init() {
	ClickEventFuncMap = make(map[string]actionplugin.HandlerFunc)
	ClickEventFuncMap["10001"] = func(msg *dto.WXBizMsg) (*dto.WechatReplyMsg, error) {
		emoticons := []string{"富强", "民主", "文明", "和谐", "自由", "平等", "公正", "法治", "爱国", "敬业",
			"诚信", "友善", "( •̀ .̫ •́ )✧", "(つд⊂)", " (•౪• )",
			" (๑•̀ㅂ•́) ✧", "ლ(╹◡╹ლ)", "_(:з」∠)_", "( •̥́ ˍ •̀ू )", "Ծ‸Ծ", "～﹃～)~zZ",
			"(⺣◡⺣)", "(๑•́ ∀ •̀๑)", "(ง •̀_•́)ง", "ฅ^ω^ฅ", "( ´°̥̥̥̥̥̥̥̥ω°̥̥̥̥̥̥̥̥`)",
			"(⇀‸↼‶)", "(⃘  ̂͘₎o̮₍ ̂͘ )⃘", "ฅʕ•̫͡•ʔฅ", "(๑॔ᵒ̴̶̷◡  ˂̶๑॓)ゞ❣", "ᕙ(⇀‸↼‵‵)ᕗ",
			"(ᵒ̤̑ ₀̑ ᵒ̤̑)", "ヽ(ｏ`皿′ｏ)ﾉ", "ー( ´ ▽ ` )ﾉ", "↺  ♫   ☼",
			"(*´・ω・`)⊃", "⊂(˃̶͈̀ε ˂̶͈́ ⊂ )))Σ≡=─", "_(´ཀ`」 ∠)_", "( ⸝⸝⸝°_°⸝⸝⸝ )", "∠( ᐛ 」∠)＿"}
		rand.Seed(time.Now().UTC().UnixNano())
		replyMsg := dto.NewTextWechatReplyMsg(emoticons[rand.Intn(len(emoticons))])

		replyMsg.ToUserName = msg.FromUserName
		replyMsg.FromUserName = msg.ToUserName
		return replyMsg, nil
	}
}
func main() {}
