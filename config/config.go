package config

import (
	"log"

	"github.com/spf13/viper"
)

//Init initialize config
func Init(path string) *viper.Viper {
	viper.SetConfigName("config")
	viper.AddConfigPath(path)
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Fatal error config file: %v \n", err)
	}

	return viper.GetViper()
}
