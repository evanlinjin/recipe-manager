package conn

import (
	"fmt"
	"sync"
	"time"
)

const (
	_TypeRequest  = 0
	_TypeResponse = 1
)

type MessageMeta struct {
	ID        int64 `json:"id"`
	Timestamp int64 `json:"ts"`
}

type Message struct {
	Command string       `json:"cmd"`
	Type    int64        `json:"typ"`
	Data    interface{}  `json:"data,omitempty"`
	Meta    *MessageMeta `json:"meta"`
	ReqMeta *MessageMeta `json:"req,omitempty"`
}

type MessageManager struct {
	sync.Mutex
	outgoingID int64
	incomingID int64
}

func MakeMessageManager() MessageManager {
	return MessageManager{
		outgoingID: 0,
		incomingID: 0,
	}
}

func (m *MessageManager) MakeRequestMessage(cmd *string, data interface{}) (msg *Message) {
	m.Lock()
	m.outgoingID -= 1
	m.Unlock()
	msg = &Message{
		Command: *cmd,
		Type:    _TypeRequest,
		Data:    data,
		Meta: &MessageMeta{
			ID:        m.outgoingID,
			Timestamp: time.Now().Unix(),
		},
	}
	return
}

func (m *MessageManager) MakeResponseMessage(reqMsg *Message, data interface{}) (msg *Message, e error) {
	if reqMsg == nil {
		e = fmt.Errorf("reqMsg is nil")
		return
	} else if reqMsg.Meta == nil {
		e = fmt.Errorf("reqMsg.Meta is nil")
		return
	} else if reqMsg.Meta.ID < 1 {
		e = fmt.Errorf("reqMsg.Meta.ID is < 1, at %v", reqMsg.Meta.ID)
		return
	}
	m.Lock()
	m.outgoingID -= 1
	m.Unlock()
	msg = &Message{
		Command: reqMsg.Command,
		Type:    _TypeResponse,
		Data:    data,
		Meta: &MessageMeta{
			ID:        m.outgoingID,
			Timestamp: time.Now().Unix(),
		},
		ReqMeta: reqMsg.Meta,
	}
	return
}

func (m *MessageManager) CheckIncomingMessage(msg *Message) error {
	switch {
	case msg == nil:
		return fmt.Errorf("nil msg")
	case msg.Meta == nil:
		return fmt.Errorf("msg has nil meta")
	case msg.Meta.ID < 0:
		return fmt.Errorf("incoming msg has -ve id; %v", msg.Meta.ID)
	case msg.Type == _TypeResponse:
		switch {
		case msg.ReqMeta == nil:
			return fmt.Errorf("response msg has nil reqMeta")
		case msg.ReqMeta.ID > 0:
			return fmt.Errorf("reqMeta of incoming msg is invalid; %v", msg.ReqMeta.ID)
		}
	}
	if msg.Meta.ID != m.incomingID+1 {
		fmt.Println(
			"[MessageManager.CheckIncomingMessage] Unexpected ID.",
			"GOT:", msg.Meta.ID,
			"EXPECTED:", m.incomingID+1,
		)
	}
	m.Lock()
	m.incomingID = msg.Meta.ID
	m.Unlock()
	return nil
}
