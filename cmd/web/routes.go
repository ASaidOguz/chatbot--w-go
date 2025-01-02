package main

import (
	"net/http"
	"websocket/internal/handlers"

	"github.com/bmizerany/pat"
)

// routes is a function that returns a http.Handler.
func routes() http.Handler {
	mux := pat.New()

	mux.Get("/", http.HandlerFunc(handlers.Home))
	mux.Get("/ws", http.HandlerFunc(handlers.WsEndPoint))

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return mux
}
