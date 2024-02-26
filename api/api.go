package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
)

func makeHTTPHandlerFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
		}
	}
}

func setCORS(w http.ResponseWriter, s string) {
	w.Header().Set("Access-Control-Allow-Origin", s)
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func (s *APIServer) Run() {
	http.Handle("/vault/", makeHTTPHandlerFunc(s.handleVault))
	log.Info().Msgf("server listening on port %v", s.srv.Addr)

	go func() {
		if err := s.srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Msgf("server closed: %v", err)
		}
		log.Debug().Msg("halting service...")
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := s.srv.Shutdown(shutdownCtx); err != nil {
		log.Fatal().Msgf("server err: %v", err)
	}
	log.Info().Msgf("server on port %v closed", s.srv.Addr)
}

func NewAPIServer(listenerAddress string, CORS string) *APIServer {
	return &APIServer{
		CORS: CORS,
		srv: &http.Server{
			Addr: listenerAddress,
		},
	}
}
