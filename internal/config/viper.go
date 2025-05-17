package config

import (
	"log"

	"github.com/spf13/viper"
)

func LoadConfig(configPath string, configName string, out interface{}) {
	viper.SetConfigFile(configPath)
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	err = viper.Unmarshal(out)
	if err != nil {
		log.Fatalf("error unmarshalling config: %v", err)
	}
}
