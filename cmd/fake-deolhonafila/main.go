package main

import (
	_ "embed"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

const (
	defaultPort = "8082"
)

var (
	//go:embed stub.json
	deolhonafilaStub []byte
)

type fakeDeOlhoNaFilaServer struct {
	*mux.Router
}

type httpHandler struct {
	ll *log.Logger
}

func (h *httpHandler) log(args ...interface{}) {
	h.ll.Println(args...)
}

func (h *httpHandler) dados(w http.ResponseWriter, _ *http.Request) {
	h.log("request dados")
	w.Write(deolhonafilaStub)
}

func main() {
	ll := log.New(os.Stdout, "Fake DeOlhoNaFila - ", 0)

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = defaultPort
	}
	addr := net.JoinHostPort("", port)

	s := newHTTPServer(ll)
	ll.Println("Starting on port", port)
	if err := http.ListenAndServe(addr, s); err != nil {
		ll.Fatal("HTTP(s) server failed")
	}
}

func newHTTPServer(ll *log.Logger) *fakeDeOlhoNaFilaServer {
	handler := &httpHandler{ll}

	r := mux.NewRouter()
	r.HandleFunc("/processadores/dados.php", handler.dados).Methods(http.MethodPost)

	return &fakeDeOlhoNaFilaServer{r}
}