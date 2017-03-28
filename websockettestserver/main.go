package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/kabukky/httpscerts"
	"github.com/evanlinjin/recipe-manager/server/talkrelay"
)

func main() {

	if httpscerts.Check("cert.pem", "key.pem") != nil {
		httpscerts.Generate("cert.pem", "key.pem", "localhost")
	}

	var (
		upg = websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}
		talkGroup = talkrelay.MakeTalkGroup()
	)

	http.HandleFunc("/", makeHandler(&upg, &talkGroup))

	if e := http.ListenAndServeTLS(":8182", "cert.pem", "key.pem", nil); e != nil {
		panic(e)
	}
}

func makeHandler(upgrader *websocket.Upgrader, talkGroup *talkrelay.TalkGroup) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		wsc, e := upgrader.Upgrade(w, r, nil)
		if e != nil {
			fmt.Println(e)
			return
		}

		readMux, writeMux := sync.Mutex{}, sync.Mutex{}
		quitChan := make(chan int)

		fmt.Println("Connection established with", r.RemoteAddr)
		defer fmt.Println("Connection with", r.RemoteAddr, "closed")
		defer func() { quitChan <- 1 }()

		go func() {
			msgChan := make(chan string)
			talkGroup.AddChef(r.RemoteAddr, msgChan)
			defer talkGroup.RmChef(r.RemoteAddr)

			for {
				select {
				case m := <-msgChan:
					writeMux.Lock()
					wsc.WriteMessage(websocket.TextMessage, []byte(m))
					writeMux.Unlock()

				case <-quitChan:
					return
				}
			}
		}()

		for {
			readMux.Lock()
			_, data, e := wsc.ReadMessage()
			readMux.Unlock()
			if e != nil {
				fmt.Println(e)
				return
			}
			go talkGroup.Talk(fmt.Sprintf("%s", data))
		}
	}
}
