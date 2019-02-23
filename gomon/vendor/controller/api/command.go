package api

import (
	"controller/command"
	"encoding/json"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"os"
)

func CreateCommand(w http.ResponseWriter, r *http.Request) {
	var c command.Command
	var err error

	err = json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		http.Error(w, "CreateCommand(), Internal Server Error: "+err.Error(), http.StatusInternalServerError)
	}
	err = c.Post()
	if err != nil {
		http.Error(w, "CreateCommand(), Internal Server Error: "+err.Error(), http.StatusInternalServerError)
	} else {
		json.NewEncoder(w).Encode(c)
	}
}

func ReadCommand(w http.ResponseWriter, r *http.Request) {
	var err error
	c := new(command.Command)

	param := r.URL.Query()
	id := param["id"]
	if len(id) > 0 {
		c.Id, err = primitive.ObjectIDFromHex(id[0])
		if err != nil {
			c.Id = primitive.NilObjectID
		}
	}
	commands, err := c.Get()
	if err != nil {
		http.Error(w, "ReadCommand() - Internal Server Error: "+err.Error(), http.StatusInternalServerError)
	} else {
		err := json.NewEncoder(os.Stdout).Encode(commands)
		if err != nil {
			log.Printf("Error encoding commands array: %v", err)
			http.Error(w, "ReadCommand() - Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		} else {
			err = json.NewEncoder(w).Encode(commands)
		}
	}
}

func UpdateCommand(w http.ResponseWriter, r *http.Request) {
	var c command.Command
	var err error

	err = json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		_ = json.NewEncoder(w).Encode(err)
	}
	err = c.Put()
	if err != nil {
		http.Error(w, "UpdateCommand(), Internal Server Error: "+err.Error(), http.StatusInternalServerError)
	} else {
		json.NewEncoder(w).Encode(c)
	}
}

func DeleteCommand(w http.ResponseWriter, r *http.Request) {
	var c command.Command
	var err error

	param := r.URL.Query()
	id := param["id"]
	if len(id) > 0 {
		c.Id, err = primitive.ObjectIDFromHex(id[0])
		if err != nil {
			http.Error(w, "DeleteCommand(), Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		} else {
			json.NewEncoder(w).Encode(c.Delete())
		}
	}
}
