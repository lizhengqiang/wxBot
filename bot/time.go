package bot

import "time"

const IdleDuration = 1 * time.Second

func (this *WeixinBot) Idle() {
	time.Sleep(IdleDuration)
}
