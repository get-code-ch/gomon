package events

import (
	"controller/command"
	"controller/history"
	"controller/host"
	"encoding/json"
	"strings"
	"time"
)

func (w Worker) Builtin(cmd command.Command) {
	var state = ""
	var message = ""
	var metric = ""

	host := host.Host{}
	hst := history.History{}
	host.Id = w.Probe.HostId
	h, err := host.Get()
	if err != nil {
		return
	}
	if len(h) < 0 {
		return
	}

	switch strings.ToUpper(cmd.Command) {
	case "WPPAGE":
		url := "http://" + h[0].FQDN
		state, message, metric = WPPage(url)
	case "ILO":
		url := "https://" + h[0].IP
		password, _ := w.Probe.GetSecret()
		state, message, metric = GetIloHealth(url, w.Probe.Username, password)
	default:

	}

	hst = history.History{}
	hst.HostId = w.Probe.HostId
	hst.ProbeId = w.Probe.Id
	hst.Timestamp = time.Now()

	hst.State = state
	hst.Message = message
	hst.Metric = metric

	// TODO: Uncomment History writing
	err = nil
	// err = hst.Post()
	if err != nil {
		Evt.msg <- time.Now().Format(time.RFC3339) + ": Error writing history for " + w.Probe.Name + " err: " + err.Error()
	}

	w.Probe.State = hst.State
	w.Probe.Result = hst.Message
	w.Probe.Last = time.Now()
	w.Probe.Next = time.Now().Add(time.Second * time.Duration(w.Probe.Interval))
	// Updating Probe
	Mux.Lock()
	w.Probe.Put()
	Mux.Unlock()

	p, _ := json.Marshal(w.Probe)
	Broadcast <- SocketMessage{Action: "UPDATE", Object: "PROBE", Data: string(p), ErrorCode: 0}

}
