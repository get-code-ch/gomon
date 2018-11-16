package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	var err error

	router := mux.NewRouter()

	// main page
	router.HandleFunc("/", root)
	router.HandleFunc("/login", login)
	router.HandleFunc("/logout", logout)

	// Serving static files
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(config.StaticFolder))))

	// Starting server
	log.Printf("Starting server %s:%s...", config.Server, config.Port)
	if config.Ssl {
		err = http.ListenAndServeTLS(config.Server+":"+config.Port, config.Cert, config.Key, router)
	} else {
		err = http.ListenAndServe(config.Server+":"+config.Port, router)
	}
	if err != nil {
		log.Fatal("ListenAndServer:", err)
	}

}
