package api

import "gopkg.in/macaron.v1"

func RegisterRoutes(m *macaron.Macaron) {
	m.Group("/api", func() {
		m.Get("/start", Start)
		m.Get("/qrcode", GetQrCodeUrl)
	}, BotMiddleware)
}
