package model

type BlockChainConfiguration struct {
	GenesisConfPath string
}

func NewBlockChainConfiguration() BlockChainConfiguration {
	return BlockChainConfiguration{
		GenesisConfPath: "./gensis.conf",
	}
}
