package bot

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
	"github.com/qiniu/log"
)

type InitWebWeixinRequestBody struct {
	BaseRequest *BaseRequest
}

type InitWebWeixinResponseBody struct {
	BaseResponse *BaseResponse
	SyncKey      *SyncKey
	User         *User
}

func (bot *WeixinBot) InitWebWeixin() int64 {
	requestBody := InitWebWeixinRequestBody{
		BaseRequest: bot.getBaseRequest(),
	}
	u := fmt.Sprintf("/webwxinit?pass_ticket=%s&skey=%s&r=%s", bot.getProperty(passTicket), bot.getProperty(skey), bot.timestamp())

	respJson := &InitWebWeixinResponseBody{}

	bot.PostJson(u, requestBody, respJson)

	bot.marshal(me, respJson.User)

	bot.saveSyncKey(respJson.SyncKey)

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

func (bot *WeixinBot) WebWeixinStatusNotify(from, to string, code int64) (err error) {
	my := &User{}
	bot.unmarshal(me, my)
	requestBody := WebWeixinStatusNotifyRequest{
		BaseRequest:  bot.getBaseRequest(),
		Code:         code,
		FromUserName: from,
		ToUserName:   to,
		ClientMsgId:  time.Now().Unix(),
	}

	log.Println(bot.getProperty(baseUri))

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
