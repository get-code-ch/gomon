package events

import (
	"controller"
	"log"
	"time"
)

var EventsMsg chan string

func init() {
	EventsMsg = make(chan string, 10)
	ticker := time.NewTicker(time.Millisecond * 1000 * time.Duration(controller.Config.Timer))
	go func() {
		for t := range ticker.C {
			EventsMsg <- "timer elapsed --> " + t.Format("15:04:05")
		}
	}()

	go func() {
		for {
			item, ok := <-EventsMsg
			if !ok {
				log.Printf("Error eventsMsg...")
				break
			}
			log.Printf("eventsMsg: %v", item)
		}
	}()
}
