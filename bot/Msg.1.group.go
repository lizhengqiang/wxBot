package bot

import "strings"

func (this *WeixinBot) groupMessage(msg *AddMsg) {
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

	if userName == this.getMe().UserName {
		return
	}

	if this.getProperty("tuling.group") == TRUE {
		r, err := this.callTuling(msg.Content, msg.FromUserName)
		if err != nil {
			return
		}
		this.SendMsg(r, msg.FromUserName)
	}

	this.Println(msg.FromUserName, userName, content)

}

