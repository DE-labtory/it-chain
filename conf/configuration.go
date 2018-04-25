package conf

import (
	"os"

	"github.com/it-chain/it-chain-Engine/conf/model"
	"github.com/it-chain/it-chain-Engine/conf/model/common"
	"github.com/spf13/viper"
)

type Configuration struct {
	Common         common.CommonConfiguration
	Txpool         model.TxpoolConfiguration
	Consensus      model.ConsensusConfiguration
	Blockchain     model.BlockChainConfiguration
	Peer           model.PeerConfiguration
	Authentication model.AuthenticationConfiguration
	Icode          model.ICodeConfiguration
}

func GetConfiguration() (*Configuration, error) {
	path, _ := os.Getwd()
	viper.SetConfigName("config")
	viper.AddConfigPath(path + "/conf")
	var configuration = Configuration{}
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	err := viper.Unmarshal(&configuration)
	if err != nil {
		return nil, err
	}
	return &configuration, nil
}
