package domain

import (
	"time"
	pb "it-chain/network/protos"
	"sync"
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
	Block       *Block
	PeerID      string
	MsgType     MsgType
	TimeStamp   time.Time
}

type Stage int

const (
	PrePrepared   Stage = iota           // The ReqMsgs is processed successfully. The node is ready to head to the Prepare stage.
	Prepared                			 // Same with `prepared` stage explained in the original paper.
	Committed                			 // Same with `committed-local` stage explained in the original paper.
)

type EndConsensusHandle func(ConsensusState)

//동시에 여러개의 consensus가 진행될 수 있다.
//한개의 consensus는 1개의 state를 갖는다.
type ConsensusState struct {
	ID                  string
	ViewID              string
	CurrentStage        Stage
	Block               *Block
	PrepareMsgs         []ConsensusMessage
	CommitMsgs          []ConsensusMessage
	endChannel          chan struct{}
	endConsensusHandler EndConsensusHandle
	sync.RWMutex
}

type View struct{
	ID string
}

//tested
func NewConsensusState(viewID string, consensusID string, block *Block, currentStage Stage, endConsensusHandler EndConsensusHandle) *ConsensusState{
	return &ConsensusState{
		ID:consensusID,
		ViewID:viewID,
		CurrentStage:currentStage,
		Block: block,
		PrepareMsgs: make([]ConsensusMessage,0),
		CommitMsgs: make([]ConsensusMessage,0),
		endConsensusHandler: endConsensusHandler,
	}
}

//tested
func NewConsesnsusMessage(consensusID string, viewID string,sequenceID int64, block *Block,peerID string, msgType MsgType) ConsensusMessage{

	return ConsensusMessage{
		ConsensusID: consensusID,
		ViewID: viewID,
		SequenceID: sequenceID,
		MsgType:msgType,
		TimeStamp: time.Now(),
		PeerID:peerID,
		Block: block,
	}
}

//todo block을 넣어야함
func FromConsensusProtoMessage(consensusMessage pb.ConsensusMessage) ConsensusMessage{

	return ConsensusMessage{
		ViewID: consensusMessage.ViewID,
		SequenceID: consensusMessage.SequenceID,
		PeerID: consensusMessage.PeerID,
		ConsensusID: consensusMessage.ConsensusID,
		MsgType: MsgType(consensusMessage.MsgType),
	}
}

//timer의 time이 끝나면 consensus를 종료한다.
func (cs *ConsensusState) startTimer(){

}

func (cs *ConsensusState) End(){

}

//message 종류에 따라서 다르게 넣어야함
func (cs *ConsensusState) AddMessage(consensusMessage ConsensusMessage){
	//PreprepareMsg는 block이 존재
	//commit, prepareMsg는 block 존재 안함
	//prepare가 2/3이상일 경우
	//commit이 2/3이상일 경우

	msgType := consensusMessage.MsgType

	switch msgType {
	case PreprepareMsg:
		cs.Block = consensusMessage.Block
		//prepareMsg broadcast 해야함
		break

	case PrepareMsg:
		cs.PrepareMsgs = append(cs.PrepareMsgs, consensusMessage)
		//commitMsg broadcast 해야함
		break

	case CommitMsg:
		cs.CommitMsgs = append(cs.CommitMsgs, consensusMessage)
		//block 저장해야함
		break
	default:
		break
	}
}

type Command interface{
	Execute()
}
