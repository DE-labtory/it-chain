package memory

import (
	"github.com/it-chain/it-chain-Engine/p2p"
	"errors"
	"sync"
)


var ErrNoLeader = errors.New("there is no leader")

type LeaderRepository struct {
	leader p2p.Leader
	leftTime int64
	state string
	mux sync.Mutex
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

func (lr *LeaderRepository) GetLeftTime() int64{
	lr.mux.Lock()
	defer lr.mux.Unlock()
	return lr.leftTime
}

func (lr *LeaderRepository) ResetLeftTime() int64{

	lr.mux.Lock()
	defer lr.mux.Unlock()

	return lr.leftTime
}

func (lr *LeaderRepository) SetState(state string){

	lr.mux.Lock()
	defer lr.mux.Unlock()

	lr.state=state
	
}

func (lr *LeaderRepository) GetState() string{

	lr.mux.Lock()
	defer lr.mux.Unlock()

	return lr.state
}