package bot

import (
	"io/ioutil"
	"net/url"
	"qiniupkg.com/x/errors.v7"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrSyncCheck error = errors.New("SyncCheck失败")
)


type SyncCheckResponseBody struct {
	retcode  int64
	selector int64
}

// 保存微信返回的SyncKey
func (this *WeixinBot) saveSyncKey(theSyncKey *SyncKey) {
	this.marshal(syncKey, theSyncKey)
	syncKeyList := make([]string, theSyncKey.Count)
	for i, v := range theSyncKey.List {
		syncKeyList[i] = strconv.FormatInt(v.Key, 10) + "_" + strconv.FormatInt(v.Val, 10)
	}
	this.setProperty(syncKeyString, strings.Join(syncKeyList, "|"))
}

func (bot *WeixinBot) SyncCheck() (retcode, selector int64, err error) {
	queryValues := url.Values{}
	queryValues.Add("synckey", bot.getProperty(syncKeyString))
	queryValues.Add("skey", bot.getProperty(skey))
	queryValues.Add("uin", bot.get(wxuni))
	queryValues.Add("r", bot.timestamp())
	queryValues.Add("deviceid", bot.getProperty(deviceId))
	queryValues.Add("sid", bot.getProperty(wxsid))
	queryValues.Add("_", bot.timestamp())
	u := "https://webpush.weixin.qq.com/cgi-bin/mmwebwx-bin/synccheck?" + queryValues.Encode()
	resp, err := bot.httpClient.Get(u)
	if err != nil {
		bot.log(err.Error())
		err = ErrSyncCheck
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		bot.log(err.Error())
		err = ErrSyncCheck
		return
	}

	re, _ := regexp.Compile(`window.synccheck={retcode:"(\d+)",selector:"(\d+)"}`)
	sub := re.FindSubmatch(body)
	if len(sub) >= 2 {
		retcode, _ := strconv.ParseInt(string(sub[1]), 10, 64)
		selector, _ := strconv.ParseInt(string(sub[2]), 10, 64)
		return retcode, selector, nil
	}

	err = ErrSyncCheck
	return
}
