package bot

import (
	"encoding/json"
	"github.com/cocotyty/summer"
	"github.com/gogap/ali_mns"
	"qiniupkg.com/x/errors.v7"
	"time"
	"github.com/lizhengqiang/wxBot/provider"
)

type Handler func(*Trigger, *WeixinBot, *provider.Message)
type Trigger struct {
	CacherFactory *provider.CacherFactory `sm:"*"`
	MQ            *provider.MQ            `sm:"*"`
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
	msg := &provider.Message{
		BotID: id,
		Type:  t,
		Body:  body,
	}
	bytes, err := json.Marshal(msg)
	if err != nil {
		return
	}
	err = this.MQ.Send(bytes)
	return

}

func (this *Trigger) handle(resp ali_mns.MessageReceiveResponse) (err error) {
	theMsg := &provider.Message{}
	err = json.Unmarshal(resp.MessageBody, theMsg)
	if err != nil {
		return
	}

	handlers, has := this.handlers[theMsg.Type]
	if !has {
		err = ErrNoHandler
		return
	}

	theBot := NewBot(theMsg.BotID, this.CacherFactory.NewCacher(theMsg.BotID))

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
				go this.handle(resp)
				this.MQ.Queue.DeleteMessage(resp.ReceiptHandle)

			case _ = <-errChan:
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
