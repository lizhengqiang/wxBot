package bot

import (
	"qiniupkg.com/x/errors.v7"
)

func (this *WeixinBot) Stop() {
	this.set(isRunning, FALSE)
}

var ErrBotLogining error = errors.New("Bot正在登录中...")

// 开启这个会话, 这货只能运行一次
func (this *WeixinBot) Start() (err error) {

	// 第一次启动
	if !this.IsRunning() {
		// 初始化
		this.Init()
		// 扫码登陆
		err := this.scanForLogin()
		this.set(isLogining, FALSE)
		if err != nil {
			// 登录失败
			this.set(isRunning, FALSE)

			return err
		}

	}

	// 正在登陆
	if this.IsLogining() && this.IsRunning() {
		return ErrBotLogining
	}

	// 依旧在运行,并且已经登陆过了

	// 获取联系人列表
	go this.GetContact()

	this.set("task", FALSE)

	// 开始监听消息
	this.ListenMsgMode()

	return nil
}

func (this *WeixinBot) ListenMsgMode() {
	for this.IsRunning() {
		retcode, selector, err := this.SyncCheck()
		if err != nil {
			this.Idle()
			continue
		}

		if retcode == 1100 {
			this.log("# 你退出了,债见~_~")
			return
		}

		if retcode != 0 {
			this.Println(retcode)
			this.Idle()
			continue
		}

		if selector == 0 {
			continue
		}

		if selector == 2 {
			msgList := this.WebWeixinSync()
			if msgList.AddMsgList != nil && len(msgList.AddMsgList) > 0 {
				this.handleMsg(msgList.AddMsgList)
			}
			continue
		}

		if selector == 7 {
			_ = this.WebWeixinSync()
			this.log("# 发现你玩手机了!")
			this.Idle()
			continue
		}

		if selector == 4 || selector == 6 {
			this.WebWeixinSync()
			this.Idle()
			continue
		}

		this.Println("Unknown Selector", selector)

		this.Idle()

	}
}
