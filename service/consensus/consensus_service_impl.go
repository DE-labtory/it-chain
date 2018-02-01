package consensus

import "it-chain/service/blockchain"


type RequestMsg struct {
	Timestamp  int64
	block      blockchain.Block
	SequenceID int64
}

type MsgType int

const (
	PreprepareMsg  MsgType = iota
	PrepareMsg
	CommitMsg
)

//consesnsus message can has 3 types
type ConsensusMessage struct{
	ConsensusID string
	ViewID     int64
	SequenceID int64
	RequestMsg *RequestMsg
	Stage 	   MsgType
	PeerID     string
	MsgType
}

//동시에 여러개의 consensus가 진행될 수 있다.
//한개의 consensus는 1개의 state를 갖는다.
type ConsensusState struct{
	ID string
}

type View struct{
	ID string
}


type ConsensusServiceImpl struct{
	ConsensusStates map[string]*ConsensusState
}


//Consensus 시작
func (cs *ConsensusServiceImpl) StartConsensus(transaction blockchain.Transaction){

}

//
func (cs *ConsensusServiceImpl) StopConsensus(){

}

func (cs *ConsensusServiceImpl) ReceiveConsensusMessage(consensusMsg *ConsensusMessage){

}