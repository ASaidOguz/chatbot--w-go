package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/CloudyKit/jet/v6"
	"github.com/gorilla/websocket"
)

// WsChan is a channel that holds only websocket payloads .
var wsChan = make(chan WsPayload)

// clients is a map that holds all the connected clients.
var clients = make(map[WebSocketConnection]string)

// views is a set of Jet templates that are loaded from the ./html directory.
var views = jet.NewSet(
	jet.NewOSFileSystemLoader("./html"),
	jet.InDevelopmentMode(),
)

// UpgradeConnection is a websocket.Upgrader that allows us
// to upgrade a HTTP connection to a websocket connection.
var UpgradeConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Home is a http.HandlerFunc that renders the home page.(this is public function
func Home(w http.ResponseWriter, r *http.Request) {
	err := renderPage(w, "home.jet", nil)
	if err != nil {
		log.Default().Println(err)
	}
}

type WebSocketConnection struct {
	*websocket.Conn
}

// WsJsonResponjse is a struct that represents a JSON response from the websocket server.
type WsJsonResponse struct {
	Action      string `json:"action"`
	Message     string `json:"message"`
	MessageType string `json:"message_type"`
}

type WsPayload struct {
	Action   string              `json:"action"`
	UserName string              `json:"username"`
	Message  string              `json:"message"`
	Conn     WebSocketConnection `json:"-"`
}

// WsEndPoint is a http.HandlerFunc that upgrades a HTTP connection to a websocket connection.
func WsEndPoint(w http.ResponseWriter, r *http.Request) {
	ws, err := UpgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		log.Default().Println(err)
	}
	log.Default().Println("Client Connected to Endpoint")

	response := &WsJsonResponse{
		Message: `<em><small>Connected to Server</small></em>`,
	}

	conn := WebSocketConnection{Conn: ws}
	clients[conn] = ""

	err = ws.WriteJSON(response)
	if err != nil {
		log.Default().Println(err)
	}
	// Create go-routine for listening for incoming messages
	go listenForWs(&conn)
}

// listenForWs is a function that listens for incoming messages from a websocket connection.(private function)
func listenForWs(conn *WebSocketConnection) {
	defer func() {
		if err := recover(); err != nil {
			log.Default().Println("Error:", err)
		}
	}()

	var payLoad WsPayload

	for {
		err := conn.ReadJSON(&payLoad)
		if err != nil {
			log.Default().Println(err)

		} else {
			payLoad.Conn = *conn
			wsChan <- payLoad
		}
	}
}

func ListenToWsChannel() {
	var response WsJsonResponse

	for {
		e := <-wsChan
		response.Action = "Got here"
		response.Message = fmt.Sprintf("Some Message and Action was: %s", e.Action)
		broadcastToAll(response)
	}
}

func broadcastToAll(response WsJsonResponse) {
	for client := range clients {
		err := client.WriteJSON(response)
		if err != nil {
			log.Default().Println("Error:", err)
			_ = client.Close()
			delete(clients, client)
		}
	}
}

// renderPage is a helper function that renders a Jet template.(private function)
func renderPage(w http.ResponseWriter, tmpl string, data jet.VarMap) error {
	view, err := views.GetTemplate(tmpl)
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return err
	}
	err = view.Execute(w, data, nil)
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return err
	}
	return nil
}
