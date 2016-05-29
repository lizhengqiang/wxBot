package bot

import (
	"strings"
	"strconv"
	"time"
	"fmt"
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

func (this *WeixinBot) groupMessage(msg *AddMsg) {
	contents := strings.Split(msg.Content, `:<br/>`)
	var content string
	var userName string
	if len(contents) == 1 {
		content = contents[0]
	} else if len(contents) == 2 {
		userName = contents[0]
		content = contents[1]
	} else {
		content = msg.Content
	}

	this.Println(msg.FromUserName, userName, content)

}

func (this *WeixinBot) contactMessage(msg *AddMsg) {
	this.Println(msg.FromUserName, msg.Content)
}

type SendMsgRequest struct {
	BaseRequest *BaseRequest
	Msg         *Msg
}

type SendMsgResponse struct {
	BaseResponse *BaseResponse
}

func (bot *WeixinBot) SendMsg(content, toUserName string) {
	clientMsgId := strconv.FormatInt(time.Now().Unix() * 1000 + time.Now().Unix(), 10)
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
	// name := bot.GetRemarkName(toUserName)
}

