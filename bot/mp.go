package bot

import "strconv"

// 给公众号群发消息
func (this *WeixinBot) mps(content string) {
	if this.get("task") == TRUE {
		this.fileHelperResponse("请等待上次任务结束")
		return
	}
	this.set("task", TRUE)
	mps := []*Contact{}
	this.unmarshal(mpList, &mps)

	this.fileHelperResponse("将要推送:" + strconv.Itoa(len(mps)) + "个公众号")
	this.Idle()
	for i, mp := range mps {
		this.Println(mp)
		this.set("status", "推送公众号:" + strconv.Itoa(i) + "/" + strconv.Itoa(len(mps)))
		this.SendMsg(content, mp.UserName)

		this.Idle()
	}
	this.set("task", FALSE)
	this.fileHelperResponse("推送完毕:" + strconv.Itoa(len(mps)) + "个公众号")
}
