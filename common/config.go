package common

import (
	"github.com/spf13/viper"
	"fmt"
)

func init(){
	initConfig()
}

func initConfig(){
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath("../")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}