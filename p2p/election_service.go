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
	"github.com/it-chain/iLogger"
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

	candidate := es.Election.GetCandidate()

	// if peer has already received request vote message return nil
	if candidate.PeerId.Id == "" {

		es.Election.SetCandidate(&peer)
	} else {
		return nil
	}

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
	iLogger.Info(nil, "[P2P] Broadcast leader")

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
func (es *ElectionService) DecideToBeLeader() error {

	//	1. if candidate, reset left time
	//	2. count up
	if es.Election.GetState() == CANDIDATE {

		es.Election.CountUpVoteCount()
	}

	//	3. if counted is same with num of peer-1 set leader and publish
	pLTable, _ := es.peerQueryService.GetPLTable()
	numOfPeers := len(pLTable.PeerTable)

	if es.Election.GetVoteCount() == numOfPeers-1 {

		peer := Peer{
			PeerId:    PeerId{Id: ""},
			IpAddress: es.Election.ipAddress,
		}

		if err := es.BroadcastLeader(peer); err != nil {
			return err
		}
	}

	return nil
}

func (es *ElectionService) RequestVote(connectionIds []string) error {

	// 0. be candidate
	es.Election.state = CANDIDATE
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
