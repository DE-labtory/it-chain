package consensus

import (
	"it-chain/service/blockchain"
	"it-chain/network/comm"
	"sync"
	"github.com/rs/xid"
)

type PBFTConsensusService struct {
	ConsensusStates map[string]*ConsensusState
	Comm            comm.ConnectionManager
	View            View
	SequenceID      int64
	sync.RWMutex
}

func NewPBFTConsensusService(comm comm.ConnectionManager) ConsensusService{

	pbft := &PBFTConsensusService{
		ConsensusStates: make(map[string]*ConsensusState),
		Comm:comm,
		SequenceID: 0,
	}

	return pbft
}

//Consensus 시작
//1. Consensus의 state를 추가한다.
//2. 합의할 block을 consensusMessage에 담고 prepreMsg로 전파한다.
func (cs *PBFTConsensusService) StartConsensus(block *blockchain.Block){

	cs.Lock()
	//set consensus with preprepared state
	ConsensusState := NewConsensusState(cs.View.ID,xid.New().String())
	cs.ConsensusStates[ConsensusState.ID] = ConsensusState

	//set consensus message to broadcast
	preprepareConsensusMessage := NewConsesnsusMessage(cs.View.ID,cs.SequenceID,block)
	cs.SequenceID++

	cs.Unlock()

	cs.broadcastMessage(preprepareConsensusMessage)
}

func (cs *PBFTConsensusService) StopConsensus(){

}

func (cs *PBFTConsensusService) ReceiveConsensusMessage(consensusMsg *ConsensusMessage){

}

func (cs *PBFTConsensusService) consensusMessageHandler(){

}

func (cs *PBFTConsensusService) broadcastMessage(consensusMsg ConsensusMessage){

}