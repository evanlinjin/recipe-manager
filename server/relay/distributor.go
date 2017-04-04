package relay

import (
	"errors"
	"github.com/evanlinjin/recipe-manager/server/chefs"
	"github.com/evanlinjin/recipe-manager/server/conn"
	"sync"
)

var (
	ErrOnlineChefDoesNotExist = errors.New(
		"online chef already does not exist - cannot remove")
)

// Distributor contains all online chefs and associated online sessions. It can
// distribute messages to these chefs via their channels.
type Distributor struct {
	sync.Mutex
	chefs map[string]*OnlineChef
}

// NewDistributor creates a new instance of Distributor.
func NewDistributor() *Distributor {
	return &Distributor{
		chefs: make(map[string]*OnlineChef),
	}
}

// AddSession adds a session to the Distributor.
func (d *Distributor) AddSession(info *chefs.SessionInfo, channel chan *conn.Message) {
	d.Lock()
	defer d.Unlock()

	// If we have the session already, just add one one.
	// Otherwise we create an OnlineChef and add channel to it.
	if onlineChef, has := d.chefs[info.ChefID]; has {
		onlineChef.AddChannel(info.SessionID, channel)
	} else {
		onlineChef = NewOnlineChef(info)
		onlineChef.AddChannel(info.SessionID, channel)
		d.chefs[info.ChefID] = onlineChef
	}
}

// RemoveSession removes a session from the Distributor.
func (d *Distributor) RemoveSession(info *chefs.SessionInfo) error {
	if info == nil {
		return nil
	}

	d.Lock()
	defer d.Unlock()

	onlineChef, has := d.chefs[info.ChefID]
	if has == false {
		return ErrOnlineChefDoesNotExist
	}

	// Remove channel from OnlineChef. If no more channels exist in OnlineChef,
	// remove the OnlineChef instance as the chef is no longer online.
	if onlineChef.RemoveChannel(info.SessionID) == 0 {
		delete(d.chefs, info.SessionID)
	}
	return nil
}

// SendGlobalMessage sends a message to all chefs.
func (d *Distributor) SendGlobalMessage(msg *conn.Message) {
	d.Lock()
	defer d.Unlock()

	for _, chef := range d.chefs {
		chef.SendMessage(msg)
	}
}
