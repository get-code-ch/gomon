package main

import (
	"controller/api"
	"controller/authorize"
	"controller/config"
	"controller/events"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	var err error

	router := mux.NewRouter().StrictSlash(true)

	// Serving static files
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(config.Config.StaticFolder))))

	// API REST services
	router.HandleFunc("/api/authenticate", authorize.CreateTokenEndpoint).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/api/menu", authorize.ValidateMiddlewareToken(api.GetMenu)).Methods(http.MethodGet, http.MethodOptions)

	// Host CRUD
	router.HandleFunc("/api/host", authorize.ValidateMiddlewareToken(api.CreateHost)).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/api/host", authorize.ValidateMiddlewareToken(api.ReadHost)).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/host", authorize.ValidateMiddlewareToken(api.UpdateHost)).Methods(http.MethodPut, http.MethodOptions)
	router.HandleFunc("/api/host", authorize.ValidateMiddlewareToken(api.DeleteHost)).Methods(http.MethodDelete, http.MethodOptions)

	// Probe CRUD
	router.HandleFunc("/api/probe", authorize.ValidateMiddlewareToken(api.CreateProbe)).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/api/probe", authorize.ValidateMiddlewareToken(api.ReadProbe)).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/probe", authorize.ValidateMiddlewareToken(api.UpdateProbe)).Methods(http.MethodPut, http.MethodOptions)
	router.HandleFunc("/api/probe", authorize.ValidateMiddlewareToken(api.DeleteProbe)).Methods(http.MethodDelete, http.MethodOptions)

	// Command CRUD
	router.HandleFunc("/api/command", authorize.ValidateMiddlewareToken(api.CreateCommand)).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/api/command", authorize.ValidateMiddlewareToken(api.ReadCommand)).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/command", authorize.ValidateMiddlewareToken(api.UpdateCommand)).Methods(http.MethodPut, http.MethodOptions)
	router.HandleFunc("/api/command", authorize.ValidateMiddlewareToken(api.DeleteCommand)).Methods(http.MethodDelete, http.MethodOptions)

	// websocket services
	router.HandleFunc("/socket", events.Upgrader)

	// Catch routes Handle by client app
	router.HandleFunc("/hosts/", catchAllHandler)
	router.HandleFunc("/host/", catchAllHandler)
	router.HandleFunc("/probe/", catchAllHandler)
	router.HandleFunc("/command/", catchAllHandler)
	router.HandleFunc("/liveview/", catchAllHandler)
	router.HandleFunc("/logout/", catchAllHandler)

	// Serving client app
	router.PathPrefix("/").Handler(http.FileServer(http.Dir(config.Config.AppFolder)))

	// Starting server
	log.Printf("Starting server %s:%s...", config.Config.Server, config.Config.Port)
	if config.Config.Ssl {
		err = http.ListenAndServeTLS(config.Config.Server+":"+config.Config.Port, config.Config.Cert, config.Config.Key, handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", "X-Token"}), handlers.AllowedMethods([]string{"GET", "DELETE", "POST", "PUT", "HEAD", "OPTIONS", "CONNECT"}), handlers.AllowedOrigins([]string{"*"}))(router))
	} else {
		err = http.ListenAndServe(config.Config.Server+":"+config.Config.Port, handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", "X-Token"}), handlers.AllowedMethods([]string{"GET", "DELETE", "POST", "PUT", "HEAD", "OPTIONS", "CONNECT"}), handlers.AllowedOrigins([]string{"*"}))(router))
	}
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

func catchAllHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusPermanentRedirect)
}
