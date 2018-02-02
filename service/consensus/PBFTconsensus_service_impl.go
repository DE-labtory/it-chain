package consensus

import (
	"it-chain/service/blockchain"
	"it-chain/network/comm"
	"sync"
	"github.com/rs/xid"
	"it-chain/service/peer"
	"fmt"
	"it-chain/common"
	"it-chain/auth"
)

type PBFTConsensusService struct {
	consensusStates map[string]*ConsensusState
	comm            comm.ConnectionManager
	view            View
	sequenceID      int64
	peerID          string
	peerService 	peer.PeerService
	crypto           auth.Crypto
	sync.RWMutex
}

func NewPBFTConsensusService(comm comm.ConnectionManager, peerService peer.PeerService,crypto auth.Crypto) ConsensusService{

	pbft := &PBFTConsensusService{
		consensusStates: make(map[string]*ConsensusState),
		comm:comm,
		sequenceID: 0,
		peerService: peerService,
		crypto:crypto,
	}

	return pbft
}

//not tested
//Consensus 시작
//1. Consensus의 state를 추가한다.
//2. 합의할 block을 consensusMessage에 담고 prepreMsg로 전파한다.
func (cs *PBFTConsensusService) StartConsensus(block *blockchain.Block){

	cs.Lock()
	//set consensus with preprepared state
	ConsensusState := NewConsensusState(cs.view.ID,xid.New().String(),block,PrePrepared)
	cs.consensusStates[ConsensusState.ID] = ConsensusState

	//set consensus message to broadcast
	preprepareConsensusMessage := NewConsesnsusMessage(cs.view.ID,cs.sequenceID,ConsensusState.Block,cs.peerID,PreprepareMsg)
	cs.sequenceID++

	cs.Unlock()

	cs.broadcastMessage(preprepareConsensusMessage)
}

func (cs *PBFTConsensusService) StopConsensus(){

}

func (cs *PBFTConsensusService) ReceiveConsensusMessage(consensusMsg *ConsensusMessage){

}

func (cs *PBFTConsensusService) consensusMessageHandler(){

}

//not tested
func (cs *PBFTConsensusService) broadcastMessage(consensusMsg ConsensusMessage){

	peerTable := cs.peerService.GetPeerTable()
	myInfo := peerTable.GetMyInfo()
	pubkey := myInfo.PubKey
	peerList := peerTable.GetPeerList()

	for _, peer := range peerList{
		envelope := common.ToEnvelope(consensusMsg,cs.crypto,pubkey)
		cs.comm.SendStream(envelope,nil,peer.PeerID)
	}
}