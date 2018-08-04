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

func (cApi StateApi) StartConsensus(userId pbft.MemberID, proposedBlock pbft.ProposedBlock) error {

	peerList, _ := cApi.parliamentService.RequestPeerList()

	if !cApi.parliamentService.IsNeedConsensus() {
		return ConsensusCreateError
	}

	createdConsensus, _ := pbft.CreateConsensus(peerList, proposedBlock)
	createdConsensus.Start()
	if err := cApi.repo.Save(*createdConsensus); err != nil {
		return err
	}

	createdPrepareMsg := pbft.NewPrePrepareMsg(createdConsensus)
	cApi.propagateService.BroadcastPrePrepareMsg(*createdPrepareMsg)

	return nil
}

func (cApi StateApi) HandlePrePrepareMsg(msg pbft.PrePrepareMsg) error {

	lid, _ := cApi.parliamentService.RequestLeader()
	if lid.ToString() != msg.SenderID {
		return pbft.InvalidLeaderIdError
	}

	constructedConsensus, _ := pbft.ConstructConsensus(msg)
	constructedConsensus.ToPrepareStage()
	if err := cApi.repo.Save(*constructedConsensus); err != nil {
		return err
	}

	prepareMsg := pbft.NewPrepareMsg(constructedConsensus)
	cApi.propagateService.BroadcastPrepareMsg(*prepareMsg)

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

	if loadedConsensus.CheckPrepareCondition() {
		newCommitMsg := pbft.NewCommitMsg(loadedConsensus)
		loadedConsensus.ToCommitStage()
		cApi.propagateService.BroadcastCommitMsg(*newCommitMsg)
	}

	return nil
}

func (cApi StateApi) HandleCommitMsg(msg pbft.CommitMsg) error {

	loadedConsensus, _ := cApi.repo.Load()

	err := loadedConsensus.SaveCommitMsg(&msg)
	if err != nil {
		return err
	}
	representativeNum := len(loadedConsensus.Representatives)
	commitMsgNum := len(loadedConsensus.PrepareMsgPool.Get())
	satisfyNum := representativeNum / 3

	if commitMsgNum > (satisfyNum + 1) {
		cApi.confirmService.ConfirmBlock(loadedConsensus.Block)

	}
	return nil
}
