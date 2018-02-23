package service

import (
	"it-chain/network/comm"
	"it-chain/domain"
	"sync"
	pb "it-chain/network/protos"
	"it-chain/common"
	"it-chain/network/comm/msg"
)

var logger_raftservice = common.GetLogger("raft_service")

type RAFTelectionService struct {
	comm comm.ConnectionManager
	node *domain.Raft
	sync.RWMutex
}

func NewRAFTelectionService(comm comm.ConnectionManager, peerID string) LeaderElectionService{
	raft := &RAFTelectionService{
		comm:comm,
		node:domain.NewRaft(peerID),
	}
	return raft
}

func (r *RAFTelectionService) Run() {
	r.node.ResetElectionTimer()
	r.node.SetLastBlockHash(r.GetLastBlockHash())
	go func() {
		for {
			select {
			case <-r.node.ElectionTimer.C:
				r.node.Lock()
				r.node.SetState(domain.Candidate)
				r.node.CountTerm()
				r.node.ResetVote()
				r.node.VotesForItself()
				r.broadcastMessage(domain.NewElectionMessage(r.node))
				r.node.Unlock()
			case <-r.node.HeartbeatTimer.C:

			default:
			}
		}
	}()
}

func (r *RAFTelectionService) Stop() {

}

// 별로의 무한루프 쓰레드열고 채널로 데이터 메세지가 들어오면 리시브 메세지가 메세지를 받아서 처리하는 것인지 comm 찾아보고 결정
func (r *RAFTelectionService) ReceiveMessage(message msg.OutterMessage) {
	logger_raftservice.Infoln("Received message: ", message)

	msg := message.Message
	em 	:= msg.GetElectionMessage()

	if em == nil{
		logger_raftservice.Errorln("election Message is empty")
		return
	}

	electionMessage := domain.FromElectionProtoMessage(*em)

	switch electionMessage.MsgType {
	case domain.HeartBeatMsg:
		if r.node.GetTerm() < electionMessage.Term && r.node.GetLastBlockHash() == electionMessage.LastBlockHash {
			r.node.Lock()
			r.node.SetTerm(electionMessage.Term)
			r.node.ResetVote()
			r.node.SetLeaderId(electionMessage.SenderID)
			r.node.SetState(domain.Follower)
			r.node.ResetElectionTimer()
			// 메세지로 받은 peerID 리스트 처리를 어떻게 할지 고민..
			r.node.Unlock()
		} else if r.node.GetLeaderId() == electionMessage.SenderID {
			r.node.Lock()
			r.node.ResetElectionTimer()
			r.node.Unlock()
		}
	case domain.RequestVoteMsg:
		if r.node.GetVotedFor() == "" && r.node.GetLastBlockHash() == electionMessage.LastBlockHash {
			r.node.Lock()
			// electionMessage.SenderID에게 vote 전송
			r.sendVoteMsg(electionMessage.SenderID)
			r.node.ResetElectionTimer()
			r.node.Unlock()
		}
	case domain.VoteMsg:
		if r.node.GetState() == domain.Candidate {
			r.node.Lock()
			r.node.CountVote()
			if r.node.GetVoteCount() > (len(r.node.GetPeerId()) / 2) {
				r.node.SetState(domain.Leader)
				r.node.SetLeaderId(r.node.GetNodeId())
				r.node.ResetVote()
				r.node.StopElectionTimer()
			}
			r.node.Unlock()
		}
	default:
		break
	}
}

func (r *RAFTelectionService) broadcastMessage(electionMsg domain.ElectionMessage){
	logger_raftservice.Infoln("Election broadcast Message")
	peerIDList := r.node.GetPeerId()
	message := &pb.StreamMessage{
		Content: &pb.StreamMessage_ElectionMessage{ ElectionMessage: domain.ToElectionProtoMessage(electionMsg), },
	}
	for _, peerID := range peerIDList{
		logger_raftservice.Infoln("sending...",peerID)
		r.comm.SendStream(message,nil, nil,peerID)
	}
}

func (r *RAFTelectionService) sendVoteMsg(peerId string) {
	logger_raftservice.Infoln("Election send vote Message")
	voteMsg := domain.NewElectionMessage(r.node)
	message := &pb.StreamMessage{
		Content: &pb.StreamMessage_ElectionMessage{ElectionMessage: domain.ToElectionProtoMessage(voteMsg),},
	}
	logger_raftservice.Infoln("sending...", peerId)
	r.comm.SendStream(message, nil, nil, peerId)
}

func (r *RAFTelectionService) AddPeerId() {

}

func (r *RAFTelectionService) GetLastBlockHash() string {
	//BlockService().GetLastBlock()
	return ""
}