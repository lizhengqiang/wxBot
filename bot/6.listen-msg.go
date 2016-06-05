package bot

import (
	"qiniupkg.com/x/errors.v7"
)

func init() {

}
func (this *WeixinBot) Stop() {
	this.Set(IsRunning, FALSE)
}

var (
	ErrStopHandleMsg error = errors.New("停止处理消息")
)

func (this *WeixinBot) HandleMsg() (err error) {
	retcode, selector, err := this.SyncCheck()
	if err != nil {
		return
	}

	if retcode == 1100 {
		this.Set(IsRunning, FALSE)
		this.log("# 你退出了,债见~_~")
		return ErrStopHandleMsg
	}

	if retcode != 0 {
		this.Println(retcode)
		return ErrStopHandleMsg
	}

	if selector == 0 {
		this.Idle()
		return
	}

	if selector == 2 {
		msgList := this.WebWeixinSync()
		if msgList.AddMsgList != nil && len(msgList.AddMsgList) > 0 {
			this.handleMsg(msgList.AddMsgList)
		}
		return
	}

	if selector == 7 {
		this.WebWeixinSync()
		this.log("# 发现你玩手机了!")
		this.Idle()
		return
	}

	if selector == 4 || selector == 6 {
		this.WebWeixinSync()
		this.Idle()
		return
	}

	this.Println("Unknown Selector", selector)
	this.Idle()
	return
}

func (this *WeixinBot) ListenMsgMode() {
	for this.IsRunning() {

	}
}
