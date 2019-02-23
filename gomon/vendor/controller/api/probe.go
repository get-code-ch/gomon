package api

import (
	"controller/probe"
	"encoding/json"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"os"
)

func CreateProbe(w http.ResponseWriter, r *http.Request) {
	var p probe.Probe
	var err error

	err = json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, "CreateProbe(), Internal Server Error: "+err.Error(), http.StatusInternalServerError)
	}
	err = p.Post()
	if err != nil {
		http.Error(w, "CreateProbe(), Internal Server Error: "+err.Error(), http.StatusInternalServerError)
	} else {
		json.NewEncoder(w).Encode(p)
	}
}

func ReadProbe(w http.ResponseWriter, r *http.Request) {
	var err error
	p := new(probe.Probe)

	param := r.URL.Query()
	id := param["id"]
	if len(id) > 0 {
		p.Id, err = primitive.ObjectIDFromHex(id[0])
		if err != nil {
			p.Id = primitive.NilObjectID
		}
	}
	probes, err := p.Get()
	if err != nil {
		http.Error(w, "ReadProbe() - Internal Server Error: "+err.Error(), http.StatusInternalServerError)
	} else {
		err := json.NewEncoder(os.Stdout).Encode(probes)
		if err != nil {
			log.Printf("Error encoding probes array: %v", err)
			http.Error(w, "ReadProbe() - Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		} else {
			err = json.NewEncoder(w).Encode(probes)
		}
	}
}

func UpdateProbe(w http.ResponseWriter, r *http.Request) {
	var p probe.Probe
	var err error

	err = json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		_ = json.NewEncoder(w).Encode(err)
	}
	err = p.Put()
	if err != nil {
		http.Error(w, "UpdateProbe(), Internal Server Error: "+err.Error(), http.StatusInternalServerError)
	} else {
		json.NewEncoder(w).Encode(p)
	}
}

func DeleteProbe(w http.ResponseWriter, r *http.Request) {
	var p probe.Probe
	var err error

	param := r.URL.Query()
	id := param["id"]
	if len(id) > 0 {
		p.Id, err = primitive.ObjectIDFromHex(id[0])
		if err != nil {
			http.Error(w, "DeleteProbe(), Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		} else {
			json.NewEncoder(w).Encode(p.Delete())
		}
	}
}
