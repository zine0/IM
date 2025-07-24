package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func SetConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.config/im")
	viper.AddConfigPath("/etc/im")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %v",err))
	}

}
