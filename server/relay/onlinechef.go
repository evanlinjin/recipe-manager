package relay

import (
	"github.com/evanlinjin/recipe-manager/server/chefs"
	"github.com/evanlinjin/recipe-manager/server/conn"
	"sync"
)

type OnlineChef struct {
	sync.Mutex
	channels map[string]chan *conn.Message
	Info     *chefs.SessionInfo
}

func NewOnlineChef(info *chefs.SessionInfo) *OnlineChef {
	return &OnlineChef{
		channels: make(map[string]chan *conn.Message),
		Info:     info,
	}
}

// AddChannel adds a channel and returns channel count (aka number of sessions).
func (c *OnlineChef) AddChannel(sessionID string, channel chan *conn.Message) int {
	c.Lock()
	defer c.Unlock()
	c.channels[sessionID] = channel
	return len(c.channels)
}

// RemoveChannel removes a channel and return channel count (aka number of sessions).
func (c *OnlineChef) RemoveChannel(sessionID string) int {
	c.Lock()
	defer c.Unlock()
	delete(c.channels, sessionID)
	return len(c.channels)
}

// SendMessage sends a message to the chef (distributed to all it's sessions).
func (c *OnlineChef) SendMessage(msg *conn.Message) {
	c.Lock()
	defer c.Unlock()

	for _, channel := range c.channels {
		go func() { channel <- msg }()
	}
}
