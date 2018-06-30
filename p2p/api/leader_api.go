package api

import (
	"errors"
	"log"
	"math/rand"
	"time"

	"github.com/it-chain/it-chain-Engine/p2p"
	"github.com/it-chain/midgard"
)

var ErrEmptyPeerId = errors.New("empty peer id requested")
var ErrEmptyLeaderId = errors.New("empty leader id proposed")
var ErrEmptyConnectionId = errors.New("empty connection id proposed")

type LeaderApi struct {
	leaderRepository   ReadOnlyLeaderRepository
	peerRepository     ReadOnlyPeerRepository
	eventRepository    EventRepository
	grpcCommandService LeaderGrpcCommandService
	myInfo             *p2p.Peer
}

type Publish func(exchange string, topic string, data interface{}) (err error) // 나중에 의존성 주입을 해준다.

type ReadOnlyLeaderRepository interface {
	GetLeader() p2p.Leader
	CountDownLeftTimeBy(tick int64)
	SetState(state string)
	GetState() string
	ResetLeftTime() int64
	GetLeftTime() int64
	CountUp()
	GetVoteCount() int
}

type EventRepository interface { //midgard.Repository
	Save(aggregateID string, events ...midgard.Event) error
}

type LeaderGrpcCommandService interface {
	DeliverLeaderInfo(connectionId string, leader p2p.Leader) error
	DeliverRequestVoteMessages(connectionIds []string) error
}

func NewLeaderApi(leaderRepository ReadOnlyLeaderRepository, peerRepository ReadOnlyPeerRepository, eventRepository EventRepository, grpcCommandService LeaderGrpcCommandService, myInfo *p2p.Peer) *LeaderApi {

	return &LeaderApi{
		leaderRepository:   leaderRepository,
		peerRepository:     peerRepository,
		eventRepository:    eventRepository,
		grpcCommandService: grpcCommandService,
		myInfo:             myInfo,
	}
}

func (leaderApi *LeaderApi) UpdateLeader(leader p2p.Leader) error {

	if leader.LeaderId.Id == "" {
		return ErrEmptyLeaderId
	}

	events := make([]midgard.Event, 0)
	leaderUpdatedEvent := p2p.LeaderUpdatedEvent{
		EventModel: midgard.EventModel{
			ID:   leader.LeaderId.ToString(),
			Type: "leader.update",
		},
	}

	events = append(events, leaderUpdatedEvent)
	err := leaderApi.eventRepository.Save(leaderUpdatedEvent.GetID(), events...)

	if err != nil {
		log.Println(err.Error())
		return err
	}

	return err
}

func (leaderApi *LeaderApi) DeliverLeaderInfo(connectionId string) error {

	if connectionId == "" {
		return ErrEmptyConnectionId
	}

	leader := leaderApi.leaderRepository.GetLeader()
	leaderApi.grpcCommandService.DeliverLeaderInfo(connectionId, leader)

	return nil
}

func (leaderApi *LeaderApi) ElectLeaderWithRaft() {
	//1. Start random timeout
	//2. timed out! alter state to 'candidate'
	//3. while ticking, count down leader repo left time
	//4. Send message having 'RequestVoteProtocol' to other node
	go StartRandomTimeOut(leaderApi)
}

//todo find connectionId by peerId of make peer repo contains connectionId
func StartRandomTimeOut(leaderApi *LeaderApi) {

	timeoutNum := GenRandomInRange(150, 300)
	timeout := time.After(time.Duration(timeoutNum) * time.Microsecond)
	tick := time.Tick(1 * time.Millisecond)

	for {
		select {

		case <-timeout:
			if leaderApi.leaderRepository.GetState() == "Ticking" {

				leaderApi.leaderRepository.SetState("Candidate")

				peerList, _ := leaderApi.peerRepository.FindAll()

				connectionIds := make([]string, 0)

				for _, peer := range peerList {
					connectionIds = append(connectionIds, peer.PeerId.Id)
				}

				leaderApi.grpcCommandService.DeliverRequestVoteMessages(connectionIds)

			} else if leaderApi.leaderRepository.GetState() == "Candidate" {

				leaderApi.leaderRepository.SetState("Ticking")

			}

		case <-tick:

			leaderApi.leaderRepository.CountDownLeftTimeBy(1)
		}
	}
}

func GenRandomInRange(min, max int64) int64 {

	rand.Seed(time.Now().Unix())

	return rand.Int63n(max-min) + min

}
