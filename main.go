package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	router := mux.NewRouter()
	go router.HandleFunc("/ws", SocketHandler)

	log.Println("Starting server")
	http.ListenAndServe(":8080", router)
}
