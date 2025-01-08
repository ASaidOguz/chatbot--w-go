package main

import (
	"net/http"
	"websocket/internal/handlers"

	"github.com/bmizerany/pat"
)

// routes is a function that returns a http.Handler.
func routes() http.Handler {
	// Create a new muxer using the pat package.
	mux := pat.New()

	// Register the home page handler.
	mux.Get("/", http.HandlerFunc(handlers.Home))
	// Register the ws endpoint handler.
	mux.Get("/ws", http.HandlerFunc(handlers.WsEndPoint))

	// Register the static file server.
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return mux
}
