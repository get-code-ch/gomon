package main

import (
	"controller/admin"
	"controller/api"
	"controller/authorize"
	"controller/config"
	"controller/events"
	"controller/host"
	"controller/probes"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	var err error

	router := mux.NewRouter()

	//handlers.AllowedOrigins([]string{"*"})
	// login management
	router.HandleFunc("/login", authorize.Login).Methods(http.MethodPost)
	router.HandleFunc("/logout", authorize.Logout)

	// Probes management
	router.HandleFunc("/probes", authorize.ValidateMiddlewareCookie(probes.ListProbes)).Methods(http.MethodGet)

	// Hosts management
	router.HandleFunc("/hosts", authorize.ValidateMiddlewareCookie(host.ListHosts)).Methods(http.MethodGet)
	router.HandleFunc("/hosts", authorize.ValidateMiddlewareCookie(host.CreateHosts)).Methods(http.MethodPost)

	// Admin pages
	router.HandleFunc("/admin", admin.Admin).Methods(http.MethodGet)

	// Serving static files
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(config.Config.StaticFolder))))
	router.HandleFunc("/favicon.ico", faviconHandler)
	router.PathPrefix("/wasm/").Handler(http.StripPrefix("/wasm/", http.FileServer(http.Dir("/home/claude/go/src/wasm/."))))

	// API REST services
	router.HandleFunc("/api/authenticate", authorize.CreateTokenEndpoint).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/api/protected", authorize.ValidateMiddlewareToken(api.Api)).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/menu", authorize.ValidateMiddlewareToken(api.GetMenu)).Methods(http.MethodGet, http.MethodOptions)

	// Host CRUD
	router.HandleFunc("/api/host", authorize.ValidateMiddlewareToken(api.CreateHost)).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/api/host", authorize.ValidateMiddlewareToken(api.ReadHost)).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/host", authorize.ValidateMiddlewareToken(api.UpdateHost)).Methods(http.MethodPut, http.MethodOptions)
	router.HandleFunc("/api/host", authorize.ValidateMiddlewareToken(api.DeleteHost)).Methods(http.MethodDelete, http.MethodOptions)

	// websocket services
	router.HandleFunc("/ws", authorize.ValidateMiddlewareCookie(events.HandleConnections)).Methods(http.MethodOptions)
	router.HandleFunc("/socket", events.Upgrader)

	// main page
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(config.Config.StaticFolder+"/client/"))))
	//router.HandleFunc("/", root.Root)

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

func faviconHandler(w http.ResponseWriter, r *http.Request) {

	http.ServeFile(w, r, "./static/favicon.ico")

}
