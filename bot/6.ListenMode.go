package bot

import (
	"github.com/qiniu/log"
	"qiniupkg.com/x/errors.v7"
)

func init() {

}
func (bot *WeixinBot) Stop() {
	bot.Set(IsRunning, FALSE)
}

var (
	ErrStopHandleMsg error = errors.New("停止处理消息")
)

func (bot *WeixinBot) ListenMode() {
	// 防止打开多个监听
	if !bot.IsLoopRunning {
		bot.IsLoopRunning = true
		go func() {
			for {
				// 非正常运行时关闭
				if !bot.IsRunning() || !bot.IsSigned() {
					return
				}

				err := bot.HandleMsg()
				if err != nil {
					bot.IsLoopRunning = false
					return
				}
			}
		}()
	}
}

func (bot *WeixinBot) syncCheck() (retcode, selector int64, err error) {
	hosts := map[string]interface{}{
		"webpush.weixin.qq.com": nil,
		"webpush.wechat.com":    nil,
		"webpush1.wechat.com":   nil,
		"webpush2.wechat.com":   nil,
		"webpush.wx.qq.com":     nil,
		"webpush2.wx.qq.com":    nil,
		"webpush.wx2.qq.com":    nil,
	}

	if _, has := hosts[bot.getProperty(syncKeyHost)]; has {
		hosts = map[string]interface{}{bot.getProperty(syncKeyHost): nil}
	}

	for host := range hosts {
		retcode, selector, err = bot.SyncCheck(host)
		log.Println("handleMsg", host, retcode, selector)
		if err != nil {
			continue
		}
		if retcode >= 1100 {
			continue
		}
		bot.setProperty(syncKeyHost, host)
		break
	}

	return
}

func (bot *WeixinBot) HandleNewMsg() (err error) {

	msgList := bot.WebWeixinSync()
	if msgList.AddMsgList == nil || len(msgList.AddMsgList) == 0 {
		return
	}
	bot.handleMsg(msgList.AddMsgList)
	return
}

func (bot *WeixinBot) HandleMsg() (err error) {

	log.Println("require", "handleMsg")
	bot.rw.Lock()
	defer func() {
		log.Println("release", "handleMsg")
		bot.rw.Unlock()
	}()
	log.Println("lock", "handleMsg")
	retcode, selector, err := bot.syncCheck()
	if err != nil || retcode >= 1100 {
		bot.Set(IsRunning, FALSE)
		bot.log("# 你退出了,债见~_~")
		return ErrStopHandleMsg
	}

	switch selector {
	case 2:
		bot.HandleNewMsg()
	case 7:
		bot.log("# 发现你玩手机了!")

	default:
	}
	return
}
