package model

type BlockChainConfiguration struct {
	RepositoryPath string
}

func NewBlockChainConfiguration() BlockChainConfiguration {
	return BlockChainConfiguration{
		RepositoryPath: "empty",
	}
}
