package api

import (
	"controller/host"
	"encoding/json"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"os"
)

func CreateHost(w http.ResponseWriter, r *http.Request) {
	var h host.Host
	var err error

	err = json.NewDecoder(r.Body).Decode(&h)
	if err != nil {
		http.Error(w, "CreateHost(), Internal Server Error: "+err.Error(), http.StatusInternalServerError)
	}
	err = h.Post()
	if err != nil {
		http.Error(w, "CreateHost(), Internal Server Error: "+err.Error(), http.StatusInternalServerError)
	} else {
		json.NewEncoder(w).Encode(h)
	}
}

func ReadHost(w http.ResponseWriter, r *http.Request) {
	var err error
	h := new(host.Host)

	param := r.URL.Query()
	id := param["id"]
	if len(id) > 0 {
		h.Id, err = primitive.ObjectIDFromHex(id[0])
		if err != nil {
			h.Id = primitive.NilObjectID
		}
	}
	hosts, err := h.Get()
	if err != nil {
		http.Error(w, "ReadHost() - Internal Server Error: "+err.Error(), http.StatusInternalServerError)
	} else {
		err := json.NewEncoder(os.Stdout).Encode(hosts)
		if err != nil {
			log.Printf("Error encoding hosts array: %v", err)
			http.Error(w, "ReadHost() - Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		} else {
			err = json.NewEncoder(w).Encode(hosts)
		}
	}
}

func UpdateHost(w http.ResponseWriter, r *http.Request) {
	var h host.Host
	var err error

	err = json.NewDecoder(r.Body).Decode(&h)
	if err != nil {
		_ = json.NewEncoder(w).Encode(err)
	}
	err = h.Put()
	if err != nil {
		http.Error(w, "UpdateHost(), Internal Server Error: "+err.Error(), http.StatusInternalServerError)
	} else {
		json.NewEncoder(w).Encode(h)
	}
}

func DeleteHost(w http.ResponseWriter, r *http.Request) {
	var h host.Host
	var err error

	param := r.URL.Query()
	id := param["id"]
	if len(id) > 0 {
		h.Id, err = primitive.ObjectIDFromHex(id[0])
		if err != nil {
			http.Error(w, "DeleteHost(), Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		} else {
			json.NewEncoder(w).Encode(h.Delete())
		}
	}
}
