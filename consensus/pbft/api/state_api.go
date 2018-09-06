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
	"github.com/it-chain/engine/consensus/pbft"
)

type StateApi interface {
	StartConsensus(block pbft.ProposedBlock) error
	HandlePrePrepareMsg(msg pbft.PrePrepareMsg) error
	HandlePrepareMsg(msg pbft.PrepareMsg) error
	HandleCommitMsg(msg pbft.CommitMsg) error
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

	createdPrePrepareMsg := pbft.NewPrePrepareMsg(createdState, cApi.publisherID)
	if err := cApi.propagateService.BroadcastPrePrepareMsg(*createdPrePrepareMsg); err != nil {
		return err
	}

	createdState.Start()
	if err := cApi.repo.Save(*createdState); err != nil {
		return err
	}

	return nil
}

func (cApi *StateApiImpl) HandlePrePrepareMsg(msg pbft.PrePrepareMsg) error {

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

	prepareMsg := pbft.NewPrepareMsg(builtState, cApi.publisherID)
	if err := cApi.propagateService.BroadcastPrepareMsg(*prepareMsg); err != nil {
		return err
	}
	builtState.ToPrepareStage()

	if err := cApi.repo.Save(*builtState); err != nil {
		return err
	}

	return nil
}

func (cApi *StateApiImpl) HandlePrepareMsg(msg pbft.PrepareMsg) error {

	loadedState, err := cApi.repo.Load()
	if err != nil {
		return err
	}

	if err := loadedState.SavePrepareMsg(&msg); err != nil {
		return err
	}

	if !loadedState.CheckPrepareCondition() {
		return nil
	}

	newCommitMsg := pbft.NewCommitMsg(&loadedState, cApi.publisherID)
	if err := cApi.propagateService.BroadcastCommitMsg(*newCommitMsg); err != nil {
		return err
	}
	loadedState.ToCommitStage()

	if err := cApi.repo.Save(loadedState); err != nil {
		return err
	}

	return nil
}

func (cApi *StateApiImpl) HandleCommitMsg(msg pbft.CommitMsg) error {

	loadedState, err := cApi.repo.Load()
	if err != nil {
		return err
	}

	if err := loadedState.SaveCommitMsg(&msg); err != nil {
		return err
	}

	if !loadedState.CheckCommitCondition() {
		return nil
	}

	if err := cApi.eventService.Publish("block.confirm", loadedState.Block); err != nil {
		return err
	}
	cApi.repo.Remove()

	return nil
}
