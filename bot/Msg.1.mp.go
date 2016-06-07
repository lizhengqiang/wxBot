package bot

func (this *WeixinBot) mpMessage(msg *AddMsg, mp *Contact) {

	if this.getProperty("tuling.mp") == TRUE {
		r, err := this.callTuling(msg.Content, msg.FromUserName)
		if err != nil {
			return
		}
		this.SendMsg(r, msg.FromUserName)
	}

	this.Println(msg.FromUserName, msg.Content)
}
