package talkrelay

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"time"
)

type Msg struct {
	Msg string `json:"msg"`
}

// WSManager manages a WebSocket connection.
type WSManager struct {
	conn *websocket.Conn
	msgs MessageManager
	enc  Encryptor

	readMux  sync.RWMutex
	writeMux sync.Mutex

	QuitChan chan bool
}

func MakeWSManager(upgrader *websocket.Upgrader, w http.ResponseWriter, r *http.Request) (wsm WSManager, e error) {
	wsm = WSManager{}
	wsm.enc = MakeEncryptor()
	wsm.msgs = MakeMessageManager()
	wsm.conn, e = upgrader.Upgrade(w, r, nil)
	if e != nil {
		return
	}
	wsm.QuitChan = make(chan bool)
	return
}

func (m *WSManager) Handshake(wait time.Duration) (e error) {
	key, _ := m.enc.makeKey()

	var resChan = make(chan *Message)

	go func() {
		res, _ := m.GetMessage()
		select {
		case resChan <- res:
		default:
		}
	}()

	req, _ := m.SendRequestMessage("handshake", key)
	timer := time.NewTimer(wait)
	for {
		select {
		case <-timer.C:
			e = fmt.Errorf("handshake failed, timeout")
			fmt.Println("timer dome.")
			timer.Stop()
			goto DoneHandShake

		case res := <-resChan:
			timer.Stop()
			switch {
			case res == nil, res.Meta == nil, res.ReqMeta == nil:
				break
			case res.Type != _TypeResponse, res.Command != "handshake":
				break
			case res.ReqMeta.Timestamp != req.Meta.Timestamp:
				break
			case res.ReqMeta.ID != req.Meta.ID:
				break
			default:
				m.enc.setKey(key)
				return nil
			}
			e = fmt.Errorf("handshake failed, invalid response")
			goto DoneHandShake
		}
	}
DoneHandShake:
	m.Close(-1, e.Error())
	return e
}

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

func (m *WSManager) SendRequestMessage(cmd string, data interface{}) (*Message, error) {
	msg := m.msgs.MakeRequestMessage(&cmd, data)
	return msg, m.sendMessage(msg)
}

func (m *WSManager) SendResposneMessage(reqMsg *Message, data interface{}) error {
	msg, e := m.msgs.MakeResponseMessage(reqMsg, data)
	if e != nil {
		return e
	}
	return m.sendMessage(msg)
}

func (m *WSManager) sendMessage(msg *Message) error {
	if msg == nil {
		return nil
	}
	// Make random Signature.
	sig, _ := m.enc.makeKey()
	// Put Message into Package with Signature.
	pkg, _ := MakePackage(msg, sig)
	// Encrypt Data and Signature.
	encSig, _ := m.enc.Encrypt(sig)
	encPkg, _ := m.enc.Encrypt(pkg)
	// Join with dot.
	out := append(encSig, byte('.'))
	out = append(out, encPkg...)

	m.writeMux.Lock()
	defer m.writeMux.Unlock()
	return m.conn.WriteMessage(websocket.TextMessage, out)
}

func (m *WSManager) GetMessage() (msg *Message, e error) {
	m.readMux.Lock()
	typ, encPkg, e := m.conn.ReadMessage()
	m.readMux.Unlock()

	// Check type of ws msg.
	switch typ {
	case websocket.CloseMessage, -1:
		return
	case websocket.TextMessage:
		break
	default:
		e = fmt.Errorf("unexpected message type", typ)
		return
	}

	// Split message into signature and data.
	split := bytes.Split(encPkg, []byte("."))
	if len(split) != 2 {
		e = fmt.Errorf("expected %v parts, got %v from %v", 2, len(split), string(encPkg))
		return
	}
	// Decrypt signature and package.
	sig, e := m.enc.Decrypt(split[0])
	if e != nil {
		e = fmt.Errorf("while decrypting signature, %v", e)
		return
	}
	pkg, e := m.enc.Decrypt(split[1])
	if e != nil {
		e = fmt.Errorf("while decrypting package, %v", e)
		return
	}
	// Vertify data with signature.
	msg = &Message{}
	e = ReadPackage(pkg, sig, msg)
	if e != nil {
		return
	}
	test, _ := json.Marshal(msg)
	fmt.Println(string(test))
	// Check msg.
	e = m.msgs.CheckIncomingMessage(msg)
	return
}
