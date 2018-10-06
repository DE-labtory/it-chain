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
	"strings"
	"time"

	"github.com/it-chain/engine/common"
	"github.com/it-chain/engine/common/command"
	"github.com/it-chain/engine/common/event"
	"github.com/it-chain/engine/consensus/pbft"
	"github.com/it-chain/iLogger"
	"github.com/rs/xid"
)

type ElectionApi struct {
	ElectionService      *pbft.ElectionService
	parliamentRepository pbft.ParliamentRepository
	eventService         common.EventService
	quit                 chan struct{}
}

func NewElectionApi(electionService *pbft.ElectionService, parliamentRepository pbft.ParliamentRepository, eventService common.EventService) *ElectionApi {

	return &ElectionApi{
		ElectionService:      electionService,
		parliamentRepository: parliamentRepository,
		eventService:         eventService,
		quit:                 make(chan struct{}, 1),
	}
}

func (e *ElectionApi) Vote(connectionId string) error {

	parliament := e.parliamentRepository.Load()

	representative, err := parliament.FindRepresentativeByID(connectionId)
	if err != nil {
		iLogger.Errorf(nil, "[PBFT] Representative who has (Id: %s) is not found", connectionId)
		return err
	}

	e.ElectionService.SetCandidate(representative)
	e.ElectionService.ResetLeftTime()

	voteLeaderMessage := pbft.VoteMessage{}
	grpcDeliverCommand, _ := common.CreateGrpcDeliverCommand("VoteLeaderProtocol", voteLeaderMessage)
	grpcDeliverCommand.RecipientList = append(grpcDeliverCommand.RecipientList, connectionId)

	iLogger.Infof(nil, "[PBFT] Vote to %s", connectionId)
	e.ElectionService.SetVoted(true)

	return e.eventService.Publish("message.deliver", grpcDeliverCommand)
}

// broadcast leader to other peers
func (e *ElectionApi) broadcastLeader(rep pbft.Representative) error {
	iLogger.Infof(nil, "[PBFT] Broadcast leader - ID: [%s]", rep.ID)

	updateLeaderMessage := pbft.UpdateLeaderMessage{
		Representative: rep,
	}
	grpcDeliverCommand, err := common.CreateGrpcDeliverCommand("UpdateLeaderProtocol", updateLeaderMessage)
	if err != nil {
		iLogger.Errorf(nil, "[PBFT] Cannot create grpc command - Error: [%s]", err.Error())
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
		iLogger.Infof(nil, "[PBFT] Leader has fully voted")

		e.EndRaft()
		representative := pbft.Representative{
			ID: e.ElectionService.NodeId,
		}

		parliament := e.parliamentRepository.Load()
		parliament.SetLeader(e.ElectionService.NodeId)
		e.parliamentRepository.Save(parliament)

		e.eventService.Publish("leader.updated", event.LeaderUpdated{
			LeaderId: e.ElectionService.NodeId,
		})

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
	parliament := e.parliamentRepository.Load()
	if !parliament.IsNeedConsensus() {
		e.ElectLeaderWithLargestRepresentativeId()
		return
	}

	e.ElectionService.SetState(pbft.TICKING)
	e.ElectionService.InitLeftTime()

	tick := time.Tick(1 * time.Millisecond)
	timeout := time.After(time.Second * 10)

	for {
		select {
		case <-tick:
			e.ElectionService.CountDownLeftTimeBy(1)
			if e.ElectionService.GetLeftTime() == 0 {
				e.HandleRaftTimeout()
			}
		case <-e.quit:
			iLogger.Infof(nil, "[PBFT] Raft has end")
			return
		case <-timeout:
			iLogger.Errorf(nil, "[PBFT] Raft Time out")
			return
		}
	}
}

func (e *ElectionApi) ElectLeaderWithLargestRepresentativeId() {
	representatives := e.parliamentRepository.Load().GetRepresentatives()

	// TODO: This logic needs to hide into domain
	ids := make([]string, 0)
	for _, rep := range representatives {
		ids = append(ids, rep.ID)
	}

	largestId := common.FindEarliestString(ids)
	e.SetLeader(largestId)
}

func (e *ElectionApi) EndRaft() {
	e.ElectionService.SetState(pbft.NORMAL)
	e.quit <- struct{}{}
}

func (e *ElectionApi) HandleRaftTimeout() error {
	if e.ElectionService.GetState() == pbft.TICKING {
		e.ElectionService.SetState(pbft.CANDIDATE)
		e.ElectionService.ResetLeftTime()
		connectionIds := make([]string, 0)
		parliament := e.parliamentRepository.Load()
		for _, r := range parliament.GetRepresentatives() {
			if r.ID != e.ElectionService.NodeId {
				connectionIds = append(connectionIds, r.ID)

			}
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
	grpcDeliverCommand, _ := common.CreateGrpcDeliverCommand("RequestVoteProtocol", requestVoteMessage)

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

func (e *ElectionApi) GetParliament() pbft.Parliament {
	parliament := e.parliamentRepository.Load()

	return parliament
}

func (e *ElectionApi) SetLeader(representativeId string) {
	parliament := e.parliamentRepository.Load()
	parliament.SetLeader(representativeId)
	e.parliamentRepository.Save(parliament)

	e.eventService.Publish("leader.updated", event.LeaderUpdated{
		LeaderId: representativeId,
	})
}
