package bot

func (bot *WeixinBot) contactMessage(msg *AddMsg) {

	if mp := bot.FindMp(msg.FromUserName); mp != nil {
		bot.mpMessage(msg, mp)
		return
	}
	if bot.getProperty("tuling.contact") == TRUE {
		r, err := bot.callTuling(msg.Content, msg.FromUserName)
		if err != nil {
			return
		}
		bot.SendMsg(r, msg.FromUserName)
	}

	bot.Println(msg.FromUserName, msg.Content)
}
