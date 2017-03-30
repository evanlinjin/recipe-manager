package talkrelay

import (
	"bytes"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
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
}

func MakeWSManager(upgrader *websocket.Upgrader, w http.ResponseWriter, r *http.Request) (wsm *WSManager, e error) {
	wsm = &WSManager{}
	wsm.enc = MakeEncryptor()
	wsm.msgs = MakeMessageManager()
	wsm.conn, e = upgrader.Upgrade(w, r, nil)
	return
}

func (m *WSManager) DeprecatedWriteMessage(data []byte) error {
	// Make random Signature.
	sig, _ := m.enc.makeKey()

	// Make Data into Package with Signature.
	pkg, _ := MakePackage(Msg{string(data)}, sig)

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

func (m *WSManager) DeprecatedReadMessage() (data []byte, e error) {
	m.readMux.Lock()
	_, msg, e := m.conn.ReadMessage()
	m.readMux.Unlock()

	// Split message into signature and data.
	split := bytes.Split(msg, []byte("."))
	if len(split) != 2 {
		e = fmt.Errorf("expected %v parts, got %v from %v", 2, len(split), string(msg))
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
	var txt Msg
	e = ReadPackage(pkg, sig, &txt)
	if e != nil {
		return
	}
	data = []byte(txt.Msg)
	return
}

func (m *WSManager) SendRequestMessage(cmd string, data interface{}) error {
	msg := m.msgs.MakeRequestMessage(&cmd, data)
	return m.sendMessage(msg)
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
	_, encPkg, e := m.conn.ReadMessage()
	m.readMux.Unlock()

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
	e = ReadPackage(pkg, sig, msg)
	if e != nil {
		return
	}
	// Check msg.
	e = m.msgs.CheckIncomingMessage(msg)
	return
}
