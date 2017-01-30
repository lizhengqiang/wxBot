package bot

import "strconv"

// 给公众号群发消息
func (bot *WeixinBot) contacts(content string) {
	if bot.Get("task") == TRUE {
		bot.fileHelperResponse("请等待上次任务结束")
		return
	}
	bot.Set("task", TRUE)
	contacts := []*Contact{}
	bot.unmarshal(contactList, &contacts)

	bot.fileHelperResponse("将要推送:" + strconv.Itoa(len(contacts)) + "个联系人")
	bot.Idle()
	for i, mp := range contacts {
		bot.Println(mp)
		bot.Set("status", "推送联系人:" + strconv.Itoa(i) + "/" + strconv.Itoa(len(contacts)))
		bot.SendMsg(content, mp.UserName)

		bot.Idle()
	}
	bot.Set("task", FALSE)
	bot.fileHelperResponse("推送完毕:" + strconv.Itoa(len(contacts)) + "个联系人")
}
