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
	"errors"

	"github.com/it-chain/engine/common"
	"github.com/it-chain/engine/common/event"
	"github.com/it-chain/engine/consensus/pbft"
)

type StateApi interface {
	StartConsensus(block pbft.ProposedBlock) error
	HandleProposeMsg(msg pbft.ProposeMsg) error
	HandlePrevoteMsg(msg pbft.PrevoteMsg) error
	HandlePreCommitMsg(msg pbft.PreCommitMsg) error
}

type StateApiImpl struct {
	publisherID       string
	propagateService  pbft.PropagateService
	eventService      common.EventService
	parliamentService pbft.ParliamentService
	repo              *pbft.StateRepository
}

var ConsensusCreateError = errors.New("Consensus can't be created")

func NewStateApi(publisherID string, propagateService pbft.PropagateService,
	eventService common.EventService, parliamentService pbft.ParliamentService, repo *pbft.StateRepository) StateApiImpl {
	return StateApiImpl{
		publisherID:       publisherID,
		propagateService:  propagateService,
		eventService:      eventService,
		parliamentService: parliamentService,
		repo:              repo,
	}
}

func (cApi *StateApiImpl) StartConsensus(proposedBlock pbft.ProposedBlock) error {

	peerList, err := cApi.parliamentService.RequestPeerList()
	if err != nil {
		return err
	}

	if !cApi.parliamentService.IsNeedConsensus() {
		return ConsensusCreateError
	}

	createdState, err := pbft.NewState(peerList, proposedBlock)
	if err != nil {
		return err
	}

	createdProposeMsg := pbft.NewProposeMsg(createdState, cApi.publisherID)
	if err := cApi.propagateService.BroadcastProposeMsg(*createdProposeMsg, createdState.Representatives); err != nil {
		return err
	}

	createdState.Start()
	if err := cApi.repo.Save(*createdState); err != nil {
		return err
	}

	return nil
}

func (cApi *StateApiImpl) HandleProposeMsg(msg pbft.ProposeMsg) error {

	lid, err := cApi.parliamentService.RequestLeader()
	if err != nil {
		return err
	}

	if lid.ToString() != msg.SenderID {
		return pbft.InvalidLeaderIdError
	}

	builtState, err := pbft.BuildState(msg)
	if err != nil {
		return err
	}

	prevoteMsg := pbft.NewPrevoteMsg(builtState, cApi.publisherID)
	if err := cApi.propagateService.BroadcastPrevoteMsg(*prevoteMsg, builtState.Representatives); err != nil {
		return err
	}
	builtState.ToPrevoteStage()

	if err := cApi.repo.Save(*builtState); err != nil {
		return err
	}

	return nil
}

func (cApi *StateApiImpl) HandlePrevoteMsg(msg pbft.PrevoteMsg) error {

	loadedState, err := cApi.repo.Load()
	if err != nil {
		return err
	}

	if err := loadedState.SavePrevoteMsg(&msg); err != nil {
		return err
	}

	if !loadedState.CheckPrevoteCondition() {
		return nil
	}

	newCommitMsg := pbft.NewPreCommitMsg(&loadedState, cApi.publisherID)
	if err := cApi.propagateService.BroadcastPreCommitMsg(*newCommitMsg, loadedState.Representatives); err != nil {
		return err
	}
	loadedState.ToPreCommitStage()

	if err := cApi.repo.Save(loadedState); err != nil {
		return err
	}

	return nil
}

func (cApi *StateApiImpl) HandlePreCommitMsg(msg pbft.PreCommitMsg) error {

	loadedState, err := cApi.repo.Load()
	if err != nil {
		return err
	}

	if err := loadedState.SavePreCommitMsg(&msg); err != nil {
		return err
	}

	if !loadedState.CheckPreCommitCondition() {
		return nil
	}

	//TODO ConsensusFinished 인자 추가
	if err := cApi.eventService.Publish("block.confirm", event.ConsensusFinished{}); err != nil {
		return err
	}
	cApi.repo.Remove()

	return nil
}
