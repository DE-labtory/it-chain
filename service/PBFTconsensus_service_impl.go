package service

import (
	"it-chain/network/comm"
	"sync"
	"github.com/rs/xid"
	"it-chain/domain"
	pb "it-chain/network/protos"
	"golang.org/x/text/message"
	"github.com/gogo/protobuf/proto"
)

//todo peerID를 어디서 가져올 것인가??
type PBFTConsensusService struct {
	consensusStates map[string]*domain.ConsensusState
	comm            comm.ConnectionManager
	view            domain.View
	sequenceID      int64
	peerID          string
	peerService 	PeerService
	sync.RWMutex
}

func NewPBFTConsensusService(comm comm.ConnectionManager, peerService PeerService) ConsensusService{

	pbft := &PBFTConsensusService{
		consensusStates: make(map[string]*domain.ConsensusState),
		comm:comm,
		sequenceID: 0,
		peerService: peerService,
	}

	return pbft
}

//tested
//Consensus 시작
//1. Consensus의 state를 추가한다.
//2. 합의할 block을 consensusMessage에 담고 prepreMsg로 전파한다.
func (cs *PBFTConsensusService) StartConsensus(block *domain.Block){

	cs.Lock()
	//set consensus with preprepared state
	ConsensusState := domain.NewConsensusState(cs.view.ID,xid.New().String(),block,domain.PrePrepared)
	cs.consensusStates[ConsensusState.ID] = ConsensusState

	//set consensus message to broadcast
	preprepareConsensusMessage := domain.NewConsesnsusMessage(cs.view.ID,cs.sequenceID,ConsensusState.Block,cs.peerID,domain.PreprepareMsg)
	cs.sequenceID++

	cs.Unlock()

	cs.broadcastMessage(preprepareConsensusMessage)
}

func (cs *PBFTConsensusService) GetCurrentConsensusState() map[string]*domain.ConsensusState{
	return cs.consensusStates
}

func (cs *PBFTConsensusService) StopConsensus(){

}

//consensusMessage가 들어옴
func (cs *PBFTConsensusService) ReceiveConsensusMessage(outterMessage comm.OutterMessage){

	message := outterMessage.Message

	message.gET

	switch msgType{
	case domain.PreprepareMsg:
		return
	case domain.PrepareMsg:
		return
	case domain.CommitMsg:
		return
	default:
		return
	}

}

func (cs *PBFTConsensusService) consensusMessageHandler(){

}

//tested
func (cs *PBFTConsensusService) broadcastMessage(consensusMsg domain.ConsensusMessage){

	peerTable := cs.peerService.GetPeerTable()
	peerList := peerTable.GetPeerList()

	for _, peer := range peerList{
		cs.comm.SendStream(consensusMsg,nil,peer.PeerID)
	}
}