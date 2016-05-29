// 主要使用文件传输助手来控制
package bot

import (
	"fmt"
	"github.com/cocotyty/json"
	"strconv"
	"strings"
	"time"
)

func (this *WeixinBot) fileHelperResponse(content string) {
	clientMsgId := strconv.FormatInt(time.Now().Unix()*1000+time.Now().Unix(), 10)
	request := SendMsgRequest{
		BaseRequest: this.getBaseRequest(),
		Msg: &Msg{
			Type:         1,
			Content:      "#Bot#" + content,
			FromUserName: this.getMe().UserName,
			ToUserName:   "filehelper",
			LocalID:      clientMsgId,
			ClientMsgId:  clientMsgId,
		},
	}
	response := SendMsgResponse{}
	this.PostJson(fmt.Sprintf("/webwxsendmsg?pass_ticket=%s", this.getProperty(passTicket)), request, &response)
	this.Println("fileHelper", content)
}

func (this *WeixinBot) fileHelper(msg *AddMsg) {
	if msg.Content == "me" {
		bytes, _ := json.MarshalIndent(this.getMe(), "", "  ")
		this.fileHelperResponse(string(bytes))
		return
	}

	if msg.Content == "logs" {
		this.fileHelperResponse(this.Cacher.Get("logs"))
		return
	}

	if msg.Content == "now" {
		this.fileHelperResponse(this.get("status"))
		return
	}

	if msg.Content == "mps" {
		this.fileHelperResponse(this.get(mpList))
		return
	}

	if strings.Contains(msg.Content, "mps") {
		go this.mps(strings.Replace(msg.Content, "mps ", "", -1))
		return
	}

}
