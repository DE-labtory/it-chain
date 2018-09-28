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
	"github.com/it-chain/engine/consensus/pbft"
	"github.com/it-chain/iLogger"
	"github.com/rs/xid"
)

type ElectionApi struct {
	ElectionService   *pbft.ElectionService
	parliamentService pbft.ParliamentService
	eventService      common.EventService
	mux               sync.Mutex
}

func NewElectionApi(electionService *pbft.ElectionService, parliamentService pbft.ParliamentService, eventService common.EventService) *ElectionApi {

	return &ElectionApi{
		mux:               sync.Mutex{},
		ElectionService:   electionService,
		parliamentService: parliamentService,
		eventService:      eventService,
	}
}

func (ea *ElectionApi) Vote(connectionId string) error {

	candidate := ea.ElectionService.GetCandidate()
	if candidate.ID != "" {
		iLogger.Info(nil, "[consensus] peer has already received request vote message")
		return nil
	}

	representative := ea.parliamentService.GetRepresentativeById(connectionId)
	ea.ElectionService.SetCandidate(representative)
	ea.ElectionService.ResetLeftTime()

	voteLeaderMessage := pbft.VoteMessage{}
	grpcDeliverCommand, _ := CreateGrpcDeliverCommand("VoteLeaderProtocol", voteLeaderMessage)
	grpcDeliverCommand.RecipientList = append(grpcDeliverCommand.RecipientList, connectionId)

	ea.eventService.Publish("message.deliver", grpcDeliverCommand)

	return nil
}

// broadcast leader to other peers
func (es *ElectionApi) broadcastLeader(rep pbft.Representative) error {
	iLogger.Info(nil, "[Consensus] Broadcast leader")

	updateLeaderMessage := pbft.UpdateLeaderMessage{
		Representative: rep,
	}

	grpcDeliverCommand, _ := CreateGrpcDeliverCommand("UpdateLeaderProtocol", updateLeaderMessage)

	table := es.parliamentService.GetRepresentativeTable()
	for _, r := range table {
		grpcDeliverCommand.RecipientList = append(grpcDeliverCommand.RecipientList, r.ID)
	}

	if err := es.eventService.Publish("message.deliver", grpcDeliverCommand); err != nil {
		iLogger.Infof(nil, "[Consensus] Fail to publish update leader message")
		return err
	}

	return nil
}

//broadcast leader when voted fully
func (es *ElectionApi) DecideToBeLeader() error {
	if es.ElectionService.GetState() != pbft.CANDIDATE {
		return nil
	}

	es.ElectionService.CountUpVoteCount()

	if es.isFullyVoted() {
		representative := pbft.Representative{
			ID: es.ElectionService.NodeId,
		}

		if err := es.broadcastLeader(representative); err != nil {
			return err
		}
	}

	return nil
}

func (ea *ElectionApi) isFullyVoted() bool {
	numOfPeers := len(ea.parliamentService.GetParliament().RepresentativeTable)
	if ea.ElectionService.GetVoteCount() == numOfPeers-1 {
		return true
	}
	return false
}

//1. Start random timeout
//2. timed out! alter state to 'candidate'
//3. while ticking, count down leader repo left time
//4. Send message having 'RequestVoteProtocol' to other node
func (e *ElectionApi) ElectLeaderWithRaft() {

	go func() {
		e.ElectionService.SetState(pbft.TICKING)
		e.ElectionService.InitLeftTime()

		timeout := time.After(time.Duration(e.ElectionService.GetLeftTime()) * time.Millisecond)
		tick := time.Tick(1 * time.Millisecond)
		end := true
		for end {
			select {

			case <-timeout:
				if e.ElectionService.GetState() == pbft.TICKING {

					e.ElectionService.SetState(pbft.CANDIDATE)

					connectionIds := make([]string, 0)
					repTable := e.parliamentService.GetRepresentativeTable()
					for _, r := range repTable {
						connectionIds = append(connectionIds, r.ID)
					}

					e.RequestVote(connectionIds)

				} else if e.ElectionService.GetState() == pbft.CANDIDATE {
					//reset time and state chane candidate -> ticking when timed in candidate state
					e.ElectionService.ResetLeftTime()
					e.ElectionService.SetState(pbft.TICKING)
				}

			case <-tick:
				// count down left time while ticking
				e.ElectionService.CountDownLeftTimeBy(1)
			case <-time.After(5 * time.Second):
				end = false
			}
		}
	}()
}

func (e *ElectionApi) RequestVote(connectionIds []string) error {
	// 1. create request vote message
	// 2. send message
	requestVoteMessage := pbft.RequestVoteMessage{}
	grpcDeliverCommand, _ := CreateGrpcDeliverCommand("RequestVoteProtocol", requestVoteMessage)

	for _, connectionId := range connectionIds {
		grpcDeliverCommand.RecipientList = append(grpcDeliverCommand.RecipientList, connectionId)
	}

	return e.eventService.Publish("message.deliver", grpcDeliverCommand)
}

func (ea *ElectionApi) GetCandidate() *pbft.Representative {
	return ea.ElectionService.GetCandidate()
}

func (ea *ElectionApi) GetState() pbft.ElectionState {
	return ea.ElectionService.GetState()
}

func (ea *ElectionApi) SetState(state pbft.ElectionState) {
	ea.ElectionService.SetState(state)
}

func (ea *ElectionApi) GetVoteCount() int {
	return ea.ElectionService.GetVoteCount()
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
