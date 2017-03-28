package talkrelay

import (
	"github.com/gorilla/websocket"
	"sync"
	"net/http"
)

// WSManager manages a WebSocket connection.
type WSManager struct {
	c *websocket.Conn

	readMux sync.RWMutex
	writeMux sync.Mutex
}

func MakeWSManager(upgrader *websocket.Upgrader, w http.ResponseWriter, r *http.Request) (wsm *WSManager,e error) {
	wsm = &WSManager{}
	wsm.c, e = upgrader.Upgrade(w, r, nil)
	return
}

func (m *WSManager) WriteMessage(data []byte) error {
	m.writeMux.Lock()
	defer m.writeMux.Unlock()
	return m.c.WriteMessage(websocket.TextMessage, data)
}

func (m *WSManager) ReadMessage() (data []byte, e error) {
	m.readMux.Lock()
	defer m.readMux.Unlock()
	_, data, e = m.c.ReadMessage()
	return
}