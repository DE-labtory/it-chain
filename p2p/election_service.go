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
	"github.com/it-chain/engine/conf"
	"github.com/it-chain/midgard"
	"github.com/rs/xid"
)

type Publish func(exchange string, topic string, data interface{}) (err error) // 나중에 의존성 주입을 해준다.

type ElectionService struct {
	mux                sync.Mutex
	electionRepository ElectionRepository
	peerService        IPeerService
	peerQueryService   PeerQueryService
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

				pLTable, _ := es.peerQueryService.GetPLTable()

				peerList := pLTable.PeerTable

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

		grpcDeliverCommand.RecipientList = append(grpcDeliverCommand.RecipientList, connectionId)
	}

	es.publish("Command", "message.send", grpcDeliverCommand)

	return nil
}

func CreateGrpcDeliverCommand(protocol string, body interface{}) (command.DeliverGrpc, error) {

	data, err := common.Serialize(body)

	if err != nil {
		return command.DeliverGrpc{}, err
	}

	return command.DeliverGrpc{
		CommandModel: midgard.CommandModel{
			ID: xid.New().String(),
		},
		RecipientList: make([]string, 0),
		Body:          data,
		Protocol:      protocol,
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
	grpcDeliverCommand.RecipientList = append(grpcDeliverCommand.RecipientList, connectionId)

	es.publish("Command", "message.deliver", grpcDeliverCommand)

	return nil
}

func (es *ElectionService) BroadcastLeader(peer Peer) error {

	updateLeaderMessage := UpdateLeaderMessage{}

	grpcDeliverCommand, _ := CreateGrpcDeliverCommand("UpdateLeaderProtocol", updateLeaderMessage)

	pLTable, _ := es.peerQueryService.GetPLTable()

	for _, peer := range pLTable.PeerTable {
		grpcDeliverCommand.RecipientList = append(grpcDeliverCommand.RecipientList, peer.PeerId.Id)
	}

	es.publish("Command", "message.deliver", grpcDeliverCommand)

	return nil
}

//broad case leader when voted fully
func (es *ElectionService) DecideToBeLeader(command command.ReceiveGrpc) error {
	election := es.electionRepository.GetElection()

	//	1. if candidate, reset left time
	//	2. count up
	if election.GetState() == "candidate" {

		election.CountUp()
		es.electionRepository.SetElection(election)
	}

	//	3. if counted is same with num of peer-1 set leader and publish
	pLTable, _ := es.peerQueryService.GetPLTable()
	numOfPeers := len(pLTable.PeerTable)

	if election.GetVoteCount() == numOfPeers-1 {

		peer := Peer{
			PeerId:    PeerId{Id: ""},
			IpAddress: conf.GetConfiguration().GrpcGateway.Address + ":" + conf.GetConfiguration().GrpcGateway.Port,
		}

		es.BroadcastLeader(peer)
	}

	return nil
}
