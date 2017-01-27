package api

import (
	"github.com/go-macaron/session"
	"github.com/mougeli/beauty"
	"gopkg.in/macaron.v1"
	"github.com/lizhengqiang/wxBot/bot"
)

func GetQrCodeUrl(ctx *macaron.Context, sess session.Store, b *bot.WeixinBot, r beauty.Render) {
	// 清理掉上次的
	r.OK(b.GetQrcodeUrl())
}
