package memory

import (
	"errors"
	"math/rand"
	"sync"
	"time"

	"github.com/it-chain/it-chain-Engine/p2p"
)

var ErrNoLeader = errors.New("there is no leader")

type LeaderRepository struct {
	leader    p2p.Leader
	leftTime  int64 //left time in millisecond
	state     string
	voteCount int
	mux       sync.Mutex
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

func (lr *LeaderRepository) GetLeftTime() int64 {

	lr.mux.Lock()
	defer lr.mux.Unlock()

	return lr.leftTime

}

func (lr *LeaderRepository) ResetLeftTime() {

	lr.mux.Lock()
	defer lr.mux.Unlock()

	lr.leftTime = GenRandomInRange(150, 300)

}

//count down left time by tick millisecond  until 0
func (lr *LeaderRepository) CountDownLeftTimeBy(tick int64) {

	lr.mux.Lock()
	defer lr.mux.Unlock()
	if lr.leftTime == 0 {
		return
	}
	lr.leftTime = lr.leftTime - tick
}

func (lr *LeaderRepository) SetState(state string) {

	lr.mux.Lock()
	defer lr.mux.Unlock()

	lr.state = state

}

func (lr *LeaderRepository) GetState() string {

	lr.mux.Lock()
	defer lr.mux.Unlock()

	return lr.state
}

func (lr *LeaderRepository) GetVoteCount() int {

	lr.mux.Lock()
	defer lr.mux.Unlock()

	return lr.voteCount

}

func (lr *LeaderRepository) ResetVoteCount() {

	lr.mux.Lock()
	defer lr.mux.Unlock()

	lr.voteCount = 0

}

func (lr *LeaderRepository) CountUp() {

	lr.mux.Lock()
	defer lr.mux.Unlock()

	lr.voteCount = lr.voteCount + 1

}

func GenRandomInRange(min, max int64) int64 {

	rand.Seed(time.Now().Unix())

	return rand.Int63n(max-min) + min
}
