package domain

import (
	"time"
	"github.com/it-chain/it-chain-Engine/common"
	"github.com/pkg/errors"
	pb "github.com/it-chain/it-chain-Engine/legacy/network/protos"
	"sync"
)

type RaftState int32
type ElectionMsgType int32

const(
	Follower 	RaftState = iota
	Candidate
	Leader
	ShutDown
)

const (
	HeartBeatMsg   ElectionMsgType = iota
	RequestVoteMsg
	VoteMsg
)

type Raft struct {
	nodeId         string
	state          RaftState
	term           int64
	voteCount      int
	votedFor       string
	ElectionTimer  time.Timer
	HeartbeatTimer time.Timer
	lastBlockHash  string
	leaderId       string
	peerIds        []string
	sync.RWMutex
}

func NewRaft(nodeId string) *Raft{
	le := &Raft{
		nodeId: nodeId, state: Follower, term: 0, voteCount: 0, votedFor: "", lastBlockHash: "", leaderId: "", peerIds: make([]string, 0),
	}
	le.ElectionTimer.Stop()
	le.HeartbeatTimer.Stop()
	return le
}

func (r *Raft) GetNodeId() string {
	return r.nodeId
}

func (r *Raft) SetState(state RaftState) {
	r.state = state
}

func (r *Raft) GetState() RaftState {
	return r.state
}

func (r *Raft) SetTerm(term int64) {
	r.term = term
}

func (r *Raft) GetTerm() int64 {
	return r.term
}

func (r *Raft) CountTerm() {
	r.term++
}

func (r *Raft) GetVoteCount() int {
	return r.voteCount
}

func (r *Raft) CountVote() {
	r.voteCount++
}

func (r *Raft) Voting(peerId string) {
	r.votedFor = peerId
}

func (r *Raft) GetVotedFor() string {
	return r.votedFor
}

func (r *Raft) GetLastBlockHash() string {
	return r.lastBlockHash
}

func (r *Raft) SetLastBlockHash(hash string) {
	r.lastBlockHash = hash
}

func (r *Raft) VoteValidate(candidateRaft *Raft) (bool, error) {
	if r.votedFor != ""{
		return false, errors.New("already voted")
	}
	if r.GetLastBlockHash() != candidateRaft.GetLastBlockHash() {
		return false, errors.New("lastblockhash is different")
	}
	return true, nil
}

func (r *Raft) VotesForItself() {
	r.votedFor = r.nodeId
	r.voteCount++
}

func (r *Raft) SetLeaderId(leaderId string) {
	r.leaderId = leaderId
}

func (r *Raft) GetLeaderId() string {
	return r.leaderId
}

func (r *Raft) ResetVote() {
	r.votedFor = ""
	r.voteCount = 0
}

func (r *Raft) ResetElectionTimer() {
	r.ElectionTimer.Reset(time.Duration(common.CryptoRandomGeneration(150, 300)) * time.Millisecond)
}

func (r *Raft) SetHeartbeatTimer(t time.Duration) {
	r.HeartbeatTimer.Reset(t)
}

func (r *Raft) GetPeerId() []string {
	return r.peerIds
}

func (r *Raft) SetPeerId(peerids []string) {
	r.peerIds = peerids
}

func (r *Raft) AppendPeerId(peerId string) {
	r.peerIds = append(r.peerIds, peerId)
}

type ElectionMessage struct {
	LastBlockHash string
	SenderID      string
	MsgType       ElectionMsgType
	Term          int64
	PeerIDs       []string
}

func FromElectionProtoMessage(electionMessage pb.ElectionMessage) *ElectionMessage{

	return &ElectionMessage{
		LastBlockHash: electionMessage.LastBlockHash,
		SenderID:      electionMessage.SenderID,
		MsgType:       ElectionMsgType(electionMessage.MsgType),
		Term:          int64(electionMessage.Term),
		PeerIDs:       electionMessage.PeerIDs,
	}
}

func ToElectionProtoMessage(electionMessage ElectionMessage) *pb.ElectionMessage{

	return &pb.ElectionMessage{
		LastBlockHash: electionMessage.LastBlockHash,
		SenderID:      electionMessage.SenderID,
		MsgType:       int32(electionMessage.MsgType),
		Term:          int64(electionMessage.Term),
		PeerIDs:       electionMessage.PeerIDs,
	}
}

func (r *Raft) StopElectionTimer() {
	r.ElectionTimer.Stop()
}

func (r *Raft) StopHeartbeatTimer() {
	r.HeartbeatTimer.Stop()
}

func NewElectionMessage(r *Raft) ElectionMessage {
	msgType := func() ElectionMsgType {
		switch r.GetState() {
		case Follower:
			return VoteMsg
		case Candidate:
			return RequestVoteMsg
		case Leader:
			return HeartBeatMsg
		default:
			return -1
		}
	}

	return ElectionMessage{
		LastBlockHash: r.GetLastBlockHash(),
		SenderID:      r.GetNodeId(),
		MsgType:       msgType(),
		Term:          r.GetTerm(),
		PeerIDs:       r.peerIds,
	}
}