package api

import (
	"encoding/json"
	"log"
	"net/http"
)

func makeHTTPHandlerFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
		}
	}
}

func enableCors(w http.ResponseWriter, s string) {
	w.Header().Set("Access-Control-Allow-Origin", s)
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func (s *APIServer) Run() {
	http.Handle("/vault", makeHTTPHandlerFunc(s.handleVault))
	log.Println("Running server on port:", s.ListenerAddress)
	http.ListenAndServe(":"+s.ListenerAddress, nil)
}
