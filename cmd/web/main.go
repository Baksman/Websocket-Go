package main

import (
	"log"
	"net/http"

	"websocket/internal/handlers"
)

func main() {
	mux := routes()

	log.Println("starting channel listener")

	go handlers.ListenToWsChannel()

	_ = http.ListenAndServe(":8080", mux)
}
