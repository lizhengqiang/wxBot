package bot

import (
	"fmt"
	"time"
)
type SyncKey struct {
	Count int64
	List  []struct {
		Key int64
		Val int64
	}
}

type WebWeixinSyncRequest struct {
	BaseRequest *BaseRequest
	SyncKey     *SyncKey
	rr          int64
}

type WebWeixinSyncResponse struct {
	BaseResponse *BaseResponse
	SyncKey      *SyncKey
	AddMsgList   []*AddMsg
}

func (bot *WeixinBot) WebWeixinSync() *WebWeixinSyncResponse {
	theSyncKey := &SyncKey{}
	bot.unmarshal(syncKey, theSyncKey)
	request := WebWeixinSyncRequest{
		BaseRequest: bot.getBaseRequest(),
		SyncKey:     theSyncKey,
		rr:          time.Now().Unix(),
	}

	response := &WebWeixinSyncResponse{}

	u := fmt.Sprintf("/webwxsync?sid=%s&skey=%s&pass_ticket=%s", bot.get(wxsid), bot.get(skey), bot.get(passTicket))
	bot.PostJson(u, request, response)

	if response.BaseResponse != nil && response.BaseResponse.Ret == 0 {
		bot.saveSyncKey(response.SyncKey)
	}

	return response
}
