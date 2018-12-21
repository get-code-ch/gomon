package api

import (
	"encoding/json"
	"model/menu"
	"net/http"
)

func GetMenu(w http.ResponseWriter, r *http.Request) {
	m, _ := menu.ReadMenu()
	json.NewEncoder(w).Encode(m)
}
