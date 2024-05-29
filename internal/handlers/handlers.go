package handlers

import (
	"fmt"
	"log"
	"net/http"
	"sort"

	"github.com/CloudyKit/jet/v6"

	"github.com/gorilla/websocket"
)

var wsChan = make(chan WSJsonPayload)

var clients = make(map[WebsocketConnection]string)

var views = jet.NewSet(jet.NewOSFileSystemLoader("./html"), jet.InDevelopmentMode())
var upgradeConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Home(w http.ResponseWriter, r *http.Request) {
	err := renderPage(w, "home.jet", nil)

	if err != nil {
		log.Print(err)
	}
}

func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	ws, err := upgradeConnection.Upgrade(w, r, nil)
	log.Print("Client connected to endpoint")
	if err != nil {
		log.Print(err)
		return
	}
	var response WSJsonResponse

	response.Message = `<em><small>Connected to server</em></small>`

	conn := WebsocketConnection{
		Conn: ws,
	}

	clients[conn] = ""
	err = ws.WriteJSON(response)

	if err != nil {
		log.Print(err)
		return
	}

	go ListenForWS(&conn)

}

func ListenForWS(conn *WebsocketConnection) {

	defer func() {

		if r := recover(); r != nil {
			log.Print("Error", fmt.Sprintf("%v", r))
		}
	}()

	var payload WSJsonPayload

	for {
		err := conn.ReadJSON(&payload)

		if err != nil {
			log.Printf("error occured reading json %s", err)
		} else {
			payload.conn = *conn
			wsChan <- payload
		}

	}
}

func broadCastToAll(response WSJsonResponse) {
	for client := range clients {
		err := client.WriteJSON(response)

		if err != nil {
			log.Println("webscoket err")
			_ = client.Close()
			delete(clients, client)
		}
	}
}

func ListenToWsChannel() {
	var response WSJsonResponse

	for {
		e := <-wsChan

		switch e.Action {
		case "username":
			// get a list of usernames and send it back
			clients[e.conn] = e.Username
			users := getUserList()
			response.Action = "list_users"
			response.ConnectedUsers = users
			broadCastToAll(response)
		case "left":
			response.Action = "list_users"
			delete(clients,e.conn)
			users := getUserList()
			response.ConnectedUsers = users
			broadCastToAll(response)
		}

	}
}

func getUserList() []string {
	var userLists = []string{}
	for _, x := range clients {
		userLists = append(userLists, x)
	}

	sort.Strings(userLists)
	return userLists
}

type WebsocketConnection struct {
	*websocket.Conn
}

type WSJsonResponse struct {
	Action         string   `json:"action"`
	Message        string   `json:"message"`
	MessageType    string   `json:"message_type"`
	ConnectedUsers []string `json:"connected_users"`
}

type WSJsonPayload struct {
	Action   string              `json:"action"`
	Message  string              `json:"message"`
	Username string              `json:"username"`
	conn     WebsocketConnection `json:"-"`
}

func renderPage(w http.ResponseWriter, tmpl string, data jet.VarMap) error {
	view, err := views.GetTemplate(tmpl)

	if err != nil {
		log.Println(err)
		return err
	}
	err = view.Execute(w, data, nil)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
