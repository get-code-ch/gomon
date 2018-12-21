package api

import (
	"encoding/json"
	"net/http"
)

// Handle api sample - if user is authenticate return username
func Api(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("{'Message': 'gomon RESTApi Service'}")
}
