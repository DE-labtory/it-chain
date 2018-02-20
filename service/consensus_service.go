package service

import (
	"it-chain/domain"
	"it-chain/network/comm/msg"
)

type ConsensusService interface{

	//Consensus 시작
	StartConsensus(view *domain.View, block *domain.Block)

	StopConsensus()
	//consensus메세지는 모두 이쪽으로 받는다.
	ReceiveConsensusMessage(message msg.OutterMessage)

	GetCurrentConsensusState() map[string]*domain.ConsensusState
}