package events

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strings"
	"sync"
)

type SocketMessage struct {
	Action    string `json:"action,omitempty"`
	Object    string `json:"object,omitempty"`
	Data      string `json:"data"`
	ErrorCode int64  `json:"error_code"`
	Error     string `json:"error,omitempty"`
}

type SocketChannel struct {
	conn *websocket.Conn
	send chan SocketMessage
}

var clients = make(map[*websocket.Conn]bool)
var Broadcast = make(chan SocketMessage)
var mux sync.Mutex

func init() {
	go broadcast()
}
func Upgrader(w http.ResponseWriter, r *http.Request) {
	var u = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	conn, _ := u.Upgrade(w, r, nil)
	ch := NewChannel(conn)
	ch.send <- SocketMessage{Data: "Hello", ErrorCode: 0}
	Broadcast <- SocketMessage{Data: "New client is connected...", ErrorCode: 0}
}

func NewChannel(conn *websocket.Conn) SocketChannel {
	c := SocketChannel{
		conn: conn,
		send: make(chan SocketMessage, 5),
	}

	go c.reader()
	go c.writer()
	clients[conn] = true

	return c
}

func (c SocketChannel) reader() {
	var msg SocketMessage
	defer c.conn.Close()

	for {
		err := c.conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("reader error: %v", err)
			c.conn.Close()
			delete(clients, c.conn)
			break
		}
		log.Printf("Received Message: %v", msg)

		switch strings.ToUpper(msg.Action) {
		case "BROADCAST":
			Broadcast <- SocketMessage{Data: msg.Data, Action: msg.Action, ErrorCode: 0}
		case "ECHO":
			c.conn.WriteJSON(msg)
		}
	}
}

func (c SocketChannel) writer() {
	for msg := range c.send {
		mux.Lock()
		c.conn.WriteJSON(msg)
		mux.Unlock()
	}
}

func broadcast() {
	for {
		// Grab the next message from the Broadcast channel
		msg := <-Broadcast
		mux.Lock()
		Evt.msg <- msg.Data + " *** BROADCAST ***"
		mux.Unlock()
		// Send it out to every client that is currently connected
		for client := range clients {
			mux.Lock()
			err := client.WriteJSON(msg)
			mux.Unlock()
			if err != nil {
				log.Printf("broadcast() error : %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{}

	// Upgrade initial GET request to a websocket
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, "Error WebSocket Upgrade: "+err.Error(), http.StatusInternalServerError)
	}
	defer ws.Close()

	// Register our new client
	clients[ws] = true
	Broadcast <- SocketMessage{}

	for {
		var msg SocketMessage
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("HandleConnections() error: %v", err)
			// delete(clients, ws)
			// break
		} else {
			Broadcast <- msg
		}
	}
}
