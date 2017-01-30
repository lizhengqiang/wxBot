package bot

import "strings"

func (bot *WeixinBot) groupMessage(msg *AddMsg) {
	contents := strings.Split(msg.Content, `:<br/>`)
	var content string
	var userName string
	if len(contents) == 1 {
		content = contents[0]
	} else if len(contents) == 2 {
		userName = contents[0]
		content = contents[1]
	} else {
		content = msg.Content
	}
	me := bot.GetMe()
	bot.WebWeixinStatusNotify(me.UserName, msg.FromUserName, 1)

	if userName == bot.GetMe().UserName {
		return
	}

	if bot.getProperty("tuling.group") == TRUE {
		r, err := bot.callTuling(msg.Content, msg.FromUserName)
		if err != nil {
			return
		}
		bot.SendMsg(r, msg.FromUserName)
	}
	bot.Hear(msg.FromUserName, userName, content)

	bot.Println(msg.FromUserName, userName, content)

}
