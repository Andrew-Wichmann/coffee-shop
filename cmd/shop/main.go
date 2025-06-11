package main

import (
	"github.com/Andrew-Wichmann/coffee-shop/pkg/logging"
)

func main() {
	logger := logging.GetLogger()
	logger.Debug("Hello world!")
	logger.Info("Hello world!")
	logger.Warn("Hello world!")
	logger.Error("Hello world!")
}
