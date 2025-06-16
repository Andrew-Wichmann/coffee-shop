package config

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	configFile := pflag.String("config", "config.yaml", "Path to the configuration file.")
	pflag.Parse()
	viper.SetConfigFile(*configFile)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
