package main

import (
	"controller"
	"controller/authorize"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"view"
)

func main() {
	var err error

	router := mux.NewRouter()

	// main page
	router.HandleFunc("/", view.Root)
	router.HandleFunc("/login", authorize.Login)
	router.HandleFunc("/logout", authorize.Logout)
	router.HandleFunc("/probes", view.Probes)

	// Serving static files
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(controller.Config.StaticFolder))))
	router.HandleFunc("/favicon.ico", faviconHandler)

	// Starting server
	log.Printf("Starting server %s:%s...", controller.Config.Server, controller.Config.Port)
	if controller.Config.Ssl {
		err = http.ListenAndServeTLS(controller.Config.Server+":"+controller.Config.Port, controller.Config.Cert, controller.Config.Key, router)
	} else {
		err = http.ListenAndServe(controller.Config.Server+":"+controller.Config.Port, router)
	}
	if err != nil {
		log.Fatal("ListenAndServer:", err)
	}

}

func faviconHandler(w http.ResponseWriter, r *http.Request) {

	http.ServeFile(w, r, "./static/favicon.ico")

}
