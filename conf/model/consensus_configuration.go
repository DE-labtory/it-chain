package model

type ConsensusConfiguration struct {
	BatchTime       int
	MaxTransactions int
}

func NewConsensusConfiguration() ConsensusConfiguration {
	return ConsensusConfiguration{
		BatchTime:       3,
		MaxTransactions: 100,
	}
}
