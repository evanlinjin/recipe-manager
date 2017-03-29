package talkrelay

import (
	"github.com/gorilla/websocket"
	"sync"
	"net/http"
	"bytes"
	"fmt"
)

type Msg struct {
	Msg string `json:"msg"`
}

// WSManager manages a WebSocket connection.
type WSManager struct {
	c *websocket.Conn
	enc Encryptor

	readMux sync.RWMutex
	writeMux sync.Mutex
}

func MakeWSManager(upgrader *websocket.Upgrader, w http.ResponseWriter, r *http.Request) (wsm *WSManager, e error) {
	wsm = &WSManager{}
	wsm.enc = MakeEncryptor()
	wsm.c, e = upgrader.Upgrade(w, r, nil)
	return
}

func (m *WSManager) WriteMessage(data []byte) error {
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
	return m.c.WriteMessage(websocket.TextMessage, out)
}

func (m *WSManager) ReadMessage() (data []byte, e error) {
	m.readMux.Lock()
	_, msg, e := m.c.ReadMessage()
	m.readMux.Unlock()

	// Split message into signature and data.
	split := bytes.Split(msg, []byte("."))
	if len(split) != 2 {
		e = fmt.Errorf("expected %v parts, got %v from %v", 2, len(split), string(msg))
		return
	}
	fmt.Println("ENC SIGNATURE:", string(split[0]), len(split[0]))
	fmt.Println("ENC PACKAGE:", string(split[1]), len(split[1]))

	// Decrypt signature and package.
	sig, e := m.enc.Decrypt(split[0])
	if e != nil {
		e = fmt.Errorf("while decrypting signature, %v", e)
		return
	}
	fmt.Println("RAW SIGNATURE:", string(sig), len(sig))
	pkg, e := m.enc.Decrypt(split[1])
	if e != nil {
		e = fmt.Errorf("while decrypting package, %v", e)
		return
	}
	fmt.Println("RAW PACKAGE:", string(pkg), len(pkg))

	// Vertify data with signature.
	var txt Msg
	e = ReadPackage(pkg, sig, &txt)
	if e != nil {
		return
	}
	data = []byte(txt.Msg)
	return
}