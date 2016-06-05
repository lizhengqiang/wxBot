package mns

import (
	"github.com/cocotyty/summer"
	"wxBot/bot"
)

func (this *Trigger) Router() {
	summer.GetStoneWithName("Trigger").(*Trigger).When("start", func(t *Trigger, b *bot.WeixinBot, m *Message) {
		// 第一次启动
		if !b.IsRunning() {
			// 初始化
			b.Init()
			t.Send(b.ID, "handleQRCodeScanStatus", nil)
		}

	})
	summer.GetStoneWithName("Trigger").(*Trigger).When("handleQRCodeScanStatus", func(t *Trigger, b *bot.WeixinBot, m *Message) {
		err := b.HandleQRCodeScanStatus()
		if err == bot.WaitScanQRCode {
			t.Send(b.ID, "handleQRCodeScanStatus", nil)
			return
		}

		if err != nil {
			b.Set(bot.IsRunning, bot.FALSE)
			return
		}

		t.Send(b.ID, "login", nil)
		return

	})

	summer.GetStoneWithName("Trigger").(*Trigger).When("login", func(t *Trigger, b *bot.WeixinBot, m *Message) {
		err := b.Login()
		if err != nil {
			b.Println("登陆失败. ")
			b.Set(bot.IsRunning, bot.FALSE)
			return
		}

		// 初始化信息
		b.InitWebWeixin()

		b.Println("初始化信息完毕-1")

		b.WebWeixinStatusNotify()

		b.Println("初始化信息完毕-2")

		b.Set("task", bot.FALSE)

		t.Send(b.ID, "contact", nil)

		t.Send(b.ID, "listen", nil)

		return
	})

	summer.GetStoneWithName("Trigger").(*Trigger).When("contact", func(t *Trigger, b *bot.WeixinBot, m *Message) {
		b.GetContact()

		return
	})

	summer.GetStoneWithName("Trigger").(*Trigger).When("handleMsg", func(t *Trigger, b *bot.WeixinBot, m *Message) {
		err := b.HandleMsg()
		if err != nil {
			return
		}
		t.Send(b.ID, "handleMsg", nil)
		return
	})
}
