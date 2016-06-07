// 主要使用文件传输助手来控制
package bot

import (
	"github.com/cocotyty/json"
	"github.com/cocotyty/summer"
	"strings"
)

func (this *WeixinBot) fileHelperResponse(content string) {
	this.SendMsg("#Bot#" + content, "filehelper")

}

func (this *WeixinBot) fileHelper(msg *AddMsg) {
	if msg.Content == "me" {
		bytes, _ := json.MarshalIndent(this.GetMe(), "", "  ")
		this.fileHelperResponse(string(bytes))
		return
	}

	if msg.Content == "logs" {
		this.fileHelperResponse(this.Cacher.Get("logs"))
		return
	}

	if msg.Content == "now" {
		this.fileHelperResponse(this.Get("status"))
		return
	}

	if msg.Content == "mps" {
		this.fileHelperResponse(this.Get(mpList))
		return
	}

	if strings.Contains(msg.Content, "mps") {
		go this.mps(strings.Replace(msg.Content, "mps ", "", -1))
		return
	}

	// 转发消息
	args := strings.Split(msg.Content, " ")
	summer.GetStoneWithName("Trigger").(*Trigger).Send(this.ID, args[0], msg.Content)
	return

}
