package events

import (
	"bytes"
	"controller/command"
	"controller/history"
	"controller/host"
	"encoding/json"
	"log"
	"os/exec"
	"regexp"
	"time"
)

func (w Worker) Bash(cmd command.Command) {
	var r *regexp.Regexp

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

	// TODO: Parse and convert parameters / move to a function
	r = regexp.MustCompile(`(?mi)~ip~`)
	c := r.ReplaceAllString(cmd.Command, h[0].IP)
	r = regexp.MustCompile(`(?mi)~hostname~`)
	c = r.ReplaceAllString(c, h[0].Name)
	r = regexp.MustCompile(`(?mi)~username~`)
	c = r.ReplaceAllString(c, w.Probe.Username)
	r = regexp.MustCompile(`(?mi)~password~`)
	pwd, _ := w.Probe.GetSecret()
	c = r.ReplaceAllString(c, pwd)
	w.Probe.Password = ""

	e := exec.Command("sh", "-c", c)
	cmdOutput := &bytes.Buffer{}
	e.Stdout = cmdOutput
	err = e.Run()
	if err != nil {
		Evt.msg <- time.Now().Format(time.RFC3339) + ": Bash error for --> " + w.Probe.Name + " err: " + err.Error() + " Out : " + string(cmdOutput.Bytes())
		//Broadcast <- SocketMessage{Data: "Bash error for --> " + w.Probe.Name + " err: " + err.Error(), ErrorCode: 0}
		//return
	}

	// TODO: Parse bash script response
	hst = history.History{}
	hst.HostId = w.Probe.HostId
	hst.ProbeId = w.Probe.Id
	hst.Timestamp = time.Now()
	r = regexp.MustCompile(`(?Um)^(.*):\s*(.*)(?:\s*\|\s*(.*))?$`)
	match := r.FindSubmatch(cmdOutput.Bytes())
	log.Printf("%v", string(cmdOutput.Bytes()))
	if match != nil && len(match) == 4 {
		hst.State = string(match[1])
		hst.Message = string(match[2])
		hst.Metric = string(match[3])
	}

	// TODO: Uncomment History writing
	err = nil
	// err = hst.Post()
	if err != nil {
		Evt.msg <- time.Now().Format(time.RFC3339) + ": Error writing history for " + w.Probe.Name + " err: " + err.Error()
	}

	w.Probe.State = hst.State
	w.Probe.Result = hst.Message
	if hst.Metric != "" {
		w.Probe.Result += " | " + hst.Metric
	}
	w.Probe.Last = time.Now()
	w.Probe.Next = time.Now().Add(time.Second * time.Duration(w.Probe.Interval))
	// Updating Probe
	Mux.Lock()
	w.Probe.Put()
	Mux.Unlock()

	p, _ := json.Marshal(w.Probe)
	Broadcast <- SocketMessage{Action: "UPDATE", Object: "PROBE", Data: string(p), ErrorCode: 0}
}
