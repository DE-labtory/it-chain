package domain

import (
	"time"
	"it-chain/common"
	"github.com/pkg/errors"
)

type RaftState int32

const(
	Follower 	RaftState = iota
	Candidate
	Leader
	ShutDown
)

type Raft struct {
	nodeId         string
	state          RaftState
	term           int
	voteCount      int
	votedFor       string
	electionTimer  time.Timer
	heartbeatTimer time.Timer
	lastBlockHash  string
	leaderId       string
	peerIds        []string
	//ElectionTimeout  time.Duration
	//HeartbeatTimeout time.Duration
}

func NewRaft(nodeId string, lastBlockHash string) *Raft{
	le := &Raft{
		nodeId: nodeId, state: Follower, term: 0, voteCount: 0, votedFor: "", lastBlockHash: lastBlockHash, leaderId: "", peerIds: make([]string, 0),
	}
	return le
}

func (r *Raft) SetRaftState(state RaftState) {
	r.state = state
}

func (r *Raft) GetRaftState() RaftState {
	return r.state
}

func (r *Raft) GetRaftTerm() int {
	return r.term
}

func (r *Raft) CountTerm() {
	r.term++
}

func (r *Raft) CountVote(votes int) {
	r.voteCount += votes
}

func (r *Raft) GetLastBlockHash() string {
	return r.lastBlockHash
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

func (r *Raft) ResetElectionTimer() {
	r.electionTimer.Reset(time.Duration(common.CryptoRandomGeneration(150, 300)) * time.Millisecond)
}

