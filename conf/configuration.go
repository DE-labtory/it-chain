package conf

import (
	"os"
	"sync"

	"github.com/it-chain/it-chain-Engine/conf/model"
	"github.com/it-chain/it-chain-Engine/conf/model/common"
	"github.com/spf13/viper"
)

var instance *Configuration
var once sync.Once

type Configuration struct {
	Common         common.CommonConfiguration
	Txpool         model.TxpoolConfiguration
	Consensus      model.ConsensusConfiguration
	Blockchain     model.BlockChainConfiguration
	Peer           model.PeerConfiguration
	Authentication model.AuthenticationConfiguration
	Icode          model.ICodeConfiguration
	GrpcGateway    model.GrpcGatewayConfiguration
}

func GetConfiguration() *Configuration {
	once.Do(func() {
		instance = &Configuration{}
		path, _ := os.Getwd()
		viper.SetConfigName("config")
		viper.AddConfigPath(path)
		if err := viper.ReadInConfig(); err != nil {
			panic("cannot read config")
		}
		err := viper.Unmarshal(&instance)
		if err != nil {
			panic("error in read config")
		}
	})
	return instance
}
