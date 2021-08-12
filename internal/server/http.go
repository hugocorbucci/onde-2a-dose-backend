package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/hugocorbucci/onde-2a-dose-backend/internal/clients/prefeitura"
	deps "github.com/hugocorbucci/onde-2a-dose-backend/internal/dependencies"
)

const (
	mandatoryBodyFieldName = "dados"

	JSONContentType = "application/json; charset=UTF-8"
)

// Server represents the HTTP server
type Server struct {
	*mux.Router
}

type httpHandler struct {
	DeOlhoNaFilaClient deps.DeOlhoNaFila
}

// NewHTTPServer creates a new server
func NewHTTPServer(client deps.DeOlhoNaFila) *Server {
	handler := &httpHandler{DeOlhoNaFilaClient: client}

	r := mux.NewRouter()
	r.HandleFunc("/data.raw", handler.rawData).Methods(http.MethodPost)
	r.HandleFunc("/data", handler.data).Methods(http.MethodGet)

	return &Server{r}
}

func (h *httpHandler) rawData(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	if len(req.PostForm) == 0 {
		h.writeError(w, http.StatusBadRequest, "missing body", nil)
		return
	}
	val := req.PostForm.Get(mandatoryBodyFieldName)
	if len(val) == 0 {
		h.writeError(w, http.StatusBadRequest, "invalid body", nil)
		return
	}

	units, err := h.DeOlhoNaFilaClient.Fetch(req.Context())
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "error fetching data", err)
		return
	}

	w.Header().Add(prefeitura.ContentTypeHeader, JSONContentType)
	err = json.NewEncoder(w).Encode(units)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "error encoding data", err)
		return
	}
}

func (h *httpHandler) data(w http.ResponseWriter, _ *http.Request) {
	// TODO: Implement
	h.writeError(w, http.StatusNotImplemented, "not implemented yet", nil)
}

func (h *httpHandler) writeError(w http.ResponseWriter, statusCode int, baseMessage string, err error) {
	w.WriteHeader(statusCode)
	w.Header().Add(prefeitura.ContentTypeHeader, JSONContentType)

	errorMsg := baseMessage
	if err != nil {
		errorMsg = fmt.Sprintf("%s: %s", baseMessage, err)
	}
	w.Write([]byte(fmt.Sprintf("{\"error\":\"%s\"}", errorMsg)))
}