package conn

import (
	"fmt"
	"sync"
)

// TalkGroup represents a group of chefs that share the same talk.
type TalkGroup struct {
	sync.Mutex
	chefs map[string]chan string
}

// MakeTalkGroup makes a new TalkGroup.
func MakeTalkGroup() (m TalkGroup) {
	m.chefs = make(map[string]chan string)
	return
}

// AddMember adds a chef to TalkGroup.
func (g *TalkGroup) AddChef(chefID string) (chefChan chan string, e error) {
	g.Lock()
	defer g.Unlock()
	chefChan = make(chan string)
	if _, got := g.chefs[chefID]; got {
		e = fmt.Errorf("an entry already exists for '%s'", chefID)
		return
	}
	g.chefs[chefID] = chefChan
	fmt.Printf("[TalkGroup] CHEF COUNT: %v \n", len(g.chefs))
	return
}

// RmChef removes a chef from TalkGroup.
func (g *TalkGroup) RmChef(chefID string) {
	g.Lock()
	defer g.Unlock()
	delete(g.chefs, chefID)
	fmt.Printf("[TalkGroup] CHEF COUNT: %v \n", len(g.chefs))
}

// HeadCount counts the number of chef heads in the group.
func (g *TalkGroup) HeadCount() int {
	g.Lock()
	defer g.Unlock()
	return len(g.chefs)
}

// Talk talks to all chefs of TalkGroup.
func (g *TalkGroup) Talk(msg string) {
	g.Lock()
	defer g.Unlock()
	for _, v := range g.chefs {
		v <- msg
	}
}
