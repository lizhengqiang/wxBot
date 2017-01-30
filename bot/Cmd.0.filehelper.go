// 主要使用文件传输助手来控制
package bot

import (
	"github.com/cocotyty/json"
	"strings"
)

func (bot *WeixinBot) fileHelperResponse(content string) {
	bot.SendMsg("#Bot#"+content, "filehelper")

}

func (bot *WeixinBot) fileHelper(msg *AddMsg) {
	if msg.Content == "me" {
		bytes, _ := json.MarshalIndent(bot.GetMe(), "", "  ")
		bot.fileHelperResponse(string(bytes))
		return
	}
	if msg.Content == "stop" {
		bot.Set(IsRunning, FALSE)
		bot.IsLoopRunning = false
	}
	if msg.Content == "reload" {
		bot.ReloadJS()
		return
	}
	if msg.Content == "logs" {
		bot.fileHelperResponse(bot.Cacher.Get("logs"))
		return
	}

	if msg.Content == "now" {
		bot.fileHelperResponse(bot.Get("status"))
		return
	}

	if msg.Content == "mps" {
		bot.fileHelperResponse(bot.Get(mpList))
		return
	}

	if strings.Contains(msg.Content, "mps") {
		go bot.mps(strings.Replace(msg.Content, "mps ", "", -1))
		return
	}

	if strings.Contains(msg.Content, "contacts") {
		go bot.contacts(strings.Replace(msg.Content, "contacts ", "", -1))
		return
	}

	return

}
