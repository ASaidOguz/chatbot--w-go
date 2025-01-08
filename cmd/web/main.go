package main

import (
	"log"
	"net/http"
	"websocket/internal/handlers"
)

func main() {
	routes := routes()

	log.Println("Starting Channel Listener")
	// Start listening to the channel.
	go handlers.ListenToWsChannel()

	log.Println("Starting server on :8080")
	// Start the server.
	_ = http.ListenAndServe(":8080", routes)
}
