package bot

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type InitWebWeixinRequestBody struct {
	BaseRequest *BaseRequest
}

type InitWebWeixinResponseBody struct {
	BaseResponse *BaseResponse
	SyncKey      *SyncKey
	User         *User
}

func (this *WeixinBot) InitWebWeixin() int64 {
	requestBody := InitWebWeixinRequestBody{
		BaseRequest: this.getBaseRequest(),
	}
	u := fmt.Sprintf("/webwxinit?pass_ticket=%s&skey=%s&r=%s", this.getProperty(passTicket), this.getProperty(skey), this.timestamp())

	respJson := &InitWebWeixinResponseBody{}

	this.PostJson(u, requestBody, respJson)

	this.marshal(me, respJson.User)

	this.saveSyncKey(respJson.SyncKey)

	return respJson.BaseResponse.Ret

}

type WebWeixinStatusNotifyRequest struct {
	BaseRequest  *BaseRequest
	Code         int64
	FromUserName string
	ToUserName   string
	ClientMsgId  int64
}

type WebWeixinStatusNotifyResponseBody struct {
	BaseResponse *BaseResponse
}

var ErrStatusNotify error = errors.New("检查状态出错")

func (bot *WeixinBot) WebWeixinStatusNotify() (err error) {
	my := &User{}
	bot.unmarshal(me, my)
	requestBody := WebWeixinStatusNotifyRequest{
		BaseRequest:  bot.getBaseRequest(),
		Code:         int64(3),
		FromUserName: my.UserName,
		ToUserName:   my.UserName,
		ClientMsgId:  time.Now().Unix(),
	}

	respBody, _ := bot.SimplePostJson(fmt.Sprintf("/webwxstatusnotify?lang=zh_CN&pass_ticket=%s", bot.getProperty(passTicket)), requestBody)
	respJson := WebWeixinStatusNotifyResponseBody{}

	err = json.Unmarshal(respBody, &respJson)

	if err != nil {
		bot.log(err.Error())
		return
	}

	if respJson.BaseResponse.Ret == int64(0) {
		return nil
	}
	return ErrStatusNotify
}
