package main

import (
	"github.com/Andrew-Wichmann/coffee-shop/pkg/logging"
)

func main() {
	logging.Logger.Debug("Hello world!")
	logging.Logger.Info("Hello world!")
	logging.Logger.Warn("Hello world!")
	logging.Logger.Error("Hello world!")
}
