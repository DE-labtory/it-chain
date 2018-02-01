package consensus

import (
	"it-chain/service/blockchain"
	"it-chain/network/comm"
	"sync"
	"time"
)


type RequestMsg struct {
	block      *blockchain.Block
	SequenceID int64
}

type MsgType int

const (
	PreprepareMsg  MsgType = iota
	PrepareMsg
	CommitMsg
)

//consesnsus message can has 3 types
type ConsensusMessage struct {
	ConsensusID string
	ViewID      string
	SequenceID  int64
	RequestMsg  *RequestMsg
	Stage       Stage
	PeerID      string
	MsgType     MsgType
	TimeStamp   time.Time
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

type PBFTConsensusService struct {
	consensusStates map[string]*ConsensusState
	comm            comm.ConnectionManager
	View            View
	SequenceID      int64
	sync.RWMutex
}

func NewPBFTConsensusService(comm comm.ConnectionManager) ConsensusService{

	pbft := &PBFTConsensusService{
		consensusStates: make(map[string]*ConsensusState),
		comm:comm,
		SequenceID: 0,
	}

	return pbft
}

//Consensus 시작
func (cs *PBFTConsensusService) StartConsensus(block *blockchain.Block){

	requestMsg := &RequestMsg{
		SequenceID: cs.SequenceID,
		block: block,
	}

	//temp consensusID and PeerID
	consensusMessage := ConsensusMessage{
		ConsensusID: "1",
		ViewID: cs.View.ID,
		SequenceID: cs.SequenceID,
		Stage: PrePrepared,
		MsgType:PreprepareMsg,
		TimeStamp: time.Now(),
		PeerID:"0",
		RequestMsg: requestMsg,
	}

	cs.Lock()
	cs.SequenceID++
	cs.Unlock()

	cs.broadcastMessage(consensusMessage)
}

func (cs *PBFTConsensusService) StopConsensus(){

}

func (cs *PBFTConsensusService) ReceiveConsensusMessage(consensusMsg *ConsensusMessage){

}

func (cs *PBFTConsensusService) consensusMessageHandler(){

}

func (cs *PBFTConsensusService) broadcastMessage(consensusMsg ConsensusMessage){

}