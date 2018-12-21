package events

import (
	"controller/config"
	"log"
	"time"
)

type Event struct {
	ticker *time.Ticker
	msg    chan string
}

var Msg chan string
var Evt Event

func init() {
	Msg = make(chan string, 10)
	Evt = Event{
		ticker: time.NewTicker(time.Millisecond * 1000 * time.Duration(config.Config.Timer)),
		msg:    make(chan string, 5),
	}
	go Evt.timer()
	go Evt.onMessage()

	go func() {
		for {
			m := <-Msg
			log.Printf("Message: %v", m)
		}
	}()
}

func (e Event) timer() {
	for t := range e.ticker.C {
		mux.Lock()
		Evt.msg <- "timer elapsed --> " + t.Format("15:04:05")
		Broadcast <- SocketMessage{Data: "timer elapsed --> " + t.Format("15:04:05"), ErrorCode: 0}
		mux.Unlock()
	}
}

func (e Event) onMessage() {
	for {
		m := <-e.msg
		log.Printf("Message: %v", m)
	}
}
