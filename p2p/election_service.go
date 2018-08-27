/*
 * Copyright 2018 It-chain
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package p2p

import (
	"math/rand"
	"sync"
	"time"

	"github.com/it-chain/engine/common"
	"github.com/it-chain/engine/common/command"
	"github.com/it-chain/engine/common/logger"
	"github.com/rs/xid"
)

type ElectionService struct {
	mux              sync.Mutex
	Election         *Election
	peerQueryService PeerQueryService
	client           Client
}

func NewElectionService(election *Election, peerQueryService PeerQueryService, client Client) ElectionService {

	return ElectionService{
		mux:              sync.Mutex{},
		Election:         election,
		peerQueryService: peerQueryService,
		client:           client,
	}
}

func (es *ElectionService) Vote(connectionId string) error {

	peer, _ := es.peerQueryService.FindPeerById(PeerId{Id: connectionId})

	es.Election.SetCandidate(&peer)

	//if leftTime >0, reset left time and send VoteLeaderMessage

	es.Election.ResetLeftTime()

	voteLeaderMessage := VoteMessage{}

	grpcDeliverCommand, _ := CreateGrpcDeliverCommand("VoteLeaderProtocol", voteLeaderMessage)
	grpcDeliverCommand.RecipientList = append(grpcDeliverCommand.RecipientList, connectionId)

	es.client.Call("message.deliver", grpcDeliverCommand, func() {})

	return nil
}

// broadcast leader to other peers
func (es *ElectionService) BroadcastLeader(peer Peer) error {
	logger.Info(nil, "broadcast leader!")

	updateLeaderMessage := UpdateLeaderMessage{
		Peer: peer,
	}

	grpcDeliverCommand, _ := CreateGrpcDeliverCommand("UpdateLeaderProtocol", updateLeaderMessage)

	pLTable, _ := es.peerQueryService.GetPLTable()

	for _, peer := range pLTable.PeerTable {
		grpcDeliverCommand.RecipientList = append(grpcDeliverCommand.RecipientList, peer.PeerId.Id)
	}

	es.client.Call("message.deliver", grpcDeliverCommand, func() {})

	return nil
}

//broadcast leader when voted fully
func (es *ElectionService) DecideToBeLeader(command command.ReceiveGrpc) error {

	logger.Infof(nil, "current state", es.Election)
	//	1. if candidate, reset left time
	//	2. count up
	if es.Election.GetState() == Candidate {

		es.Election.CountUp()
	}

	//	3. if counted is same with num of peer-1 set leader and publish
	pLTable, _ := es.peerQueryService.GetPLTable()
	numOfPeers := len(pLTable.PeerTable)

	if es.Election.GetVoteCount() == numOfPeers-1 {

		peer := Peer{
			PeerId:    PeerId{Id: ""},
			IpAddress: es.Election.ipAddress,
		}

		es.BroadcastLeader(peer)
	}

	return nil
}

func (es *ElectionService) ElectLeaderWithRaft() {

	//1. Start random timeout
	//2. timed out! alter state to 'candidate'
	//3. while ticking, count down leader repo left time
	//4. Send message having 'RequestVoteProtocol' to other node
	go func() {
		es.Election.state = Ticking

		es.Election.leftTime = GenRandomInRange(150, 300)

		timeout := time.After(time.Duration(es.Election.leftTime) * time.Millisecond)
		tick := time.Tick(1 * time.Millisecond)
		end := true
		for end {
			select {

			case <-timeout:
				logger.Info(nil, "timed out!")
				// when timed out
				// 1. if state is ticking, be candidate and request vote
				// 2. if state is candidate, reset state and left time
				if es.Election.GetState() == Ticking {
					logger.Infof(nil, "candidate process: %v", es.Election.candidate)
					es.Election.SetState(Candidate)

					pLTable, _ := es.peerQueryService.GetPLTable()

					peerList := pLTable.PeerTable

					connectionIds := make([]string, 0)

					for _, peer := range peerList {
						connectionIds = append(connectionIds, peer.PeerId.Id)
					}

					es.RequestVote(connectionIds)

				} else if es.Election.GetState() == Candidate {
					//reset time and state chane candidate -> ticking when timed in candidate state
					es.Election.ResetLeftTime()
					es.Election.SetState(Ticking)
				}

			case <-tick:
				// count down left time while ticking
				es.Election.CountDownLeftTimeBy(1)
			case <-time.After(5 * time.Second):
				end = false
			}

		}
	}()
}

func (es *ElectionService) RequestVote(connectionIds []string) error {

	// 0. be candidate
	es.Election.state = Candidate
	// 1. create request vote message
	// 2. send message
	requestVoteMessage := RequestVoteMessage{}

	grpcDeliverCommand, _ := CreateGrpcDeliverCommand("RequestVoteProtocol", requestVoteMessage)

	for _, connectionId := range connectionIds {

		grpcDeliverCommand.RecipientList = append(grpcDeliverCommand.RecipientList, connectionId)
	}

	es.client.Call("message.deliver", grpcDeliverCommand, func() {})

	return nil
}

func CreateGrpcDeliverCommand(protocol string, body interface{}) (command.DeliverGrpc, error) {

	data, err := common.Serialize(body)

	if err != nil {
		return command.DeliverGrpc{}, err
	}

	return command.DeliverGrpc{
		MessageId:     xid.New().String(),
		RecipientList: make([]string, 0),
		Body:          data,
		Protocol:      protocol,
	}, err
}

func GenRandomInRange(min, max int) int {

	rand.Seed(time.Now().UnixNano())

	return rand.Intn(max-min) + min
}
