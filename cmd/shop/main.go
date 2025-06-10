package main

import (
	"encoding/json"
	"flag"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	configFile := flag.String("config", "config.yaml", "Path to the configuration file.")
	flag.Parse()
	viper.SetConfigFile(*configFile)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	loggingConfig := viper.GetStringMap("logging")
	jsonData, err := json.Marshal(loggingConfig)
	var cfg zap.Config
	if err := json.Unmarshal(jsonData, &cfg); err != nil {
		panic(err)
	}
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	logger.Info("Hello world!")
}
