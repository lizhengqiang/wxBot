package bot

func (this *WeixinBot) contactMessage(msg *AddMsg) {

	if mp := this.FindMp(msg.FromUserName); mp != nil {
		this.mpMessage(msg, mp)
		return
	}
	if this.getProperty("tuling.contact") == TRUE {
		r, err := this.callTuling(msg.Content, msg.FromUserName)
		if err != nil {
			return
		}
		this.SendMsg(r, msg.FromUserName)
	}

	this.Println(msg.FromUserName, msg.Content)
}
