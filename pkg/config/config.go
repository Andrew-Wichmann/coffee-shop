package config

import (
	"flag"

	"github.com/spf13/viper"
)

func init() {
	configFile := flag.String("config", "config.yaml", "Path to the configuration file.")
	flag.Parse()
	viper.SetConfigFile(*configFile)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
