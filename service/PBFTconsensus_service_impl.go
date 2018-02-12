package service

import (
	"time"
	"it-chain/common"
	"math"
	"it-chain/domain"
	"it-chain/network/comm"
	"sync"
	"github.com/rs/xid"
)

var logger_pbftservice = common.GetLogger("pbft_service")

//todo peerID를 어디서 가져올 것인가??
type PBFTConsensusService struct {
	consensusStates map[string]*domain.ConsensusState
	comm            comm.ConnectionManager
	view            *domain.View
	peerID          string
	peerService 	PeerService
	blockService    BlockService
	sync.RWMutex
}

func NewPBFTConsensusService(view *domain.View,comm comm.ConnectionManager, blockService BlockService) ConsensusService{

	pbft := &PBFTConsensusService{
		consensusStates: make(map[string]*domain.ConsensusState),
		comm:comm,
		view:view,
		blockService: blockService,
	}

	return pbft
}

//tested
//Consensus 시작
//만약 합의에 들어가는 peerID가 없다면 바로 block에 저장
//1. Consensus의 state를 추가한다.
//2. 합의할 block을 consensusMessage에 담고 prepreMsg로 전파한다.
func (cs *PBFTConsensusService) StartConsensus(block *domain.Block){

	if len(cs.view.PeerID) <= 1 && cs.view.LeaderID == cs.peerID{
		//ADD block

		return
	}

	cs.Lock()
	//set consensus with preprepared state
	consensusState := domain.NewConsensusState(cs.view,xid.New().String(),block,domain.PrePrepared,cs.EndConsensusState,300)
	cs.consensusStates[consensusState.ID] = consensusState

	//set consensus message to broadcast
	sequenceID := time.Now().UnixNano()
	preprepareConsensusMessage := domain.NewConsesnsusMessage(consensusState.ID,*cs.view,sequenceID,consensusState.Block,cs.peerID,domain.PreprepareMsg)
	cs.broadcastMessage(preprepareConsensusMessage)
	consensusState.CurrentStage = domain.Prepared

	cs.Unlock()
}

func (cs *PBFTConsensusService) GetCurrentConsensusState() map[string]*domain.ConsensusState{
	return cs.consensusStates
}

func (cs *PBFTConsensusService) StopConsensus(){

}

//consensusMessage가 들어옴
//todo FromConsensusProtoMessage에서 block변환도 해야함
//todo time을 config로 부터 읽어야함
//todo 다음 block이 먼저 들어올 경우 고려해야함,
//todo 블록의 높이와 이전 블록 해시가 올바른지 확인
func (cs *PBFTConsensusService) ReceiveConsensusMessage(outterMessage comm.OutterMessage){

	logger_pbftservice.Infoln("Received message: ",outterMessage)

	message := outterMessage.Message
	cm:= message.GetConsensusMessage()

	if cm == nil{
		logger_pbftservice.Errorln("consensus Message is empty")
		return
	}

	consensusMessage := domain.FromConsensusProtoMessage(*cm)
	//1 time check
	t := time.Unix(0, consensusMessage.SequenceID)
	elapsed := time.Since(t)

	if math.Abs(elapsed.Minutes()) > 5.0 {
		logger_pbftservice.Errorln("time over (5min)")
		return
	}

	logger_pbftservice.Infoln("Time check OK")

	//2 consensus id check
	cs.Lock()
	defer cs.Unlock()

	consensusID := consensusMessage.ConsensusID
	consensusState, ok := cs.consensusStates[consensusID]

	if !ok{
		//consensus state생성
		//prepremessage인 경우에만 block과 view, stage를 세팅
		//var newConsensusState *domain.ConsensusState
		consensusStateBuilder := domain.NewConsensusStateBuilder()

		if consensusMessage.MsgType == domain.PreprepareMsg{

			consensusState = consensusStateBuilder.
				ConsensusID(consensusMessage.ConsensusID).
				CurrentStage(domain.Prepared).
				View(&consensusMessage.View).
				Block(consensusMessage.Block).
				EndConsensusHandler(cs.EndConsensusState).
				Period(300).Build()

		}else{

			consensusState = consensusStateBuilder.
				ConsensusID(consensusMessage.ConsensusID).
				EndConsensusHandler(cs.EndConsensusState).
				Period(300).Build()
		}

		cs.consensusStates[consensusState.ID] = consensusState
	}

	logger_pbftservice.Infoln("Add message to consensusState")
	consensusState.AddMessage(consensusMessage)

	//1. prepare stage && prepare message가 전체의 2/3이상 -> commitMsg전파
	if consensusState.CurrentStage == domain.Prepared && consensusState.PrepareReady(){
		sequenceID := time.Now().UnixNano()
		commitConsensusMessage := domain.NewConsesnsusMessage(consensusState.ID,*cs.view,sequenceID,consensusState.Block,cs.peerID,domain.CommitMsg)
		cs.broadcastMessage(commitConsensusMessage)
		consensusState.CurrentStage = domain.Committed
		logger_pbftservice.Infoln("ConsensusState is Committed")
	}

	//2. commit state && commit message가 전체의 2/3이상 -> 블록저장
	if consensusState.CurrentStage == domain.Committed && consensusState.CommitReady(){
		//block 저장
		//todo block에 저장
		logger_pbftservice.Infoln("ConsesnsusState is End")
	}
}

func (cs *PBFTConsensusService) EndConsensusState(consensusState domain.ConsensusState){
	cs.Lock()
	defer cs.Unlock()

	delete(cs.consensusStates,consensusState.ID)
}

//tested
func (cs *PBFTConsensusService) broadcastMessage(consensusMsg domain.ConsensusMessage){

	peerIDList := cs.view.PeerID
	for _, peerID := range peerIDList{
		cs.comm.SendStream(consensusMsg,nil,peerID)
	}
}