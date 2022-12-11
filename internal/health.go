package internal

import (
	"github.com/rs/zerolog/log"
	"net/http"
)

func RegisterHealthCheck() {
	log.Info().Msg("registering route: health check")
	http.Handle("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info().Msg("received health check request")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(`{"status": "OK"}`))
	}))
}
