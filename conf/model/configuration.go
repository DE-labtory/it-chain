package model

import "github.com/it-chain/it-chain-Engine/conf/model/common"

type Configuration struct {
	Common         common.CommonConfiguration
	Txpool         TxpoolConfiguration
	Consensus      ConsensusConfiguration
	Blockchain     BlockChainConfiguration
	Peer           PeerConfiguration
	Authentication AuthenticationConfiguration
	Icode          ICodeConfiguration
}

func NewConfiguration() Configuration {
	return Configuration{
		Common:         common.NewCommonConfiguration(),
		Txpool:         NewTxpoolConfiguration(),
		Consensus:      NewConsensusConfiguration(),
		Blockchain:     NewBlockChainConfiguration(),
		Peer:           NewPeerConfiguration(),
		Authentication: NewAuthenticationConfiguration(),
		Icode:          NewIcodeConfiguration(),
	}
}
