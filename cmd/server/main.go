package main

import (
	"log"
	"net"
	"net/http"
	"os"

	"github.com/hugocorbucci/onde-2a-dose-backend/internal/server"
)

const (
	defaultPort = "8080"
)

func main() {
	ll := log.New(os.Stdout, "Onde2aDose - ", 0)
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = defaultPort
	}
	addr := net.JoinHostPort("", port)

	ll.Println("Starting server on port", port)
	s := server.NewHTTPServer()
	if err := http.ListenAndServe(addr, s); err != nil {
		ll.Fatal("HTTP(s) server failed")
	}
}
