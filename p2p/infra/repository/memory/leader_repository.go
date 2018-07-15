package memory

import (
	"github.com/it-chain/it-chain-Engine/p2p"
	"sync"
)

type LeaderRepository struct {
	leader p2p.Leader
	mux    sync.Mutex
}

// set leader method
func (lr *LeaderRepository) SetLeader(leader p2p.Leader) {


	lr.mux.Lock()
	defer lr.mux.Unlock()

	lr.leader = leader
}
