package mns

import (
	"encoding/json"
	"github.com/cocotyty/summer"
	"github.com/gogap/ali_mns"
	"qiniupkg.com/x/errors.v7"
	"qiniupkg.com/x/log.v7"
	"time"
	"wxBot/bot"
	"wxBot/provider"
)

type Handler func(*Trigger, *bot.WeixinBot, *Message)
type Trigger struct {
	CacherFactory *provider.CacherFactory `sm:"*"`
	MQ            *MQ                     `sm:"*"`
	handlers      map[string][]Handler
}

func (this *Trigger) Init() {
	this.handlers = map[string][]Handler{}
}

func (this *Trigger) idle() {
	time.Sleep(3 * time.Second)
}

var (
	ErrNoHandler error = errors.New("找不到Handler")
)

func (this *Trigger) When(t string, h Handler) {
	_, has := this.handlers[t]
	if !has {
		this.handlers[t] = []Handler{}
	}
	this.handlers[t] = append(this.handlers[t], h)
}

func (this *Trigger) Send(id, t string, body interface{}) (err error) {
	msg := &Message{
		BotID: id,
		Type:  t,
		Body:  body,
	}
	bytes, err := json.Marshal(msg)
	if err != nil {
		return
	}
	this.MQ.Send(bytes)
	return

}

func (this *Trigger) handle(resp ali_mns.MessageReceiveResponse) (err error) {
	theMsg := &Message{}
	err = json.Unmarshal(resp.MessageBody, theMsg)
	if err != nil {
		return
	}
	log.Println(theMsg)
	handlers, has := this.handlers[theMsg.Type]
	if !has {
		err = ErrNoHandler
		return
	}

	theBot := bot.NewBot(theMsg.BotID, this.CacherFactory.NewCacher(theMsg.BotID))

	for _, h := range handlers {
		h(this, theBot, theMsg)
	}
	return
}
func (this *Trigger) Ready() {
	respChan := make(chan ali_mns.MessageReceiveResponse)
	errChan := make(chan error)
	this.MQ.Recv(respChan, errChan)
	this.Router()
	go func() {
		for {
			select {
			case resp := <-respChan:

				err := this.handle(resp)
				if err != nil {
					log.Println(err)
				}

				err = this.MQ.Queue.DeleteMessage(resp.ReceiptHandle)
				if err != nil {
					log.Println(err)
				}

			case err := <-errChan:
				log.Println(err)
				this.idle()
				continue
				//if ali_mns.ERR_MNS_QUEUE_NOT_EXIST.IsEqual(err) {
				//}
				//if ali_mns.ERR_MNS_MESSAGE_NOT_EXIST.IsEqual(err) {
				//}
			}
		}
	}()
}

func init() {
	summer.Add("Trigger", &Trigger{})
}
