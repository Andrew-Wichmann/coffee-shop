package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Andrew-Wichmann/coffee-shop/pkg/logging"
	"github.com/Andrew-Wichmann/coffee-shop/pkg/models"
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

type customer struct {
	coffee models.Coffee
	ticket models.Ticket
}

func (c *customer) orderCoffee() error {
	coffee, err := models.CoffeeStore.OrderCoffee(models.HOUSE)
	if err != nil {
		return err
	}
	c.coffee = coffee
	return nil
}

func (c *customer) drinkCoffee() (string, error) {
	if c.coffee == nil {
		return "You have no coffee to drink!", nil
	}
	response := fmt.Sprintf("You enjoy your %s", c.coffee.Name())
	c.coffee = nil
	return response, nil
}

func websocketHandler(rw http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(rw, req, nil)
	if err != nil {
		logging.Logger.Error("Could not upgrade the connection to a websocket", zap.Error(err))
		return
	}
	defer conn.Close()
	logging.Logger.Debug("Connection established")
	var customer customer
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
			response = "Welcome. There are X people in line. Would you like to stand in line?"
		case "take ticket":
			var ticket models.Ticket
			ticket, err = models.CoffeeStore.TakeTicket(customer.orderCoffee)
			customer.ticket = ticket
			response = "Ticket taken"
		case "drink coffee":
			response, err = customer.drinkCoffee()
		case "check ticket":
			response = "checking ticket"
		case "order":
			response = "You ordered a coffee."
		case "sit":
			response = "You are sitting."
		case "leave":
			response = "Thank you for visiting!"
		default:
			response = "Unknown action."
		}

		if err != nil {
			logging.Logger.Info("Something went wrong with this customer. We're kicking them out and closing the door/connection", zap.Error(err))
			err = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Goodbye"))
			if err != nil {
				logging.Logger.Error("Could not send the close websocket message")
				return
			}
			logging.Logger.Info("Connection closed gracefully")
			return
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
	logging.Logger.Info("Opening the store!")
	models.CoffeeStore.Open()
	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/ws", websocketHandler)
	logging.Logger.Info("Starting server!")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	server := http.Server{
		Addr:    ":8080",
		Handler: loggingMiddleware(mux),
	}
	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			logging.Logger.Panic("Server exited unexpectedly", zap.Error(err))
		}
	}()
	<-stop

	logging.Logger.Info("Shutting down the server...")

	// According to Chat-Gipity, typical Kubernetes deployments five 30 seconds
	// for the server to shut down.
	var wg sync.WaitGroup
	wg.Add(2)
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()
	go func() {
		defer wg.Done()

		if err := server.Shutdown(ctx); err != nil {
			logging.Logger.Error("Server forced to shutdown", zap.Error(err))
		} else {
			logging.Logger.Info("Server shutdown gracefully")
		}
	}()

	go func() {
		defer wg.Done()
		models.CoffeeStore.Close(ctx)
	}()

	wg.Wait()
	logging.Logger.Info("Everything shutdown. Exiting.")
}
