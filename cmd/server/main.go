package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/hugocorbucci/onde-2a-dose-backend/internal/clients/prefeitura"
	"github.com/hugocorbucci/onde-2a-dose-backend/internal/server"
)

const (
	defaultPort = "8080"
	fetchTimeout = 3 * time.Second
)

func main() {
	ll := log.New(os.Stdout, "Onde2aDose - ", 0)
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = defaultPort
	}
	addr := net.JoinHostPort("", port)

	httpClient := http.DefaultClient
	httpClient.Timeout = fetchTimeout
	httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	prefeituraClient := &prefeitura.Client{HTTPClient: httpClient}

	ll.Println("Starting server on port", port)
	s := server.NewHTTPServer(prefeituraClient)
	if err := http.ListenAndServe(addr, s); err != nil {
		ll.Fatal("HTTP(s) server failed")
	}
}
