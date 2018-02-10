package domain

import (
	"time"
	pb "it-chain/network/protos"
	"sync"
	"sync/atomic"
)

type MsgType int

const (
	PreprepareMsg  MsgType = iota
	PrepareMsg
	CommitMsg
)

type ConsensusStateBuilder interface {
	ConsensusID(string) ConsensusStateBuilder
	CurrentStage(Stage) ConsensusStateBuilder
	View(*View) ConsensusStateBuilder
	EndConsensusHandler(EndConsensusHandle) ConsensusStateBuilder
	Period(int32) ConsensusStateBuilder
	Block(*Block) ConsensusStateBuilder
	Build() *ConsensusState
}

type consensusStateBuilder struct {
	id                  string
	view                *View
	currentStage        Stage
	block               *Block
	endConsensusHandler EndConsensusHandle
	period              int32
}

func NewConsensusStateBuilder() ConsensusStateBuilder {
	csb := &consensusStateBuilder{}
	csb.view = nil
	csb.block = nil
	return csb
}

func (csb consensusStateBuilder) ConsensusID(id string) ConsensusStateBuilder{
	csb.id = id
	return csb
}

func (csb consensusStateBuilder) CurrentStage(stage Stage) ConsensusStateBuilder{
	csb.currentStage = stage
	return csb
}

func (csb consensusStateBuilder) View(view *View) ConsensusStateBuilder{
	csb.view = view
	return csb
}

func (csb consensusStateBuilder) EndConsensusHandler(endConsensusHandler EndConsensusHandle) ConsensusStateBuilder{
	csb.endConsensusHandler = endConsensusHandler
	return csb
}

func (csb consensusStateBuilder) Period(period int32) ConsensusStateBuilder{
	csb.period = period
	return csb
}

func (csb consensusStateBuilder) Block(block *Block) ConsensusStateBuilder{
	csb.block = block
	return csb
}

func (csb consensusStateBuilder) Build() *ConsensusState{

	cs := &ConsensusState{
		ID:csb.id,
		View:csb.view,
		CurrentStage:csb.currentStage,
		Block: csb.block,
		PrepareMsgs: make(map[string]ConsensusMessage),
		CommitMsgs: make(map[string]ConsensusMessage),
		endConsensusHandler: csb.endConsensusHandler,
		IsEnd: int32(0),
		period: csb.period,
	}

	go cs.start()

	return cs
}

//consesnsus message can has 3 types
type ConsensusMessage struct {
	ConsensusID string
	View        View
	SequenceID  int64
	Block       *Block
	SenderID    string
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
	View                *View
	CurrentStage        Stage
	Block               *Block
	PrepareMsgs         map[string]ConsensusMessage
	CommitMsgs          map[string]ConsensusMessage
	endChannel          chan struct{}
	endConsensusHandler EndConsensusHandle
	IsEnd               int32
	period              int32
	sync.RWMutex
}

//현재 Consensus에 참여하는 leader와 peer정보
type View struct {
	ID       string
	LeaderID string
	PeerID   []string
}

//tested
func NewConsensusState(view *View, consensusID string, block *Block, currentStage Stage, endConsensusHandler EndConsensusHandle, periodSeconds int32) *ConsensusState{

	cs := &ConsensusState{
		ID:consensusID,
		View:view,
		CurrentStage:currentStage,
		Block: block,
		PrepareMsgs: make(map[string]ConsensusMessage),
		CommitMsgs: make(map[string]ConsensusMessage),
		endConsensusHandler: endConsensusHandler,
		IsEnd: int32(0),
		period: periodSeconds,
	}

	go cs.start()

	return cs
}

//tested
func NewConsesnsusMessage(consensusID string, view View,sequenceID int64, block *Block,peerID string, msgType MsgType) ConsensusMessage{

	return ConsensusMessage{
		ConsensusID: consensusID,
		View: view,
		SequenceID: sequenceID,
		MsgType:msgType,
		TimeStamp: time.Now(),
		SenderID:peerID,
		Block: block,
	}
}



//todo block을 넣어야함
//todo View를 넣어야함
func FromConsensusProtoMessage(consensusMessage pb.ConsensusMessage) ConsensusMessage{

	return ConsensusMessage{
		SequenceID: consensusMessage.SequenceID,
		SenderID: consensusMessage.SenderID,
		ConsensusID: consensusMessage.ConsensusID,
		MsgType: MsgType(consensusMessage.MsgType),
	}
}

//timer의 time이 끝나면 consensus를 종료한다.
//Consensus timer는 new 했을 때 시작된다.
func (cs *ConsensusState) start(){
	time.Sleep(time.Duration(cs.period)*time.Second)
	cs.Lock()
	defer cs.Unlock()

	//time out
	//consensus did not end
	//need to delete
	if cs.IsEnd == 0{
		cs.endConsensusHandler(*cs)
	}
}

func (cs *ConsensusState) End(){
	atomic.StoreInt32(&(cs.IsEnd), int32(1))
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
		cs.View = &consensusMessage.View
		cs.CurrentStage = Prepared
		break

	case PrepareMsg:
		_ ,ok := cs.PrepareMsgs[consensusMessage.SenderID]

		if !ok{
			cs.PrepareMsgs[consensusMessage.SenderID] = consensusMessage
		}
		break

	case CommitMsg:

		_ ,ok := cs.CommitMsgs[consensusMessage.SenderID]

		if !ok{
			cs.CommitMsgs[consensusMessage.SenderID] = consensusMessage
		}

		break
	default:
		break
	}
}

func (cs *ConsensusState) PrepareReady() bool{
	totalVotes := len(cs.View.PeerID)
	nowVotes := len(cs.PrepareMsgs)

	if nowVotes/totalVotes > 0.3{
		return true
	}
	return false
}
func (cs *ConsensusState) CommitReady() bool{
	totalVotes := len(cs.View.PeerID)
	nowVotes := len(cs.CommitMsgs)

	if nowVotes/totalVotes > 0.3{
		return true
	}
	return false
}

