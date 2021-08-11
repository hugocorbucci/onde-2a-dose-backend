package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Server represents the HTTP server
type Server struct {
	*mux.Router
}

type httpHandler struct {
}

// NewHTTPServer creates a new server
func NewHTTPServer() *Server {
	handler := &httpHandler{}

	r := mux.NewRouter()
	r.HandleFunc("/data.raw", handler.rawData).Methods(http.MethodPost)
	r.HandleFunc("/data", handler.data).Methods(http.MethodGet)

	return &Server{r}
}

func (h *httpHandler) rawData(w http.ResponseWriter, req *http.Request) {
	// TODO: Implement
}

func (h *httpHandler) data(w http.ResponseWriter, req *http.Request) {
	// TODO: Implement
}
