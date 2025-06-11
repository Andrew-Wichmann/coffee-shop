package main

import (
	"net/http"
	"time"

	"github.com/Andrew-Wichmann/coffee-shop/pkg/logging"
	"go.uber.org/zap"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		requestLogger := logging.Logger.With(
			zap.String("method", req.Method), zap.String("requestURI", req.RequestURI),
		)
		requestLogger.Info("Request received")
		start := time.Now()
		next.ServeHTTP(w, req)
		duration := time.Since(start)
		requestLogger.Info("Request finished", zap.Duration("duration", duration))
	})
}

func rootHandler(rw http.ResponseWriter, req *http.Request) {
	rw.WriteHeader(200)
	rw.Write([]byte("Hello world"))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)
	logging.Logger.Info("Starting server!")
	http.ListenAndServe(":8080", loggingMiddleware(mux))
}
