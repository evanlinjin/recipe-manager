package main

import (
	"fmt"
	"net/http"

	"github.com/evanlinjin/recipe-manager/server/talkrelay"
	"github.com/gorilla/websocket"
	"github.com/kabukky/httpscerts"
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
		wsm, e := talkrelay.MakeWSManager(upgrader, w, r)
		if e != nil {
			fmt.Println(e)
			return
		}

		fmt.Println("Connection established with", r.RemoteAddr)
		defer fmt.Println("Connection with", r.RemoteAddr, "closed")

		quitChan := make(chan int)
		defer func() { quitChan <- 1 }()

		msgChan, e := talkGroup.AddChef(r.RemoteAddr)
		if e != nil {
			return
		}
		defer talkGroup.RmChef(r.RemoteAddr)

		go func() {
			for {
				select {
				case m := <-msgChan:
					//wsm.DeprecatedWriteMessage([]byte(m))
					wsm.SendRequestMessage("unspecified", m)
				case <-quitChan:
					return
				}
			}
		}()

		for {
			data, e := wsm.DeprecatedReadMessage()
			if e != nil {
				fmt.Println(e)

				return
			}
			go talkGroup.Talk(fmt.Sprintf("%s", data))
		}
	}
}
