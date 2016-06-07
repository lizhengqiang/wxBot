package provider

import (
	"github.com/cocotyty/summer"
	"github.com/gogap/ali_mns"
	"os"
	"qiniupkg.com/x/log.v7"
)

type MQ struct {
	AccessKey string
	SecretKey string
	MnsUrl    string
	QueueName string
	Queue     ali_mns.AliMNSQueue
	Client    ali_mns.MNSClient
}

func (this *MQ) Init() {
	this.AccessKey = os.Getenv("aliyun.accessKey")
	this.SecretKey = os.Getenv("aliyun.secretKey")
	this.MnsUrl = os.Getenv("mq.url")
	this.QueueName = os.Getenv("mq.queue")
}

func (this *MQ) Ready() {
	log.Println(this.MnsUrl, this.AccessKey, this.SecretKey)
	this.Client = ali_mns.NewAliMNSClient(this.MnsUrl, this.AccessKey, this.SecretKey)
	this.Queue = ali_mns.NewMNSQueue(this.QueueName, this.Client)
}
func (this *MQ) Send(body []byte) (err error) {
	msg := ali_mns.MessageSendRequest{
		MessageBody:  body,
		DelaySeconds: 0,
		Priority:     8,
	}
	_, err = this.Queue.SendMessage(msg)
	return
}

func (this *MQ) Recv(respChan chan ali_mns.MessageReceiveResponse, errChan chan error) (err error) {
	go this.Queue.ReceiveMessage(respChan, errChan)
	return
}

func init() {
	summer.Add("MQ", &MQ{})
}
