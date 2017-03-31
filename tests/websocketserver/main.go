package main

import (
	"fmt"
	"net/http"

	"github.com/evanlinjin/recipe-manager/server/conn"
	"github.com/gorilla/websocket"
	"github.com/kabukky/httpscerts"
	"time"
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
		talkGroup = conn.MakeTalkGroup()
	)

	http.HandleFunc("/", makeHandler(&upg, &talkGroup))

	if e := http.ListenAndServeTLS(":8182", "cert.pem", "key.pem", nil); e != nil {
		panic(e)
	}
}

func makeHandler(upgrader *websocket.Upgrader, talkGroup *conn.TalkGroup) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		wsm, e := conn.MakeWSManager(upgrader, w, r)
		if e != nil {
			fmt.Println(e)
			return
		}
		if e := wsm.Handshake(time.Second * 3); e != nil {
			fmt.Println(e)
			return
		}

		fmt.Println("Connection established with", r.RemoteAddr)
		defer fmt.Println("Connection with", r.RemoteAddr, "closed")

		//quitChan := make(chan int)
		//defer func() { quitChan <- 1 }()

		msgChan, e := talkGroup.AddChef(r.RemoteAddr)
		if e != nil {
			wsm.Close(-1, e.Error())
			return
		}
		defer talkGroup.RmChef(r.RemoteAddr)

		go func() {
			for {
				select {
				case m := <-msgChan:
					wsm.SendRequestMessage("unspecified", m)
				case <-wsm.QuitChan:
					return
				}
			}
		}()

		for {
			msg, e := wsm.GetMessage()
			if e != nil {
				fmt.Println(e)
				wsm.Close(-1, e.Error())
				return
			}
			if msg == nil {
				break
			}
			go talkGroup.Talk(fmt.Sprintf("%s", msg.Data))
		}
		wsm.Close(0, "")
	}
}
