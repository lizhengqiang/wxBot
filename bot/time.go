package bot

import "time"

const IdleDuration = 1 * time.Second

func (bot *WeixinBot) Idle() {
	time.Sleep(IdleDuration)
}
