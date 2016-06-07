package bot

import (
	"fmt"
	"github.com/cocotyty/summer"
	"strconv"
	"strings"
	"time"
)

type AddMsg struct {
	MsgType      int64
	FromUserName string
	ToUserName   string
	Content      string
}

type Msg struct {
	Type         int64
	Content      string
	FromUserName string
	ToUserName   string
	LocalID      string
	ClientMsgId  string
}

func (this *WeixinBot) handleMsg(msgList []*AddMsg) {
	for _, msg := range msgList {
		// 机器人保留
		if strings.Contains(msg.Content, "#Bot#") {
			continue
		}

		// 文件助手
		if msg.ToUserName == "filehelper" {
			this.fileHelper(msg)
			continue
		}

		// 自己的
		if msg.FromUserName == this.getMe().UserName {

			continue
		}

		// 群消息
		if strings.Contains(msg.FromUserName, "@@") {
			this.groupMessage(msg)
			continue

		}

		// 用户消息
		this.contactMessage(msg)
	}
}




type SendMsgRequest struct {
	BaseRequest *BaseRequest
	Msg         *Msg
}

type SendMsgResponse struct {
	BaseResponse *BaseResponse
}

type SendMsgBody struct {
	Content    string `json:"content"`
	ToUserName string `json:"toUserName"`
}

func (bot *WeixinBot) DoSendMsg(content, toUserName string) {
	clientMsgId := strconv.FormatInt(time.Now().Unix()*1000+time.Now().Unix(), 10)
	request := SendMsgRequest{
		BaseRequest: bot.getBaseRequest(),
		Msg: &Msg{
			Type:         1,
			Content:      content,
			FromUserName: bot.getMe().UserName,
			ToUserName:   toUserName,
			LocalID:      clientMsgId,
			ClientMsgId:  clientMsgId,
		},
	}
	response := SendMsgResponse{}
	bot.PostJson(fmt.Sprintf("/webwxsendmsg?pass_ticket=%s", bot.getProperty(passTicket)), request, &response)
}

func (this *WeixinBot) SendMsg(content, toUserName string) {
	summer.GetStoneWithName("Trigger").(*Trigger).Send(this.ID, "sendMsg", &SendMsgBody{content, toUserName})
}