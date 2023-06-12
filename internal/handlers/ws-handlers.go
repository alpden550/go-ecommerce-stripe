package handlers

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

type WebSocketConnection struct {
	*websocket.Conn
}

type WSPayload struct {
	Action      string               `json:"action"`
	UserName    string               `json:"user_name"`
	Message     string               `json:"message"`
	MessageType string               `json:"message_type"`
	Conn        *WebSocketConnection `json:"-"`
}

type WSJsonResponse struct {
	Action  string `json:"action"`
	Message string `json:"message"`
	UserID  int    `json:"user_id"`
}

var upgradeConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(request *http.Request) bool { return true },
}

var clients = make(map[WebSocketConnection]string)

var wsChan = make(chan WSPayload)

func WsHandler(writer http.ResponseWriter, request *http.Request) {
	ws, err := upgradeConnection.Upgrade(writer, request, nil)
	if err != nil {
		app.ErrorLog.Printf("%w", fmt.Errorf("%e", err))
		return
	}
	app.InfoLog.Println(request.RemoteAddr)

	response := WSJsonResponse{
		Message: "Connected to server",
	}
	err = ws.WriteJSON(response)
	if err != nil {
		app.ErrorLog.Printf("%w", fmt.Errorf("%e", err))
		return
	}

	conn := WebSocketConnection{ws}
	clients[conn] = ""

	go listenForWS(&conn)
}

func listenForWS(conn *WebSocketConnection) {
	defer func() {
		if r := recover(); r != nil {
			app.ErrorLog.Printf("ERROR: %v", r)
		}
	}()

	var payload WSPayload

	for {
		err := conn.ReadJSON(&payload)
		if err != nil {
			app.ErrorLog.Printf("%w", fmt.Errorf("%e", err))
		} else {
			payload.Conn = conn
			wsChan <- payload
		}
	}
}

func ListenToWsChannel() {
	var response WSJsonResponse

	for {
		e := <-wsChan
		switch e.Action {
		case "deleteUser":
			response.Action = "logout"
			response.Message = "Your account has been deleted"
			broadcastToAll(response)

		default:
		}
	}
}

func broadcastToAll(response WSJsonResponse) {
	for client := range clients {
		err := client.WriteJSON(response)
		if err != nil {
			app.ErrorLog.Printf("%w", fmt.Errorf("%e", err))
			_ = client.Close()
			delete(clients, client)
		}
	}
}
