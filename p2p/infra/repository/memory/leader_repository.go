package memory

import (
	"errors"
	"github.com/it-chain/it-chain-Engine/p2p"
	"sync"
)

var ErrNoLeader = errors.New("there is no leader")

type LeaderRepoitory interface{

	GetLeader() p2p.Leader
	SetLeader(leader p2p.Leader) error
}

type LeaderRepository struct {
	leader p2p.Leader
	mux    sync.Mutex
}

func NewLeaderRepository(leader p2p.Leader) *LeaderRepository {

	return &LeaderRepository{
		leader: leader,
	}
}

// get leader method
func (lr *LeaderRepository) GetLeader() p2p.Leader {

	lr.mux.Lock()
	defer lr.mux.Unlock()

	return lr.leader
}

// set leader method
func (lr *LeaderRepository) SetLeader(leader p2p.Leader) {


	lr.mux.Lock()
	defer lr.mux.Unlock()

	lr.leader = leader
}
