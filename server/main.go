package main

import (
	"fmt"
	"github.com/evanlinjin/recipe-manager/server/chefs"
	"github.com/evanlinjin/recipe-manager/server/conn"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
	"github.com/evanlinjin/recipe-manager/server/handle"
)

// ObjectGroup groups a bunch of objects together to convenient passing between
// function. Note that members are all references.
type ObjectGroup struct {
	OutgoingMessages chan *conn.Message
	Upgrader         *websocket.Upgrader
	ChefsDB          *chefs.ChefsDB
}

// MakeObjectGroup makes a new instance of ObjectGroup.
func MakeObjectGroup() (g ObjectGroup, e error) {

	g.OutgoingMessages = make(chan *conn.Message)

	g.Upgrader = &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	cdb, e := chefs.MakeChefsDB(chefs.Config{
		DomainName: "http://localhost:8080",
		BotEmail: "noreply.recipemanager.io@gmail.com",
		BotEmailPwd: "",
	})
	g.ChefsDB = &cdb

	return
}

func main() {
	g, e := MakeObjectGroup()
	if e != nil {
		panic(e)
	}

	http.HandleFunc("/action/", chefs.MakeActivationEndpoint(g.ChefsDB))
	http.HandleFunc("/ws", MakeWebSocketEndpoint(g))
	http.ListenAndServe(":8080", nil)
}

// MakeWebSocketEndpoint makes a WebSocket connection endpoint.
func MakeWebSocketEndpoint(g ObjectGroup) func(
	w http.ResponseWriter, r *http.Request) {

	ep := func(w http.ResponseWriter, r *http.Request) {

		// Initiate upgrade from HTTPS to WebSocket connection.
		ws, e := conn.MakeWSManager(g.Upgrader, w, r)
		if e != nil {
			fmt.Println(e)
			return
		}

		// Initiate handshake to agree on encryption key for custom protocol.
		if e := ws.Handshake(time.Second * 3); e != nil {
			fmt.Println(e)
			return
		}

		// Some debug messages.
		fmt.Println("Connection established with", r.RemoteAddr)
		defer fmt.Println("Connection with", r.RemoteAddr, "closed")

		// Handle outgoing messages.
		go func() {
			for {
				select {
				case m := <-g.OutgoingMessages:
					ws.SendMessage(m)

				case <-ws.QuitChan:
					return
				}
			}
		}()

		// Handle incoming messages.
		for {
			m, e := ws.GetMessage()
			if e != nil {

				// Error switch.
				switch e.(type) {

				case *conn.ErrCloseMessage:
					// Close connection on close message.
					ws.Close(0, e.Error())
					return

				case *conn.ErrUnexpectedMessage:
					// Print error on unexpected message.
					fmt.Sprintf("from %v, got %v", r.RemoteAddr, e)

				default:
					// Print error on undetermined error.
					fmt.Sprintf("from %v, got %v", r.RemoteAddr, e)
				}

			} else {

				// Message switch.
				switch m.Type {

				case conn.TypeRequest:
					// For request messages.

					switch m.Command {
					case "new_chef":
						d, e := handle.GetNewChefData(m.Data)
						if e != nil {
							msg := fmt.Sprintf("invalid data structure: %v", e)
							fmt.Printf("[%v] %v", r.RemoteAddr, msg)
							ws.SendResponseMessage(m, msg)
							break
						}
						if e := g.ChefsDB.AddChef(d.Email, d.Password); e != nil {
							msg := fmt.Sprintf("invalid request: %v", e)
							fmt.Printf("[%v] %v", r.RemoteAddr, msg)
							ws.SendResponseMessage(m, msg)
							break
						}
						msg := fmt.Sprintf("please check your email to activate your account for %v", d.Email)
						fmt.Printf("[%v] %v", r.RemoteAddr, msg)
						ws.SendResponseMessage(m, msg)
					}

				case conn.TypeResponse:
					// For response messages.

					switch m.Command {

					}
				}
			}
		}
	}
	return ep
}
