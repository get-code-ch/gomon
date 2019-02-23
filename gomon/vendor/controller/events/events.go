package events

import (
	"controller/command"
	"controller/config"
	"controller/probe"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"strings"
	"time"
)

type Event struct {
	ticker *time.Ticker
	msg    chan string
}

type Worker struct {
	Ticker *time.Ticker
	Probe  probe.Probe
}

var Msg chan string
var Evt Event
var Workers map[primitive.ObjectID]Worker

func init() {
	Msg = make(chan string, 2)

	Evt = Event{
		ticker: time.NewTicker(time.Millisecond * 1000 * time.Duration(config.Config.Timer)),
		msg:    make(chan string, 2),
	}
	go Evt.timer()
	go Evt.onMessage()

	go func() {
		for {
			m := <-Msg
			log.Printf("Message: %v", m)
		}
	}()

	Workers = make(map[primitive.ObjectID]Worker)
	p := probe.Probe{}
	probes, err := p.Get()
	if err != nil {
		log.Printf("Error getting probes: %v", err)
		return
	}

	for _, p := range probes {
		if p.Interval >= 1 {
			w := Worker{
				Ticker: time.NewTicker(time.Millisecond * 1000 * time.Duration(p.Interval)),
				Probe:  p,
			}
			Workers[p.Id] = w
			go Workers[p.Id].Timer()
		}
	}

}

func (w Worker) Stop() {
	log.Printf("Stop worker %v", w)
}

func (w Worker) Timer() {
	for t := range w.Ticker.C {

		// Reading and running command
		cPtr := command.Command{}
		cPtr.Id = w.Probe.CommandId
		c, err := cPtr.Get()
		if err != nil {
			Evt.msg <- "Error getting command for Probe --> " + w.Probe.Name + " error: " + err.Error()
			Broadcast <- SocketMessage{Action: "LOG", Data: t.Format(time.RFC3339) + ": Error getting command for Probe --> " + w.Probe.Name + " error: " + err.Error(), ErrorCode: 0}
			continue
		}

		if len(c) < 1 {
			Evt.msg <- "Error getting command for Probe --> " + w.Probe.Name + " error: Command not Found"
			Broadcast <- SocketMessage{Action: "LOG", Data: t.Format(time.RFC3339) + ": Error getting command for Probe --> " + w.Probe.Name + " error: Command not Found", ErrorCode: 0}
			continue
		}

		// TODO: For test purpose only, in this place we must handling command
		switch strings.ToUpper(c[0].CommandType) {
		case "BASH":
			go w.Bash(c[0])
		case "BUILTIN":
			go w.Builtin(c[0])
		default:
			Evt.msg <- "WRONG commandType for Probe --> " + w.Probe.Name + " Command --> " + c[0].Key
			Broadcast <- SocketMessage{Action: "LOG", Data: t.Format(time.RFC3339) + ": WRONG commandType for Probe --> " + w.Probe.Name + " Command --> " + c[0].Key, ErrorCode: 0}
		}

	}
}

func (e Event) timer() {
}

func (e Event) onMessage() {
	for {
		m := <-e.msg
		log.Printf("Message: %v", m)
	}
}
