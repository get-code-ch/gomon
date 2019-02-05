package events

import (
	"controller/authorize"
	"controller/command"
	"controller/history"
	"controller/host"
	"controller/probe"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type SocketMessage struct {
	Action    string `json:"action,omitempty"`
	Object    string `json:"object,omitempty"`
	Data      string `json:"data"`
	ErrorCode int64  `json:"error_code"`
	Error     string `json:"error,omitempty"`
	ClientId  int64  `json:"client_id,omitempty"`
	Username  string `json:"username,omitempty"`
	Token     string `json:"token,omitempty"`
}

type SocketChannel struct {
	conn *websocket.Conn
	send chan SocketMessage
}

var clients = make(map[*websocket.Conn]int64)
var authClients = make(map[int64]bool)
var Broadcast = make(chan SocketMessage)
var Mux sync.Mutex

func init() {
	go broadcast()
}

// Upgrade http request to websocket
func Upgrader(w http.ResponseWriter, r *http.Request) {
	var u = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	conn, _ := u.Upgrade(w, r, nil)
	c, cid := NewChannel(conn)
	c.conn.WriteJSON(SocketMessage{Action: "AUTHENTICATE", Data: "Please login", ErrorCode: 0, ClientId: cid})
}

// Create a new socket channel
func NewChannel(conn *websocket.Conn) (SocketChannel, int64) {
	c := SocketChannel{
		conn: conn,
		send: make(chan SocketMessage, 200),
	}

	go c.reader()

	var cid = int64(0)
	for {
		var exist = false
		cid = rand.Int63n(999999999)
		for _, v := range clients {
			if v == cid {
				exist = true
			}
		}
		if exist {
			continue
		}
		break
	}
	clients[conn] = cid

	// By default client is not authenticated
	authClients[cid] = false
	return c, cid
}

// Reading receiving message on Sockets
func (c SocketChannel) reader() {
	var msg SocketMessage
	defer c.conn.Close()

	for {
		err := c.conn.ReadJSON(&msg)
		if err != nil {
			cid := clients[c.conn]
			c.conn.Close()
			// LockManager release editing locks for this clients
			UnlockByClientId(cid)
			delete(clients, c.conn)
			Broadcast <- SocketMessage{Action: "LOG", Data: "Client " + strconv.FormatInt(cid, 10) + " as gone", ErrorCode: 0}
			break
		}
		log.Printf("Received Message: %v", msg)

		if strings.ToUpper(msg.Action) == "AUTHENTICATE" {
			auth, err := authorize.ValidateSocketToken(msg.Token)
			if err != nil {
				c.conn.WriteJSON(SocketMessage{Action: "LOG",
					Data: time.Now().Format(time.RFC3339) + ": Authentication Error", ErrorCode: 1, Error: err.Error()})
			} else {
				if auth {
					authClients[msg.ClientId] = true
					c.conn.WriteJSON(SocketMessage{Action: "LOG", Data: "Successfull authentication",
						ErrorCode: 0, Token: msg.Token, ClientId: msg.ClientId})
					continue
				} else {
					c.conn.WriteJSON(SocketMessage{Action: "LOG",
						Data: time.Now().Format(time.RFC3339) + ": Authentication Error", ErrorCode: 1, Error: err.Error()})
				}
			}
		}

		if !authClients[clients[c.conn]] {
			if ok, _ := authorize.ValidateSocketToken(msg.Token); !ok {
				c.conn.WriteJSON(SocketMessage{Action: "LOG", Data: "Error unauthorized Access", ErrorCode: 0, ClientId: clients[c.conn]})
				continue
			} else {
				authClients[clients[c.conn]] = true
			}
		}

		switch strings.ToUpper(msg.Object) {
		case "PROBE":
			handleProbeMessage(msg, c)
		case "HOST":
			handleHostMessage(msg, c)
		case "COMMAND":
			handleCommandMessage(msg, c)
		case "HISTORY":
			handleHistoryMessage(msg, c)
		default:
			switch strings.ToUpper(msg.Action) {
			case "BROADCAST":
				Broadcast <- SocketMessage{Data: msg.Data, Action: "LOG", ErrorCode: 0}
			case "ECHO":
				c.conn.WriteJSON(SocketMessage{Data: msg.Data, Action: "LOG", ErrorCode: 0})
			}
		}
	}
}

func broadcast() {
	for {
		// Grab the next message from the Broadcast channel
		msg := <-Broadcast
		for client, key := range clients {
			if authClients[key] {
				Mux.Lock()
				log.Printf("%v", key)
				m := SocketMessage{msg.Action, msg.Object, msg.Data, msg.ErrorCode, msg.Error, key, "", ""}
				err := client.WriteJSON(m)
				Mux.Unlock()
				if err != nil {
					log.Printf("broadcast() error : %v", err)
					client.Close()
					delete(clients, client)
				}
			}
		}
	}
}

func handleProbeMessage(msg SocketMessage, sc SocketChannel) {
	var p = probe.Probe{}
	action := strings.ToUpper(msg.Action)

	if action != "GET" {
		err := json.Unmarshal([]byte(msg.Data), &p)
		if err != nil {
			sc.conn.WriteJSON(SocketMessage{Object: "PROBE", Action: "LOG",
				Data: time.Now().Format(time.RFC3339) + ": Error Unmarshaling probe data", ErrorCode: 1, Error: err.Error()})
			return
		}
	}

	switch action {
	case "GET":
		probes, err := p.Get()
		if err != nil {
			sc.conn.WriteJSON(SocketMessage{Object: "PROBE", Action: "LOG",
				Data: time.Now().Format(time.RFC3339) + ": Error getting probe", ErrorCode: 0})
		} else {
			// Check locked items
			for i, p := range probes {
				probes[i].Locked = IsLocked(p.Id)
			}

			j, _ := json.Marshal(probes)
			Broadcast <- SocketMessage{Object: "PROBE", Action: "GET",
				Data: string(j), ErrorCode: 0}
		}
	case "CREATE":
		if p.Interval < 1 {
			p.Interval = 99999
		}
		Mux.Lock()
		err := p.Post()
		Mux.Unlock()
		if err != nil {
			sc.conn.WriteJSON(SocketMessage{Object: "PROBE", Action: "LOG",
				Data: time.Now().Format(time.RFC3339) + ": Error Creating probe -->" + p.Name, ErrorCode: 1, Error: err.Error()})
		} else {
			sc.conn.WriteJSON(SocketMessage{Object: "PROBE", Action: "LOG",
				Data: time.Now().Format(time.RFC3339) + ": Probe " + p.Name + " created", ErrorCode: 0})
			j, _ := json.Marshal(p)
			Broadcast <- SocketMessage{Object: "PROBE", Action: "CREATE",
				Data: string(j), ErrorCode: 0}
			w := Worker{
				Ticker: time.NewTicker(time.Millisecond * 1000 * time.Duration(p.Interval)),
				Probe:  p,
			}
			Workers[p.Id] = w
			go Workers[p.Id].Timer()
		}
	case "UPDATE":
		if p.Interval < 1 {
			p.Interval = 99999
		}
		Mux.Lock()
		err := p.Put()
		Mux.Unlock()
		if err != nil {
			sc.conn.WriteJSON(SocketMessage{Object: "PROBE", Action: "LOG",
				Data: time.Now().Format(time.RFC3339) + ": Error updating probe -->" + p.Name, ErrorCode: 1, Error: err.Error()})
		} else {
			sc.conn.WriteJSON(SocketMessage{Object: "PROBE", Action: "LOG",
				Data: time.Now().Format(time.RFC3339) + ": Probe " + p.Name + " updated", ErrorCode: 0})
			j, _ := json.Marshal(p)
			Broadcast <- SocketMessage{Object: "PROBE", Action: "UPDATE",
				Data: string(j), ErrorCode: 0}
			w := Worker{
				Ticker: time.NewTicker(time.Millisecond * 1000 * time.Duration(p.Interval)),
				Probe:  p,
			}
			Workers[p.Id].Ticker.Stop()
			Workers[p.Id] = w
			go Workers[p.Id].Timer()
		}
	case "DELETE":
		Mux.Lock()
		err := p.Delete()
		Mux.Unlock()
		if err != nil {
			sc.conn.WriteJSON(SocketMessage{Object: "PROBE", Action: "LOG",
				Data: time.Now().Format(time.RFC3339) + ": Error deleting probe -->" + p.Name, ErrorCode: 1, Error: err.Error()})
		} else {
			sc.conn.WriteJSON(SocketMessage{Object: "PROBE", Action: "LOG",
				Data: time.Now().Format(time.RFC3339) + ": Probe " + p.Name + " deleted", ErrorCode: 0})
			Workers[p.Id].Ticker.Stop()
			delete(Workers, p.Id)
			j, _ := json.Marshal(p)
			Broadcast <- SocketMessage{Object: "PROBE", Action: "DELETE",
				Data: string(j), ErrorCode: 0}
		}
	case "LOCK", "UNLOCK":
		l := Lock{time.Now(), "PROBE", p.Id, clients[sc.conn]}
		if action == "LOCK" {
			AddLock(l)
		} else {
			RemoveLock(p.Id)
		}
		Broadcast <- SocketMessage{Object: "PROBE", Action: action, Data: msg.Data, ErrorCode: 0}
	}
}

func handleHostMessage(msg SocketMessage, sc SocketChannel) {
	h := host.Host{}
	action := strings.ToUpper(msg.Action)

	if action != "GET" {
		err := json.Unmarshal([]byte(msg.Data), &h)
		if err != nil {
			sc.conn.WriteJSON(SocketMessage{Object: "HOST", Action: "LOG",
				Data: time.Now().Format(time.RFC3339) + ": Error Unmarshaling host data", ErrorCode: 1, Error: err.Error()})
			return
		}
	}

	switch action {
	case "GET":
		hosts, err := h.Get()
		if err != nil {
			sc.conn.WriteJSON(SocketMessage{Object: "HOST", Action: "LOG",
				Data: time.Now().Format(time.RFC3339) + ": Error getting host", ErrorCode: 0})
		} else {
			// Check locked items
			for i, h := range hosts {
				hosts[i].Locked = IsLocked(h.Id)
			}

			j, _ := json.Marshal(hosts)
			Broadcast <- SocketMessage{Object: "HOST", Action: "GET",
				Data: string(j), ErrorCode: 0}
		}
	case "CREATE":
		Mux.Lock()
		err := h.Post()
		Mux.Unlock()
		if err != nil {
			sc.conn.WriteJSON(SocketMessage{Object: "HOST", Action: "LOG",
				Data: time.Now().Format(time.RFC3339) + ": Error Creating host -->" + h.Name, ErrorCode: 1, Error: err.Error()})
		} else {
			sc.conn.WriteJSON(SocketMessage{Object: "HOST", Action: "LOG",
				Data: time.Now().Format(time.RFC3339) + ": Host " + h.Name + " created", ErrorCode: 0})
			j, _ := json.Marshal(h)
			Broadcast <- SocketMessage{Object: "HOST", Action: "CREATE",
				Data: string(j), ErrorCode: 0}
		}
	case "UPDATE":
		Mux.Lock()
		err := h.Put()
		Mux.Unlock()
		if err != nil {
			sc.conn.WriteJSON(SocketMessage{Object: "HOST", Action: "LOG",
				Data: time.Now().Format(time.RFC3339) + ": Error updating host -->" + h.Name, ErrorCode: 1, Error: err.Error()})
		} else {
			sc.conn.WriteJSON(SocketMessage{Object: "HOST", Action: "LOG",
				Data: time.Now().Format(time.RFC3339) + ": Probe " + h.Name + " updated", ErrorCode: 0})
			j, _ := json.Marshal(h)
			Broadcast <- SocketMessage{Object: "HOST", Action: "UPDATE",
				Data: string(j), ErrorCode: 0}
		}
	case "DELETE":
		Mux.Lock()
		err := h.Delete()
		Mux.Unlock()
		if err != nil {
			sc.conn.WriteJSON(SocketMessage{Object: "HOST", Action: "LOG",
				Data: time.Now().Format(time.RFC3339) + ": Error deleting host -->" + h.Name, ErrorCode: 1, Error: err.Error()})
		} else {
			sc.conn.WriteJSON(SocketMessage{Object: "HOST", Action: "LOG",
				Data: time.Now().Format(time.RFC3339) + ": Probe " + h.Name + " deleted", ErrorCode: 0})
			j, _ := json.Marshal(h)
			Broadcast <- SocketMessage{Object: "HOST", Action: "DELETE",
				Data: string(j), ErrorCode: 0}
		}
	case "LOCK", "UNLOCK":
		l := Lock{time.Now(), "HOST", h.Id, clients[sc.conn]}
		if action == "LOCK" {
			AddLock(l)
		} else {
			RemoveLock(h.Id)
		}
		Broadcast <- SocketMessage{Object: "HOST", Action: action, Data: msg.Data, ErrorCode: 0}
	}
}

func handleCommandMessage(msg SocketMessage, sc SocketChannel) {
	c := command.Command{}
	action := strings.ToUpper(msg.Action)

	if action != "GET" {
		err := json.Unmarshal([]byte(msg.Data), &c)
		if err != nil {
			sc.conn.WriteJSON(SocketMessage{Object: "COMMAND", Action: "LOG",
				Data: time.Now().Format(time.RFC3339) + ": Error Unmarshaling command data", ErrorCode: 1, Error: err.Error()})
			return
		}
	}

	switch action {
	case "GET":
		commands, err := c.Get()
		if err != nil {
			sc.conn.WriteJSON(SocketMessage{Object: "COMMAND", Action: "LOG",
				Data: time.Now().Format(time.RFC3339) + ": Error getting host", ErrorCode: 0})
		} else {
			// Check locked items
			for i, c := range commands {
				commands[i].Locked = IsLocked(c.Id)
			}

			j, _ := json.Marshal(commands)
			Broadcast <- SocketMessage{Object: "COMMAND", Action: "GET",
				Data: string(j), ErrorCode: 0}
		}
	case "CREATE":
		Mux.Lock()
		err := c.Post()
		Mux.Unlock()
		if err != nil {
			sc.conn.WriteJSON(SocketMessage{Object: "COMMAND", Action: "LOG",
				Data: time.Now().Format(time.RFC3339) + ": Error Creating command -->" + c.Name, ErrorCode: 1, Error: err.Error()})
		} else {
			sc.conn.WriteJSON(SocketMessage{Object: "COMMAND", Action: "LOG",
				Data: time.Now().Format(time.RFC3339) + ": Command " + c.Name + " created", ErrorCode: 0})
			j, _ := json.Marshal(c)
			Broadcast <- SocketMessage{Object: "COMMAND", Action: "CREATE",
				Data: string(j), ErrorCode: 0}
		}
	case "UPDATE":
		Mux.Lock()
		err := c.Put()
		Mux.Unlock()
		if err != nil {
			sc.conn.WriteJSON(SocketMessage{Object: "COMMAND", Action: "LOG",
				Data: time.Now().Format(time.RFC3339) + ": Error updating command -->" + c.Name, ErrorCode: 1, Error: err.Error()})
		} else {
			sc.conn.WriteJSON(SocketMessage{Object: "COMMAND", Action: "LOG",
				Data: time.Now().Format(time.RFC3339) + ": Command " + c.Name + " updated", ErrorCode: 0})
			j, _ := json.Marshal(c)
			Broadcast <- SocketMessage{Object: "COMMAND", Action: "UPDATE",
				Data: string(j), ErrorCode: 0}
		}
	case "DELETE":
		Mux.Lock()
		err := c.Delete()
		Mux.Unlock()
		if err != nil {
			sc.conn.WriteJSON(SocketMessage{Object: "COMMAND", Action: "LOG",
				Data: time.Now().Format(time.RFC3339) + ": Error deleting host -->" + c.Name, ErrorCode: 1, Error: err.Error()})
		} else {
			sc.conn.WriteJSON(SocketMessage{Object: "COMMAND", Action: "LOG",
				Data: time.Now().Format(time.RFC3339) + ": Probe " + c.Name + " deleted", ErrorCode: 0})
			j, _ := json.Marshal(c)
			Broadcast <- SocketMessage{Object: "COMMAND", Action: "DELETE",
				Data: string(j), ErrorCode: 0}
		}
	case "LOCK", "UNLOCK":
		l := Lock{time.Now(), "COMMAND", c.Id, clients[sc.conn]}
		if action == "LOCK" {
			AddLock(l)
		} else {
			RemoveLock(c.Id)
		}
		Broadcast <- SocketMessage{Object: "COMMAND", Action: action, Data: msg.Data, ErrorCode: 0}
	}
}

func handleHistoryMessage(msg SocketMessage, sc SocketChannel) {
	h := history.History{}
	switch strings.ToUpper(msg.Action) {
	case "GET":
		histories, err := h.Get()
		if err != nil {
			sc.conn.WriteJSON(SocketMessage{Object: "HISTORY", Action: "LOG",
				Data: time.Now().Format(time.RFC3339) + ": Error getting host", ErrorCode: 0})
		} else {
			j, _ := json.Marshal(histories)
			Broadcast <- SocketMessage{Object: "HISTORY", Action: "GET",
				Data: string(j), ErrorCode: 0}
		}
	case "DELETE":
		err := json.Unmarshal([]byte(msg.Data), &h)
		if err != nil {
			sc.conn.WriteJSON(SocketMessage{Object: "HISTORY", Action: "LOG",
				Data: time.Now().Format(time.RFC3339) + ": Error Unmarshaling host data", ErrorCode: 1, Error: err.Error()})
		} else {
			Mux.Lock()
			err := h.Delete()
			Mux.Unlock()
			if err != nil {
				sc.conn.WriteJSON(SocketMessage{Object: "HISTORY", Action: "LOG",
					Data: time.Now().Format(time.RFC3339) + ": Error deleting history", ErrorCode: 1, Error: err.Error()})
			} else {
				sc.conn.WriteJSON(SocketMessage{Object: "HISTORY", Action: "LOG",
					Data: time.Now().Format(time.RFC3339) + ": History record deleted", ErrorCode: 0})
				j, _ := json.Marshal(h)
				Broadcast <- SocketMessage{Object: "HISTORY", Action: "DELETE",
					Data: string(j), ErrorCode: 0}
			}
		}
	}
}
