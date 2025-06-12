package main

import (
	"net/http"
	"time"

	"github.com/Andrew-Wichmann/coffee-shop/pkg/logging"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins
		return true
	},
}

func websocketHandler(rw http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(rw, req, nil)
	if err != nil {
		logging.Logger.Error("Could not upgrade the connection to a websocket", zap.Error(err))
		return
	}
	defer conn.Close()
	logging.Logger.Debug("Connection established")

	for {
		_, _message, err := conn.ReadMessage()
		if err != nil {
			logging.Logger.Error("Error reading message", zap.Error(err))
			return
		}
		message := string(_message)
		logging.Logger.Debug("Received message", zap.String("message_received", message))
		var response string
		switch message {
		case "enter":
			response = "Welcome"
		case "leave":
			response = "Thank you for visiting!"
		default:
			response = "Unknown action."
		}

		logging.Logger.Debug("Sending response", zap.String("response sent", response))
		err = conn.WriteMessage(websocket.TextMessage, []byte(response))
		if err != nil {
			logging.Logger.Error("Could not send 'world' response", zap.Error(err))
			return
		}

		if message == "leave" {
			err = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Goodbye"))
			if err != nil {
				logging.Logger.Error("Could not send the close websocket message")
				return
			}
			logging.Logger.Info("Connection closed gracefully")
			return
		}
	}
}

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
	mux.HandleFunc("/ws", websocketHandler)
	logging.Logger.Info("Starting server!")
	http.ListenAndServe(":8080", loggingMiddleware(mux))
}
