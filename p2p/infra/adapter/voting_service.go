package adapter

import (
	"math/rand"
	"sync"
	"time"
	"github.com/it-chain/it-chain-Engine/p2p"
)

type VotingService struct {
	leftTime           int64 //left time in millisecond
	state              string
	voteCount          int
	mux                sync.Mutex
	publish            Publish
	leaderRepository   p2p.LeaderRepository
	peerRepository     p2p.PeerRepository
	grpcCommandService GrpcCommandService
}

func (vs *VotingService) GetLeftTime() int64 {

	vs.mux.Lock()
	defer vs.mux.Unlock()

	return vs.leftTime

}

func (vs *VotingService) ResetLeftTime() {

	vs.mux.Lock()
	defer vs.mux.Unlock()

	vs.leftTime = GenRandomInRange(150, 300)

}

//count down left time by tick millisecond  until 0
func (vs *VotingService) CountDownLeftTimeBy(tick int64) {

	vs.mux.Lock()
	defer vs.mux.Unlock()
	if vs.leftTime == 0 {
		return
	}
	vs.leftTime = vs.leftTime - tick
}

func (vs *VotingService) SetState(state string) {

	vs.mux.Lock()
	defer vs.mux.Unlock()

	vs.state = state
}

func (vs *VotingService) GetState() string {

	vs.mux.Lock()
	defer vs.mux.Unlock()

	return vs.state
}

func (vs *VotingService) GetVoteCount() int {

	vs.mux.Lock()
	defer vs.mux.Unlock()

	return vs.voteCount
}

func (vs *VotingService) ResetVoteCount() {

	vs.mux.Lock()
	defer vs.mux.Unlock()

	vs.voteCount = 0
}

func (vs *VotingService) CountUp() {

	vs.mux.Lock()
	defer vs.mux.Unlock()

	vs.voteCount = vs.voteCount + 1
}

func (vs *VotingService) DeliverRequestVoteMessages(connectionIds []string) error {

	requestVoteMessage := p2p.RequestVoteMessage{}

	grpcDeliverCommand, _ := CreateGrpcDeliverCommand("PeerTableDeliver", requestVoteMessage)

	for _, connectionId := range connectionIds {
		grpcDeliverCommand.Recipients = append(grpcDeliverCommand.Recipients, connectionId)
	}

	vs.publish("Command", "message.send", grpcDeliverCommand)

	return nil

}

func (vs *VotingService) DeliverVoteLeaderMessage(connectionId string) error {
	voteLeaderMessage := p2p.VoteLeaderMessage{}

	grpcDeliverCommand, _ := CreateGrpcDeliverCommand("VoteLeaderProtocol", voteLeaderMessage)

	grpcDeliverCommand.Recipients = append(grpcDeliverCommand.Recipients, connectionId)

	vs.publish("Command", "message.send", grpcDeliverCommand)

	return nil

}

func (vs *VotingService) DeliverUpdateLeaderMessage(connectionId string, peer p2p.Peer) error {

	updateLeaderMessage := p2p.UpdateLeaderMessage{}

	grpcDeliverCommand, _ := CreateGrpcDeliverCommand("UpdateLeaderProtocol", updateLeaderMessage)

	vs.publish("Command", "message.deliver", grpcDeliverCommand)

	return nil

}

func (vs *VotingService) ElectLeaderWithRaft() {
	//1. Start random timeout
	//2. timed out! alter state to 'candidate'
	//3. while ticking, count down leader repo left time
	//4. Send message having 'RequestVoteProtocol' to other node
	go StartRandomTimeOut(vs)
}

func StartRandomTimeOut(vs *VotingService) {

	timeoutNum := GenRandomInRange(150, 300)
	timeout := time.After(time.Duration(timeoutNum) * time.Microsecond)
	tick := time.Tick(1 * time.Millisecond)

	for {
		select {

		case <-timeout:
			if vs.GetState() == "Ticking" {

				vs.SetState("Candidate")

				peerList, _ := vs.peerRepository.FindAll()

				connectionIds := make([]string, 0)

				for _, peer := range peerList {
					connectionIds = append(connectionIds, peer.PeerId.Id)
				}

				vs.DeliverRequestVoteMessages(connectionIds)

			} else if vs.GetState() == "Candidate" {

				//reset time and state chane candidate -> ticking when timed in candidate state
				vs.ResetLeftTime()
				vs.SetState("Ticking")

			}

		case <-tick:

			vs.CountDownLeftTimeBy(1)
		}
	}
}


func GenRandomInRange(min, max int64) int64 {

	rand.Seed(time.Now().Unix())

	return rand.Int63n(max-min) + min
}

