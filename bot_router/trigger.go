package bot_router

import (
	"github.com/cocotyty/summer"
	"github.com/lizhengqiang/wxBot/bot"
	"github.com/lizhengqiang/wxBot/domain"
	"github.com/lizhengqiang/wxBot/provider"
	"qiniupkg.com/x/errors.v7"
	"time"
)

type Handler func(*Trigger, *bot.WeixinBot, *domain.Message) error
type Trigger struct {
	CacherFactory *provider.CacherFactory `sm:"*"`
	MQ            domain.MessageQueue     `sm:"*"`
	handlers      map[string][]Handler
	BotManager    *bot.BotManager `sm:"*"`
}

func (trigger *Trigger) Init() {
	trigger.handlers = map[string][]Handler{}

}

func (trigger *Trigger) idle() {
	time.Sleep(3 * time.Second)
}

var (
	ErrNoHandler error = errors.New("找不到Handler")
)

func (trigger *Trigger) When(t string, h Handler) {
	_, has := trigger.handlers[t]
	if !has {
		trigger.handlers[t] = []Handler{}
	}
	trigger.handlers[t] = append(trigger.handlers[t], h)
}

func (trigger *Trigger) Send(id, t string, body interface{}) (err error) {
	return trigger.MQ.Send(&domain.Message{
		BotID: id,
		Type:  t,
		Body:  body,
	})

}

func (trigger *Trigger) handle(msg *domain.Message) (err error) {
	handlers, has := trigger.handlers[msg.Type]
	if !has {
		err = ErrNoHandler
		return
	}
	self := trigger.BotManager.Get(msg.BotID)
	for _, h := range handlers {
		err = h(trigger, self, msg)
		if err != nil {
			return
		}
	}
	return
}

func (trigger *Trigger) ListenMode() {

}
func (trigger *Trigger) Ready() {
	trigger.MQ.RegisterHandler(trigger.handle)
	go trigger.ListenMode()
}

func init() {
	summer.Add("Trigger", &Trigger{})
}
