package bot

func (bot *WeixinBot) mpMessage(msg *AddMsg, mp *Contact) {

	if bot.getProperty("tuling.mp") == TRUE {
		r, err := bot.callTuling(msg.Content, msg.FromUserName)
		if err != nil {
			return
		}
		bot.SendMsg(r, msg.FromUserName)
	}

	bot.Println(msg.FromUserName, msg.Content)
}
