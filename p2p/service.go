package p2p

import (
	"math/rand"
	"sync"
	"time"

	"github.com/it-chain/it-chain-Engine/common"
	"github.com/it-chain/midgard"
	"github.com/rs/xid"
)

type PeerService interface {
	Dial(ipAddress string) error
	DeliverLeaderInfo(connectionId string, leader Leader) error
	DeliverPeerLeaderTable(connectionId string, peerLeaderTable PeerLeaderTable) error
}

type Publish func(exchange string, topic string, data interface{}) (err error) // 나중에 의존성 주입을 해준다.

type VotingService struct {
	leftTime       int64 //left time in millisecond
	state          string
	voteCount      int
	mux            sync.Mutex
	peerRepository PeerRepository
	publish        Publish
}

func (vs *VotingService) getLeftTime() int64 {

	vs.mux.Lock()
	defer vs.mux.Unlock()

	return vs.leftTime
}

func (vs *VotingService) resetLeftTime() {

	vs.mux.Lock()
	defer vs.mux.Unlock()

	vs.leftTime = genRandomInRange(150, 300)
}

//count down left time by tick millisecond  until 0
func (vs *VotingService) countDownLeftTimeBy(tick int64) {

	vs.mux.Lock()
	defer vs.mux.Unlock()

	if vs.leftTime == 0 {
		return
	}

	vs.leftTime = vs.leftTime - tick
}

func (vs *VotingService) setState(state string) {

	vs.mux.Lock()
	defer vs.mux.Unlock()

	vs.state = state
}

func (vs *VotingService) getState() string {

	vs.mux.Lock()
	defer vs.mux.Unlock()

	return vs.state
}

func (vs *VotingService) getVoteCount() int {

	vs.mux.Lock()
	defer vs.mux.Unlock()

	return vs.voteCount
}

func (vs *VotingService) resetVoteCount() {

	vs.mux.Lock()
	defer vs.mux.Unlock()

	vs.voteCount = 0
}

func (vs *VotingService) countUp() {

	vs.mux.Lock()
	defer vs.mux.Unlock()

	vs.voteCount = vs.voteCount + 1
}

func (vs *VotingService) ElectLeaderWithRaft() {
	//1. Start random timeout
	//2. timed out! alter state to 'candidate'
	//3. while ticking, count down leader repo left time
	//4. Send message having 'RequestVoteProtocol' to other node
	go StartRandomTimeOut(vs)
}

func StartRandomTimeOut(vs *VotingService) {

	timeoutNum := genRandomInRange(150, 300)
	timeout := time.After(time.Duration(timeoutNum) * time.Microsecond)
	tick := time.Tick(1 * time.Millisecond)

	for {
		select {

		case <-timeout:
			if vs.getState() == "Ticking" {

				vs.setState("Candidate")

				peerList, _ := vs.peerRepository.FindAll()

				connectionIds := make([]string, 0)

				for _, peer := range peerList {
					connectionIds = append(connectionIds, peer.PeerId.Id)
				}

				vs.deliverRequestVoteMessages(connectionIds)

			} else if vs.getState() == "Candidate" {
				//reset time and state chane candidate -> ticking when timed in candidate state
				vs.resetLeftTime()
				vs.setState("Ticking")
			}

		case <-tick:
			vs.countDownLeftTimeBy(1)
		}
	}
}

func (vs *VotingService) deliverRequestVoteMessages(connectionIds []string) error {

	requestVoteMessage := RequestVoteMessage{}

	grpcDeliverCommand, _ := CreateGrpcDeliverCommand("PeerTableDeliver", requestVoteMessage)

	for _, connectionId := range connectionIds {
		grpcDeliverCommand.Recipients = append(grpcDeliverCommand.Recipients, connectionId)
	}

	vs.publish("Command", "message.send", grpcDeliverCommand)

	return nil
}

func CreateGrpcDeliverCommand(protocol string, body interface{}) (GrpcDeliverCommand, error) {

	data, err := common.Serialize(body)

	if err != nil {
		return GrpcDeliverCommand{}, err
	}

	return GrpcDeliverCommand{
		CommandModel: midgard.CommandModel{
			ID: xid.New().String(),
		},
		Recipients: make([]string, 0),
		Body:       data,
		Protocol:   protocol,
	}, err
}

func genRandomInRange(min, max int64) int64 {

	rand.Seed(time.Now().Unix())

	return rand.Int63n(max-min) + min
}
