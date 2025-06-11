package logging

import (
	"encoding/json"

	_ "github.com/Andrew-Wichmann/coffee-shop/pkg/config"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var Logger *zap.Logger

func init() {
	loggingConfig := viper.GetStringMap("logging")
	jsonData, err := json.Marshal(loggingConfig)
	var cfg zap.Config
	if err := json.Unmarshal(jsonData, &cfg); err != nil {
		panic(err)
	}
	Logger, err = cfg.Build()
	if err != nil {
		panic(err)
	}
}
