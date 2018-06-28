package leveldb

import (
	"github.com/it-chain/it-chain-Engine/p2p"
	"sync"
	"errors"
)

//inmemory leader info
var leader p2p.Leader
var once sync.Once

var ErrNoLeader = errors.New("there is no leader")

type LeaderRepository struct {}

func NewLeaderRepository(firstLeader p2p.Leader) *LeaderRepository {

	once.Do(func() {
		leader = firstLeader
	})
	return &LeaderRepository{}
}

// get leader method
func (lr *LeaderRepository) GetLeader() p2p.Leader {
	return leader
}

// set leader method
func (lr *LeaderRepository) SetLeader(newLeader p2p.Leader) {
	leader = newLeader
}
