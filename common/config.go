package common

import   (
	"github.com/spf13/viper"
	"fmt"
	"os"
)

func init(){
	initConfig()
}

func initConfig(){
	pathSeparator := string(os.PathSeparator)

	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath("$GOPATH" + pathSeparator + "src" + pathSeparator + "it-chain")

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}