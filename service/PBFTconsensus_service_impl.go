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

func NewPBFTConsensusService(view *domain.View,comm comm.ConnectionManager, peerService PeerService, blockService BlockService) ConsensusService{

	pbft := &PBFTConsensusService{
		consensusStates: make(map[string]*domain.ConsensusState),
		comm:comm,
		view:view,
		peerService: peerService,
		blockService: blockService,
	}

	return pbft
}

//tested
//Consensus 시작
//1. Consensus의 state를 추가한다.
//2. 합의할 block을 consensusMessage에 담고 prepreMsg로 전파한다.
//todo sequence 를 nano로 수정
func (cs *PBFTConsensusService) StartConsensus(block *domain.Block){

	cs.Lock()
	//set consensus with preprepared state
	consensusState := domain.NewConsensusState(cs.view,xid.New().String(),block,domain.PrePrepared,cs.HandleEndConsensus,300)
	cs.consensusStates[consensusState.ID] = consensusState

	//set consensus message to broadcast
	sequenceID := time.Now().UnixNano()
	preprepareConsensusMessage := domain.NewConsesnsusMessage(consensusState.ID,cs.view.ID,sequenceID,consensusState.Block,cs.peerID,domain.PreprepareMsg)

	cs.Unlock()

	cs.broadcastMessage(preprepareConsensusMessage)
}

func (cs *PBFTConsensusService) GetCurrentConsensusState() map[string]*domain.ConsensusState{
	return cs.consensusStates
}

func (cs *PBFTConsensusService) StopConsensus(){

}

//consensusMessage가 들어옴
//todo FromConsensusProtoMessage에서 block변환도 해야함
//todo 언제 Message를 무시해야 하는가 일단은 time laps는 5분
//todo time을 config로 부터 읽어야함
//todo 다음 block이 먼저 들어올 경우 고려해야함,
//todo 블록의 높이와 이전 블록 해시가 올바른지 확인
func (cs *PBFTConsensusService) ReceiveConsensusMessage(outterMessage comm.OutterMessage){

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
		logger_pbftservice.Errorln("time over(5min)")
		return
	}

	//2 consensus id check
	consensusID := consensusMessage.ConsensusID
	msgType := consensusMessage.MsgType
	consensusState, ok := cs.consensusStates[consensusID]

	cs.Lock()

	if !ok{
		//consensus state생성
		//prepremessage인 경우에만 block과 view를 세팅
		//var newConsensusState *domain.ConsensusState
		//if consensusMessage.MsgType == domain.PreprepareMsg{
		//	newConsensusState = domain.NewConsensusState(&consensusMessage.View,consensusMessage.ConsensusID,consensusMessage.Block,domain.Stage(msgType),cs.HandleEndConsensus,300)
		//}
		//
		//cs.consensusStates[newConsensusState.ID] = newConsensusState
	}

	//if !ok{
	//	//이미 state가 존재함
	//	//prepare or commit message
	//	consensusState.AddMessage(consensusMessage)
	//}else{
	//	//preprepare message를 받는경우
	//	//id가 다르면 check안함
	//	//block check
	//	newConsensusState := domain.NewConsensusState(cs.view.ID,consensusMessage.ConsensusID,nil,domain.Stage(msgType),cs.HandleEndConsensus)
	//	newConsensusState.AddMessage(consensusMessage)
	//	cs.consensusStates[newConsensusState.ID] = newConsensusState
	//}
	cs.Unlock()
}

func (cs *PBFTConsensusService) HandleEndConsensus(consensusState domain.ConsensusState){

}

//tested
func (cs *PBFTConsensusService) broadcastMessage(consensusMsg domain.ConsensusMessage){

	peerTable := cs.peerService.GetPeerTable()
	peerList := peerTable.GetPeerList()

	for _, peer := range peerList{
		cs.comm.SendStream(consensusMsg,nil,peer.PeerID)
	}
}