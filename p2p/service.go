package p2p

import (
	"math/rand"
	"sync"
	"time"

	"github.com/it-chain/it-chain-Engine/common"
	"github.com/it-chain/it-chain-Engine/conf"
	"github.com/it-chain/midgard"
	"github.com/rs/xid"
)

type CommunicationService interface {
	Dial(ipAddress string) error
	DeliverLeaderInfo(connectionId string, leader Leader) error
	DeliverPLTable(connectionId string, peerLeaderTable PLTable) error

}

type PeerService interface{
	
}

type Publish func(exchange string, topic string, data interface{}) (err error) // 나중에 의존성 주입을 해준다.

type ElectionService struct {
	mux                sync.Mutex
	electionRepository ElectionRepository
	peerRepository     PeerRepository
	publish            Publish
}

func (es *ElectionService) ElectLeaderWithRaft() {
	//1. Start random timeout
	//2. timed out! alter state to 'candidate'
	//3. while ticking, count down leader repo left time
	//4. Send message having 'RequestVoteProtocol' to other node
	go StartRandomTimeOut(es)
}

func StartRandomTimeOut(es *ElectionService) {

	timeoutNum := genRandomInRange(150, 300)
	timeout := time.After(time.Duration(timeoutNum) * time.Microsecond)
	tick := time.Tick(1 * time.Millisecond)
	election := es.electionRepository.GetElection()

	for {
		select {

		case <-timeout:
			// when timed out
			// 1. if state is ticking, be candidate and request vote
			// 2. if state is candidate, reset state and left time
			if election.GetState() == "Ticking" {

				election.SetState("Candidate")
				es.electionRepository.SetElection(election)

				peerList, _ := es.peerRepository.FindAll()

				connectionIds := make([]string, 0)

				for _, peer := range peerList {
					connectionIds = append(connectionIds, peer.PeerId.Id)
				}

				es.requestVote(connectionIds)

			} else if election.GetState() == "Candidate" {
				//reset time and state chane candidate -> ticking when timed in candidate state
				election.ResetLeftTime()
				election.SetState("Ticking")
			}

			es.electionRepository.SetElection(election)

		case <-tick:
			// count down left time while ticking
			election.CountDownLeftTimeBy(1)

			es.electionRepository.SetElection(election)

		}
	}
}

func (es *ElectionService) requestVote(connectionIds []string) error {

	// 1. create request vote message
	// 2. send message
	requestVoteMessage := RequestVoteMessage{}

	grpcDeliverCommand, _ := CreateGrpcDeliverCommand("PeerTableDeliver", requestVoteMessage)

	for _, connectionId := range connectionIds {

		grpcDeliverCommand.Recipients = append(grpcDeliverCommand.Recipients, connectionId)
	}

	es.publish("Command", "message.send", grpcDeliverCommand)

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

func (es *ElectionService) Vote(connectionId string) error {

	//if leftTime >0, reset left time and send VoteLeaderMessage
	election := es.electionRepository.GetElection()

	if election.GetLeftTime() < 0 {
		return nil
	}

	election.ResetLeftTime()

	voteLeaderMessage := VoteMessage{}

	grpcDeliverCommand, _ := CreateGrpcDeliverCommand("VoteLeaderProtocol", voteLeaderMessage)
	grpcDeliverCommand.Recipients = append(grpcDeliverCommand.Recipients, connectionId)

	es.publish("Command", "message.deliver", grpcDeliverCommand)

	return nil
}

func (es *ElectionService) BroadcastLeader(peer Peer) error {

	updateLeaderMessage := UpdateLeaderMessage{}

	grpcDeliverCommand, _ := CreateGrpcDeliverCommand("UpdateLeaderProtocol", updateLeaderMessage)

	peers, _ := es.peerRepository.FindAll()

	for _, peer := range peers {
		grpcDeliverCommand.Recipients = append(grpcDeliverCommand.Recipients, peer.PeerId.Id)
	}

	es.publish("Command", "message.deliver", grpcDeliverCommand)

	return nil
}

//broad case leader when voted fully
func (es *ElectionService) DecideToBeLeader(command GrpcReceiveCommand) error {
	election := es.electionRepository.GetElection()

	//	1. if candidate, reset left time
	//	2. count up
	if election.GetState() == "candidate" {

		election.CountUp()
		es.electionRepository.SetElection(election)
	}

	//	3. if counted is same with num of peer-1 set leader and publish
	peers, _ := es.peerRepository.FindAll()
	numOfPeers := len(peers)

	if election.GetVoteCount() == numOfPeers-1 {

		peer := Peer{
			PeerId:    PeerId{Id: ""},
			IpAddress: conf.GetConfiguration().Common.NodeIp,
		}

		es.BroadcastLeader(peer)
	}

	return nil
}
