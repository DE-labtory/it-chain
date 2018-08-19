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

	"github.com/it-chain/engine/consensus/pbft"
	"github.com/it-chain/engine/consensus/pbft/infra/mem"
)

type StateApi struct {
	publisherID       string
	propagateService  pbft.PropagateService
	eventService      pbft.EventService
	parliamentService pbft.ParliamentService
	repo              mem.StateRepository
}

var ConsensusCreateError = errors.New("Consensus can't be created")

func NewStateApi(publisherID string, propagateService pbft.PropagateService,
	confirmService pbft.EventService, parliamentService pbft.ParliamentService, repo mem.StateRepository) StateApi {
	return StateApi{
		publisherID:       publisherID,
		propagateService:  propagateService,
		eventService:      confirmService,
		parliamentService: parliamentService,
		repo:              repo,
	}
}

func (cApi StateApi) StartConsensus(proposedBlock pbft.ProposedBlock) error {

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

	if err := cApi.repo.Save(*createdState); err != nil {
		return err
	}

	createdPrePrepareMsg := pbft.NewPrePrepareMsg(createdState, cApi.publisherID)
	if err := cApi.propagateService.BroadcastPrePrepareMsg(*createdPrePrepareMsg); err != nil {
		return err
	}
	createdState.Start()

	return nil
}

func (cApi StateApi) HandlePrePrepareMsg(msg pbft.PrePrepareMsg) error {

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

	if err := cApi.repo.Save(*builtState); err != nil {
		return err
	}

	prepareMsg := pbft.NewPrepareMsg(builtState, cApi.publisherID)
	if err := cApi.propagateService.BroadcastPrepareMsg(*prepareMsg); err != nil {
		return err
	}
	builtState.ToPrepareStage()

	return nil
}

func (cApi StateApi) HandlePrepareMsg(msg pbft.PrepareMsg) error {

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

	newCommitMsg := pbft.NewCommitMsg(loadedState, cApi.publisherID)
	if err := cApi.propagateService.BroadcastCommitMsg(*newCommitMsg); err != nil {
		return err
	}
	loadedState.ToCommitStage()

	return nil
}

func (cApi StateApi) HandleCommitMsg(msg pbft.CommitMsg) error {

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

	if err := cApi.eventService.ConfirmBlock(loadedState.Block); err != nil {
		return err
	}
	cApi.repo.Remove()

	return nil
}
