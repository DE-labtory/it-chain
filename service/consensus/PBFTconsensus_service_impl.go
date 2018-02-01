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

type Stage int

const (
	Idle        Stage = iota // Node is created successfully, but the consensus process is not started yet.
	PrePrepared              // The ReqMsgs is processed successfully. The node is ready to head to the Prepare stage.
	Prepared                 // Same with `prepared` stage explained in the original paper.
	Committed                // Same with `committed-local` stage explained in the original paper.
)

//동시에 여러개의 consensus가 진행될 수 있다.
//한개의 consensus는 1개의 state를 갖는다.
type ConsensusState struct {
	ID           string
	ViewID       int64
	CurrentStage Stage
}

type View struct{
	ID string
}

type PBFTConsensusService struct{
	ConsensusStates map[string]*ConsensusState
}

//Consensus 시작
func (cs *PBFTConsensusService) StartConsensus(block *blockchain.Block){

}

func (cs *PBFTConsensusService) StopConsensus(){

}

func (cs *PBFTConsensusService) ReceiveConsensusMessage(consensusMsg *ConsensusMessage){

}