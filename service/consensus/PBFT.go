package consensus

import (
	"it-chain/service/blockchain"
	"time"
)

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
	Block       *blockchain.Block
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
	ID             string
	ViewID         string
	CurrentStage   Stage
	PrepareMsgs    []*ConsensusMessage
	CommitMsgs     []*ConsensusMessage
}

type View struct{
	ID string
}

func NewConsensusState(viewID string, consensusID string) *ConsensusState{
	return &ConsensusState{
		ID:consensusID,
		ViewID:viewID,
		CurrentStage:PrePrepared,
		PrepareMsgs: make([]*ConsensusMessage,0),
		CommitMsgs: make([]*ConsensusMessage,0),
	}
}

func NewConsesnsusMessage(viewID string,sequenceID int64, block *blockchain.Block) ConsensusMessage{
	return ConsensusMessage{
		ConsensusID: "1",
		ViewID: viewID,
		SequenceID: sequenceID,
		Stage: PrePrepared,
		MsgType:PreprepareMsg,
		TimeStamp: time.Now(),
		PeerID:"0",
		Block: block,
	}
}
