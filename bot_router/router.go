package bot_router

import (
	"github.com/cocotyty/summer"
	"github.com/lizhengqiang/wxBot/bot"
	"github.com/lizhengqiang/wxBot/domain"
)

type Router struct {
	Trigger *Trigger `sm:"*"`
	running bool
	Log     *summer.SimpleLog
}

func init() {
	summer.Put(&Router{})
}

func (r *Router) Init() {
	r.Log = summer.NewSimpleLog("BotRouter", summer.InfoLevel)
}
func (r *Router) Ready() {
	trigger := r.Trigger
	trigger.When("start", func(t *Trigger, b *bot.WeixinBot, m *domain.Message) (err error) {
		r.Log.Info("开始登陆")
		// 第一次启动
		if !b.IsRunning() {
			// 初始化
			b.Set(bot.IsRunning, bot.TRUE)
			b.Set(bot.IsSigning, bot.TRUE)
			b.Init()
			t.Send(b.ID, "handleQRCodeScanStatus", nil)
		}
		return nil
	})
	trigger.When("handleQRCodeScanStatus", func(t *Trigger, b *bot.WeixinBot, m *domain.Message) (err error) {

		if b.IsSigned() {
			return nil // 删除消息
		}
		r.Log.Info("等待扫描二维码")
		err = b.HandleQRCodeScanStatus()
		// 继续等待
		if err == bot.WaitScanQRCode {
			t.Send(b.ID, "handleQRCodeScanStatus", nil)
			return nil // 删除消息
		}

		// 登录失败
		if err != nil {
			r.Log.Info("二维码扫描失败")
			b.Set(bot.IsRunning, bot.FALSE)
			b.Set(bot.IsSigning, bot.FALSE)
			return nil // 删除消息
		}
		r.Log.Info("二维码扫描成功")
		t.Send(b.ID, "login", nil)
		return

	})

	trigger.When("login", func(t *Trigger, b *bot.WeixinBot, m *domain.Message) (err error) {
		err = b.Login()
		if err != nil {
			b.Println("登陆失败. ")
			b.Set(bot.IsRunning, bot.FALSE)
			b.Set(bot.IsSigning, bot.FALSE)
			return
		}

		// 初始化信息
		b.InitWebWeixin()

		b.Println("初始化信息完毕-1")

		me := b.GetMe()
		b.WebWeixinStatusNotify(me.UserName, me.UserName, 3)

		b.Println("初始化信息完毕-2")

		b.Set("task", bot.FALSE)
		b.Set(bot.IsSigned, bot.TRUE)

		b.ReloadJS()

		b.ListenMode()

		return
	})

	trigger.When("contact", func(t *Trigger, b *bot.WeixinBot, m *domain.Message) (err error) {
		b.GetContact()
		return
	})

	trigger.When("handleMsg", func(t *Trigger, b *bot.WeixinBot, m *domain.Message) (err error) {
		b.ListenMode()
		return
	})

	trigger.When("sendMsg", func(t *Trigger, b *bot.WeixinBot, m *domain.Message) (err error) {
		b.DoSendMsg(m.Body.(map[string]interface{})["content"].(string), m.Body.(map[string]interface{})["toUserName"].(string))
		return
	})

	trigger.When("tuling", func(t *Trigger, b *bot.WeixinBot, m *domain.Message) (err error) {
		b.Tuling(m.Body.(string))
		return
	})
}
