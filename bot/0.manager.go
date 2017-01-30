package bot

import (
	"github.com/cocotyty/summer"
	"github.com/lizhengqiang/wxBot/domain"
	"github.com/lizhengqiang/wxBot/provider"
	"sync"
)

func init() {
	summer.Put(&BotManager{})
}

type BotManager struct {
	CacherFactory *provider.CacherFactory `sm:"*"`
	MQ            domain.MessageQueue     `sm:"*"`
	bots          map[string]*WeixinBot
	mutex         sync.Mutex
}

func (m *BotManager) Init() {
	m.bots = map[string]*WeixinBot{}
}

func (m *BotManager) Get(sessionID string) (self *WeixinBot) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	self, has := m.bots[sessionID]
	if !has {
		self = NewBot(sessionID, m.CacherFactory.NewCacher(sessionID), m.MQ)
		m.bots[sessionID] = self
	}

	self.ListenMode()
	return
}
