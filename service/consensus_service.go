package service

import "it-chain/domain"

type ConsensusService interface{

	//Consensus 시작
	StartConsensus(block *domain.Block)

	StopConsensus()
	//consensus메세지는 모두 이쪽으로 받는다.
	ReceiveConsensusMessage(consensusMsg *domain.ConsensusMessage)
}