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
	propagateService  pbft.PropagateService
	confirmService    pbft.ConfirmService
	parliamentService pbft.ParliamentService
	repo              mem.StateRepository
}

var ConsensusCreateError = errors.New("Consensus can't be created")

func NewStateApi(propagateService pbft.PropagateService,
	confirmService pbft.ConfirmService, parliamentService pbft.ParliamentService, repo mem.StateRepository) StateApi {
	return StateApi{
		propagateService:  propagateService,
		confirmService:    confirmService,
		parliamentService: parliamentService,
		repo:              repo,
	}
}

func (cApi StateApi) StartConsensus(userId pbft.MemberID, proposedBlock pbft.ProposedBlock) error {

	peerList, err := cApi.parliamentService.RequestPeerList()
	if err != nil {
		return err
	}

	if !cApi.parliamentService.IsNeedConsensus() {
		return ConsensusCreateError
	}

	createdConsensus, err := pbft.CreateConsensus(peerList, proposedBlock)
	if err != nil {
		return err
	}

	if err := cApi.repo.Save(*createdConsensus); err != nil {
		return err
	}

	createdPrePrepareMsg := pbft.NewPrePrepareMsg(createdConsensus)
	if err := cApi.propagateService.BroadcastPrePrepareMsg(*createdPrePrepareMsg); err != nil {
		return err
	}
	createdConsensus.Start()

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

	constructedConsensus, err := pbft.ConstructConsensus(msg)
	if err != nil {
		return err
	}

	if err := cApi.repo.Save(*constructedConsensus); err != nil {
		return err
	}

	prepareMsg := pbft.NewPrepareMsg(constructedConsensus)
	if err := cApi.propagateService.BroadcastPrepareMsg(*prepareMsg); err != nil {
		return err
	}
	constructedConsensus.ToPrepareStage()

	return nil
}

func (cApi StateApi) HandlePrepareMsg(msg pbft.PrepareMsg) error {

	loadedConsensus, err := cApi.repo.Load()
	if err != nil {
		return err
	}

	if err := loadedConsensus.SavePrepareMsg(&msg); err != nil {
		return err
	}

	if !loadedConsensus.CheckPrepareCondition() {
		return nil
	}

	newCommitMsg := pbft.NewCommitMsg(loadedConsensus)
	if err := cApi.propagateService.BroadcastCommitMsg(*newCommitMsg); err != nil {
		return err
	}
	loadedConsensus.ToCommitStage()

	return nil
}

func (cApi StateApi) HandleCommitMsg(msg pbft.CommitMsg) error {

	loadedConsensus, err := cApi.repo.Load()

	if err != nil {
		return err
	}

	if err := loadedConsensus.SaveCommitMsg(&msg); err != nil {
		return err
	}

	if !loadedConsensus.CheckCommitCondition() {
		return nil
	}

	if err := cApi.confirmService.ConfirmBlock(loadedConsensus.Block); err != nil {
		return err
	}
	cApi.repo.Remove()

	return nil
}
