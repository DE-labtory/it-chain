package model

type Configuration struct {
	txpool         TxpoolConfiguration
	consensus      ConsensusConfiguration
	blockchain     BlockChainConfiguration
	peer           PeerConfiguration
	authentication AuthenticationConfiguration
	icode          ICodeConfiguration
}
