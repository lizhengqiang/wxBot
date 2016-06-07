package bot

import (
	"wxBot/provider"
)

func (this *Trigger) Router() {
	this.When("start", func(t *Trigger, b *WeixinBot, m *provider.Message) {
		// 第一次启动
		if !b.IsRunning() {
			// 初始化
			b.Init()
			t.Send(b.ID, "handleQRCodeScanStatus", nil)
		}

	})
	this.When("handleQRCodeScanStatus", func(t *Trigger, b *WeixinBot, m *provider.Message) {
		err := b.HandleQRCodeScanStatus()
		if err == WaitScanQRCode {
			t.Send(b.ID, "handleQRCodeScanStatus", nil)
			return
		}

		if err != nil {
			b.Set(IsRunning, FALSE)
			return
		}

		t.Send(b.ID, "login", nil)
		return

	})

	this.When("login", func(t *Trigger, b *WeixinBot, m *provider.Message) {
		err := b.Login()
		if err != nil {
			b.Println("登陆失败. ")
			b.Set(IsRunning, FALSE)
			return
		}

		// 初始化信息
		b.InitWebWeixin()

		b.Println("初始化信息完毕-1")

		b.WebWeixinStatusNotify()

		b.Println("初始化信息完毕-2")

		b.Set("task", FALSE)

		t.Send(b.ID, "contact", nil)

		t.Send(b.ID, "handleMsg", nil)
		t.Send(b.ID, "handleMsg", nil)
		t.Send(b.ID, "handleMsg", nil)
		t.Send(b.ID, "handleMsg", nil)
		t.Send(b.ID, "handleMsg", nil)

		return
	})

	this.When("contact", func(t *Trigger, b *WeixinBot, m *provider.Message) {
		b.GetContact()

		return
	})

	this.When("handleMsg", func(t *Trigger, b *WeixinBot, m *provider.Message) {
		t.Send(b.ID, "handleMsg", nil)
		err := b.HandleMsg()
		if err != nil {
			return
		}
		return
	})

	this.When("sendMsg", func(t *Trigger, b *WeixinBot, m *provider.Message) {
		b.DoSendMsg(m.Body.(map[string]interface{})["content"].(string), m.Body.(map[string]interface{})["toUserName"].(string))
		return
	})

	this.When("tuling", func(t *Trigger, b *WeixinBot, m *provider.Message) {
		b.Tuling(m.Body.(string))
		return
	})
}
