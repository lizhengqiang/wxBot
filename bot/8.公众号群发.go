package bot

import "strconv"

// 给公众号群发消息
func (bot *WeixinBot) mps(content string) {
	if bot.Get("task") == TRUE {
		bot.fileHelperResponse("请等待上次任务结束")
		return
	}
	bot.Set("task", TRUE)
	mps := []*Contact{}
	bot.unmarshal(mpList, &mps)

	bot.fileHelperResponse("将要推送:" + strconv.Itoa(len(mps)) + "个公众号")
	bot.Idle()
	for i, mp := range mps {
		bot.Println(mp)
		bot.Set("status", "推送公众号:" + strconv.Itoa(i) + "/" + strconv.Itoa(len(mps)))
		bot.SendMsg(content, mp.UserName)

		bot.Idle()
	}
	bot.Set("task", FALSE)
	bot.fileHelperResponse("推送完毕:" + strconv.Itoa(len(mps)) + "个公众号")
}
