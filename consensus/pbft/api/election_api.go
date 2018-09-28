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
	"time"

	"fmt"

	"strings"

	"github.com/it-chain/engine/common"
	"github.com/it-chain/engine/common/command"
	"github.com/it-chain/engine/consensus/pbft"
	"github.com/it-chain/iLogger"
	"github.com/rs/xid"
)

type ElectionApi struct {
	ElectionService      *pbft.ElectionService
	parliamentRepository pbft.ParliamentRepository
	eventService         common.EventService
}

func NewElectionApi(electionService *pbft.ElectionService, parliamentRepository pbft.ParliamentRepository, eventService common.EventService) *ElectionApi {

	return &ElectionApi{
		ElectionService:      electionService,
		parliamentRepository: parliamentRepository,
		eventService:         eventService,
	}
}

func (e *ElectionApi) Vote(connectionId string) error {

	e.ElectionService.IncreaseTerm()

	parliament := e.parliamentRepository.Load()

	representative, err := parliament.FindRepresentativeByID(connectionId)
	if err != nil {
		return err
	}

	e.ElectionService.SetCandidate(representative)
	e.ElectionService.ResetLeftTime()

	voteLeaderMessage := pbft.VoteMessage{}
	grpcDeliverCommand, _ := CreateGrpcDeliverCommand("VoteLeaderProtocol", voteLeaderMessage)
	grpcDeliverCommand.RecipientList = append(grpcDeliverCommand.RecipientList, connectionId)
	return e.eventService.Publish("message.deliver", grpcDeliverCommand)
}

// broadcast leader to other peers
func (e *ElectionApi) broadcastLeader(rep pbft.Representative) error {
	iLogger.Infof(nil, "[PBFT] Broadcast leader id: %s", rep.ID)

	updateLeaderMessage := pbft.UpdateLeaderMessage{
		Representative: rep,
	}
	grpcDeliverCommand, err := CreateGrpcDeliverCommand("UpdateLeaderProtocol", updateLeaderMessage)
	if err != nil {
		return err
	}

	parliament := e.parliamentRepository.Load()
	for _, r := range parliament.GetRepresentatives() {
		grpcDeliverCommand.RecipientList = append(grpcDeliverCommand.RecipientList, r.ID)
	}

	return e.eventService.Publish("message.deliver", grpcDeliverCommand)
}

//broadcast leader when voted fully
func (e *ElectionApi) DecideToBeLeader() error {
	if e.ElectionService.GetState() != pbft.CANDIDATE {
		return nil
	}

	e.ElectionService.CountUpVoteCount()

	if e.isFullyVoted() {
		representative := pbft.Representative{
			ID: e.ElectionService.NodeId,
		}

		if err := e.broadcastLeader(representative); err != nil {
			return err
		}
	}

	return nil
}

func (e *ElectionApi) isFullyVoted() bool {
	parliament := e.parliamentRepository.Load()
	numOfPeers := len(parliament.Representatives)
	if e.ElectionService.GetVoteCount() == numOfPeers-1 {
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

		fmt.Println(e.ElectionService.GetLeftTime())

		tick := time.Tick(1 * time.Millisecond)
		end := true
		for end {
			select {

			case <-tick:
				// count down left time while ticking
				e.ElectionService.CountDownLeftTimeBy(1)
				if e.ElectionService.GetLeftTime() == 0 {
					e.HandleRaftTimeout()
				}
			case <-time.After(4 * time.Second):
				end = false
				fmt.Println("end raft")
			}
		}
	}()
}
func (e *ElectionApi) HandleRaftTimeout() error {
	fmt.Println("start raft")
	if e.ElectionService.GetState() == pbft.TICKING {

		e.ElectionService.SetState(pbft.CANDIDATE)
		connectionIds := make([]string, 0)
		parliament := e.parliamentRepository.Load()
		for _, r := range parliament.GetRepresentatives() {
			connectionIds = append(connectionIds, r.ID)
		}
		e.RequestVote(connectionIds)
	} else if e.ElectionService.GetState() == pbft.CANDIDATE {
		//reset time and state chane candidate -> ticking when timed in candidate state
		e.ElectionService.ResetLeftTime()
		e.ElectionService.SetState(pbft.TICKING)
	}
	return nil
}
func (e *ElectionApi) RequestVote(peerIds []string) error {
	iLogger.Infof(nil, "[PBFT] Request Vote - Peers:[%s]", strings.Join(peerIds, ", "))
	// 1. create request vote message
	// 2. send message
	requestVoteMessage := pbft.RequestVoteMessage{
		Term: e.ElectionService.GetTerm(),
	}
	grpcDeliverCommand, _ := CreateGrpcDeliverCommand("RequestVoteProtocol", requestVoteMessage)

	for _, connectionId := range peerIds {
		grpcDeliverCommand.RecipientList = append(grpcDeliverCommand.RecipientList, connectionId)
	}

	return e.eventService.Publish("message.deliver", grpcDeliverCommand)
}

func (e *ElectionApi) GetCandidate() pbft.Representative {
	return e.ElectionService.GetCandidate()
}

func (e *ElectionApi) GetState() pbft.ElectionState {
	return e.ElectionService.GetState()
}

func (e *ElectionApi) SetState(state pbft.ElectionState) {
	e.ElectionService.SetState(state)
}

func (e *ElectionApi) GetVoteCount() int {
	return e.ElectionService.GetVoteCount()
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
