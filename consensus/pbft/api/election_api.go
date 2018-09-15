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

package api

import (
	"math/rand"
	"sync"
	"time"

	"github.com/it-chain/engine/common"
	"github.com/it-chain/engine/common/command"
	"github.com/it-chain/engine/common/logger"
	"github.com/it-chain/engine/consensus/pbft"
	"github.com/rs/xid"
)

type ElectionApi struct {
	election     *pbft.ElectionService
	parliament   *pbft.Parliament
	eventService common.EventService
	mux          sync.Mutex
}

func NewElectionApi(election *pbft.ElectionService, parliament *pbft.Parliament, eventService common.EventService) *ElectionApi {

	return &ElectionApi{
		mux:          sync.Mutex{},
		election:     election,
		parliament:   parliament,
		eventService: eventService,
	}
}

func (es *ElectionApi) Vote(connectionId string) error {

	representative := es.parliament.RepresentativeTable[connectionId]

	candidate := es.election.GetCandidate()

	// if peer has no candidate set candidate
	if candidate.ID == "" {

		es.election.SetCandidate(representative)
	} else {

		logger.Info(nil, "[consensus] peer has already received request vote message")
		return nil
	}

	//if leftTime >0, reset left time and send VoteLeaderMessage

	es.election.ResetLeftTime()

	voteLeaderMessage := pbft.VoteMessage{}

	grpcDeliverCommand, _ := CreateGrpcDeliverCommand("VoteLeaderProtocol", voteLeaderMessage)
	grpcDeliverCommand.RecipientList = append(grpcDeliverCommand.RecipientList, connectionId)

	es.eventService.Publish("message.deliver", grpcDeliverCommand)

	return nil
}

// broadcast leader to other peers
func (es *ElectionApi) broadcastLeader(rep pbft.Representative) error {
	logger.Info(nil, "[consensus] Broadcast leader")

	updateLeaderMessage := pbft.UpdateLeaderMessage{
		Representative: rep,
	}

	grpcDeliverCommand, _ := CreateGrpcDeliverCommand("UpdateLeaderProtocol", updateLeaderMessage)

	for _, r := range es.parliament.RepresentativeTable {
		grpcDeliverCommand.RecipientList = append(grpcDeliverCommand.RecipientList, r.ID)
	}

	pubErr := es.eventService.Publish("message.deliver", grpcDeliverCommand)

	if pubErr != nil {
		logger.Infof(nil, "[consensus] Fail to publish update leader message")
		return pubErr
	}

	return nil
}

//broadcast leader when voted fully
func (es *ElectionApi) DecideToBeLeader() error {

	if es.election.GetState() != pbft.CANDIDATE {
		return nil
	}
	//	1. if candidate, reset left time
	//	2. count up
	es.election.CountUpVoteCount()

	//	3. if fully voted set leader and publish

	if es.isFullyVoted() {

		representative := pbft.Representative{
			ID:        "",
			IpAddress: es.election.GetIpAddress(),
		}

		if err := es.broadcastLeader(representative); err != nil {
			return err
		}
	}

	return nil
}

func (ea *ElectionApi) isFullyVoted() bool {
	numOfPeers := len(ea.parliament.RepresentativeTable)
	if ea.election.GetVoteCount() == numOfPeers-1 {
		return true
	}
	return false
}

func (es *ElectionApi) ElectLeaderWithRaft() {

	//1. Start random timeout
	//2. timed out! alter state to 'candidate'
	//3. while ticking, count down leader repo left time
	//4. Send message having 'RequestVoteProtocol' to other node
	go func() {
		es.election.SetState(pbft.TICKING)

		es.election.InitLeftTime()

		timeout := time.After(time.Duration(es.election.GetLeftTime()) * time.Millisecond)
		tick := time.Tick(1 * time.Millisecond)
		end := true
		for end {
			select {

			case <-timeout:
				logger.Info(nil, "[consensus] RAFT timer timed out")
				// when timed out
				// 1. if state is ticking, be candidate and request vote
				// 2. if state is candidate, reset state and left time
				if es.election.GetState() == pbft.TICKING {
					logger.Infof(nil, "[consensus] candidate process: %v", es.election.GetCandidate())
					es.election.SetState(pbft.CANDIDATE)

					connectionIds := make([]string, 0)

					for _, r := range es.parliament.RepresentativeTable {
						connectionIds = append(connectionIds, r.ID)
					}

					es.RequestVote(connectionIds)

				} else if es.election.GetState() == pbft.CANDIDATE {
					//reset time and state chane candidate -> ticking when timed in candidate state
					es.election.ResetLeftTime()
					es.election.SetState(pbft.TICKING)
				}

			case <-tick:
				// count down left time while ticking
				es.election.CountDownLeftTimeBy(1)
			case <-time.After(5 * time.Second):
				end = false
			}

		}
	}()
}

func (es *ElectionApi) RequestVote(connectionIds []string) error {

	// 1. create request vote message
	// 2. send message
	requestVoteMessage := pbft.RequestVoteMessage{}

	grpcDeliverCommand, _ := CreateGrpcDeliverCommand("RequestVoteProtocol", requestVoteMessage)

	for _, connectionId := range connectionIds {

		grpcDeliverCommand.RecipientList = append(grpcDeliverCommand.RecipientList, connectionId)
	}

	es.eventService.Publish("message.deliver", grpcDeliverCommand)

	return nil
}

func (ea *ElectionApi) GetIpAddress() string {
	return ea.election.GetIpAddress()
}

func (ea *ElectionApi) GetCandidate() *pbft.Representative {
	return ea.election.GetCandidate()
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
