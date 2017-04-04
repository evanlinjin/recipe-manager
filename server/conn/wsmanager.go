package conn

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

// WSManager manages a WebSocket connection. ReadMux and writeMux are used to
// ensure that writes and reads to the WebSocket connection on different
// goroutines do not collide. QuitChan tells all the associated goroutines to
// end.
type WSManager struct {
	conn *websocket.Conn
	msgs MessageManager

	readMux  sync.RWMutex
	writeMux sync.Mutex

	QuitChan chan bool
}

// MakeWSManager makes a new WSManager instance. It returns an error if it is
// unable to upgrade the connection to WebSocket.
func MakeWSManager(upgrader *websocket.Upgrader, w http.ResponseWriter, r *http.Request) (wsm WSManager, e error) {
	wsm = WSManager{}
	wsm.msgs = MakeMessageManager()
	wsm.conn, e = upgrader.Upgrade(w, r, nil)
	if e != nil {
		return
	}
	wsm.QuitChan = make(chan bool)
	return
}

// Close closes the connected, and tell all associated listening/reading
// goroutines to end, via the QuitChan channel.
func (m *WSManager) Close(code int, msg string) {
	for {
		select {
		case m.QuitChan <- false:
		default:
			goto DoneQuitChan
		}
	}

DoneQuitChan:
	if m.conn != nil {
		type Msg struct {
			Code int    `json:"code"`
			Msg  string `json:"msg"`
		}
		m.SendRequestMessage("close", Msg{code, msg})
		m.conn.Close()
	}
}

// SendMessage sends a message via WebSocket.
func (m *WSManager) SendMessage(msg *Message) error {
	if msg == nil {
		return nil
	}
	// Put Message into Package with Signature.
	pkg, e := json.Marshal(msg)
	if e != nil {
		return e
	}

	m.writeMux.Lock()
	defer m.writeMux.Unlock()
	return m.conn.WriteMessage(websocket.BinaryMessage, pkg)
}

// SendRequestMessage sends, specifically, a request message. This is a
// convenience function.
func (m *WSManager) SendRequestMessage(cmd string, data interface{}) (*Message, error) {
	msg := m.msgs.MakeRequestMessage(&cmd, data)
	return msg, m.SendMessage(msg)
}

// SendResponseMessage sends, specifically, a response message. This is a
// convenience function.
func (m *WSManager) SendResponseMessage(reqMsg *Message, data interface{}) error {
	msg, e := m.msgs.MakeResponseMessage(reqMsg, data)
	if e != nil {
		return e
	}
	return m.SendMessage(msg)
}

func (m *WSManager) GetMessage() (msg *Message, e error) {
	m.readMux.Lock()
	typ, pkg, e := m.conn.ReadMessage()
	m.readMux.Unlock()
	if e != nil {
		return
	}

	// Check type of ws msg.
	switch typ {
	case websocket.BinaryMessage:
		break
	case websocket.CloseMessage, -1:
		e = &ErrCloseMessage{typ}
		return
	default:
		e = &ErrUnexpectedMessage{typ}
		return
	}

	// Vertify data with signature.
	msg = &Message{}
	e = json.Unmarshal(pkg, msg)
	if e != nil {
		return
	}
	// Check msg.
	e = m.msgs.CheckIncomingMessage(msg)
	return
}
