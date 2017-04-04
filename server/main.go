package main

import (
	"fmt"
	"github.com/evanlinjin/recipe-manager/server/chefs"
	"github.com/evanlinjin/recipe-manager/server/conn"
	"github.com/evanlinjin/recipe-manager/server/handle"
	"github.com/evanlinjin/recipe-manager/server/relay"
	"github.com/gorilla/websocket"
	"net/http"
)

// ObjectGroup groups a bunch of objects together to convenient passing between
// functions. One object group is shared between all the chefs.
type ObjectGroup struct {
	Upgrader    *websocket.Upgrader
	ChefsDB     *chefs.ChefsDB // This is shared via reference across all chefs.
	Distributor *relay.Distributor
}

// MakeObjectGroup makes a new instance of ObjectGroup.
// One object group per chef.
func MakeObjectGroup() (g ObjectGroup, e error) {

	g.Upgrader = &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	cdb, e := chefs.MakeChefsDB(&chefs.Config{
		DomainName:  "http://localhost:8080",
		MongoUrls:   "127.0.0.1:32017",
		BotEmail:    "noreply.recipemanager.io@gmail.com",
		BotEmailPwd: "",
	})
	g.ChefsDB = &cdb

	g.Distributor = relay.NewDistributor()

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

		// Some debug messages.
		fmt.Println("Connection established with", r.RemoteAddr)
		defer fmt.Println("Connection with", r.RemoteAddr, "closed")

		// Some personal session objects.
		var personalSessionInfo *chefs.SessionInfo
		var personalChannel = make(chan *conn.Message)

		// Cleanup code.
		defer g.Distributor.RemoveSession(personalSessionInfo) // handles nil.

		// Handle outgoing messages.
		go func() {
			for {
				select {
				case m := <-personalChannel:
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
					return

				default:
					// Print error on undetermined error.
					fmt.Sprintf("from %v, got %v", r.RemoteAddr, e)
					return
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
							ws.SendResponseMessage(m, msg)
							return
						}
						if e := g.ChefsDB.AddChef(d.Email, d.Password); e != nil {
							msg := fmt.Sprintf("invalid request: %v", e)
							ws.SendResponseMessage(m, msg)
							return
						}
						msg := fmt.Sprintf("please check your email to activate your account for %v", d.Email)
						ws.SendResponseMessage(m, msg)

					case "login":
						d, e := handle.GetLoginData(m.Data)
						if e != nil {
							msg := "invalid data structure"
							ws.SendResponseMessage(m, msg)
							return
						}
						sessionInfo, e := g.ChefsDB.NewSession(d.Email, d.Password)
						if e != nil {
							msg := handle.DealNewSessionError(e)
							ws.SendResponseMessage(m, msg)
							return
						}
						personalSessionInfo = sessionInfo
						g.Distributor.AddSession(personalSessionInfo, personalChannel)
						ws.SendResponseMessage(m, personalSessionInfo)

					case "logout":
						if personalSessionInfo == nil {
							msg := "no session was active"
							fmt.Printf("[%s] %s\n", r.RemoteAddr, msg)
							ws.SendResponseMessage(m, msg)
							return
						}
						e := g.ChefsDB.DeleteSession(personalSessionInfo)
						if e != nil {
							fmt.Printf("[%s] %s : %v\n", r.RemoteAddr, "logout", e)
							ws.SendResponseMessage(m, e.Error())
						} else {
							ws.SendResponseMessage(m, "Logged out!")
						}
						return

					case "claim_session":
						d, e := handle.GetClaimSessionData(m.Data)
						if e != nil {
							msg := "invalid data structure"
							ws.SendResponseMessage(m, msg)
							return
						}
						sessionInfo, e := g.ChefsDB.ClaimSession(
							d.ChefID, d.SessionID, d.SessionKey)
						if e != nil {
							msg := handle.DealClaimSessionError(e)
							ws.SendResponseMessage(m, msg)
							return
						}
						personalSessionInfo = sessionInfo
						g.Distributor.AddSession(personalSessionInfo, personalChannel)
						ws.SendResponseMessage(m, personalSessionInfo)
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
