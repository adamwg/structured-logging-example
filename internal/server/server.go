package server

import (
	"encoding/json"
	"math/rand"
	"net/http"

	"github.com/aybabtme/log"
	"github.com/moby/moby/pkg/namesgenerator"
)

type server struct {
	log *log.Log
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ll := s.log.KV("client_addr", r.RemoteAddr).KV("req_method", r.Method).KV("req_path", r.URL.Path)

	switch r.Method {
	case http.MethodGet:
		// OK
	default:
		errorResponse(w, http.StatusMethodNotAllowed, "invalid method", ll)
		return
	}

	luck := rand.Intn(10)
	ll = ll.KV("lucky_number", luck)
	if luck < 2 {
		errorResponse(w, http.StatusServiceUnavailable, "user is unlucky", ll)
		return
	}

	name := namesgenerator.GetRandomName(0)
	resp := struct {
		Name string `json:"name"`
	}{
		Name: name,
	}
	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		ll = ll.Err(err)
		errorResponse(w, http.StatusInternalServerError, "failed to encode response", ll)
	}
	ll.KV("random_name", name).KV("resp_code", http.StatusOK).Info("returned name to user")
}

func errorResponse(w http.ResponseWriter, code int, msg string, ll *log.Log) {
	ll.KV("resp_code", code).Error(msg)
	w.WriteHeader(code)
}

func New(l *log.Log) http.Handler {
	return &server{log: l}
}
