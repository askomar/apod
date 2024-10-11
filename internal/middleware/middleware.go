package middleware

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/askomar/apod/config"
	"github.com/askomar/apod/internal/image-metadata/dtos"
	log "github.com/sirupsen/logrus"
)

func NewServer(cfg config.Config) *http.Server {
	return &http.Server{
		Addr:         ":8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
}

func AccessLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			log.Infof("Method: %v Path: %v Remote address: %v UserAgent: %v", r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
		}()
		next.ServeHTTP(w, r)
	})
}

func PanicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Error("recovered", err)
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(dtos.ErrorResponse{
					Error: "internal_error",
				})
			}
		}()
		next.ServeHTTP(w, r)
	})
}
