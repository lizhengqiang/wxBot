package provider

import (
	"encoding/json"
	"github.com/cocotyty/summer"
	"github.com/gogap/ali_mns"
	"github.com/lizhengqiang/wxBot/domain"
	"os"
	"time"
)

type AliMQ struct {
	AccessKey string
	SecretKey string
	MnsUrl    string
	QueueName string
	Queue     ali_mns.AliMNSQueue
	Client    ali_mns.MNSClient
	Log       *summer.SimpleLog

	handlers []domain.MessageHandler
}

func (p *AliMQ) Init() {
	p.Log = summer.NewSimpleLog("AliMessageQueue", summer.DebugLevel)
	p.AccessKey = os.Getenv("aliyun.accessKey")
	p.SecretKey = os.Getenv("aliyun.secretKey")
	p.MnsUrl = os.Getenv("mq.url")
	p.QueueName = os.Getenv("mq.queue")

	p.handlers = []domain.MessageHandler{}
}
func (p *AliMQ) MNSSend(body []byte) (err error) {
	msg := ali_mns.MessageSendRequest{
		MessageBody:  body,
		DelaySeconds: 0,
		Priority:     8,
	}
	_, err = p.Queue.SendMessage(msg)
	return
}

func (p *AliMQ) handle(resp ali_mns.MessageReceiveResponse) (err error) {
	msg := &domain.Message{}
	err = json.Unmarshal(resp.MessageBody, msg)
	if err != nil {
		return
	}
	p.Log.Info("handle", msg, p.handlers)

	for _, h := range p.handlers {
		err = h(msg)
		if err != nil {
			return
		}
	}
	return
}

func (p *AliMQ) idle() {
	time.Sleep(1 * time.Second)
}

func (p *AliMQ) handleAndDel(resp ali_mns.MessageReceiveResponse) {
	err := p.handle(resp)
	if err != nil {
		return
	}
	p.Queue.DeleteMessage(resp.ReceiptHandle)
}

func (p *AliMQ) Ready() {
	p.Client = ali_mns.NewAliMNSClient(p.MnsUrl, p.AccessKey, p.SecretKey)
	p.Queue = ali_mns.NewMNSQueue(p.QueueName, p.Client)
	respChan := make(chan ali_mns.MessageReceiveResponse)
	errChan := make(chan error)
	go p.Queue.ReceiveMessage(respChan, errChan)
	go func() {
		for {
			select {
			case resp := <-respChan:
				go p.handleAndDel(resp)

			case _ = <-errChan:
				p.idle()
				continue
				//if ali_mns.ERR_MNS_QUEUE_NOT_EXIST.IsEqual(err) {
				//}
				//if ali_mns.ERR_MNS_MESSAGE_NOT_EXIST.IsEqual(err) {
				//}
			}
		}
	}()
	return
}

func (p *AliMQ) Send(msg *domain.Message) (err error) {
	bytes, err := json.Marshal(msg)
	if err != nil {
		return
	}
	err = p.MNSSend(bytes)
	p.Log.Info("Send", string(bytes), err)
	return
}

func (p *AliMQ) RegisterHandler(handler domain.MessageHandler) (err error) {
	p.handlers = append(p.handlers, handler)
	return nil
}

func init() {
	summer.Add("AliMQ", &AliMQ{})
}
