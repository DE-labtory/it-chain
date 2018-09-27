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

type StateApi struct {
	publisherID          string
	propagateService     pbft.PropagateService
	eventService         common.EventService
	parliamentRepository pbft.ParliamentRepository
	repo                 pbft.StateRepository
	tempPrevoteMsgPool   pbft.PrevoteMsgPool
	tempPreCommitMsgPool pbft.PreCommitMsgPool
}

var ConsensusCreateError = errors.New("Consensus can't be created")

func NewStateApi(publisherID string, propagateService pbft.PropagateService,
	eventService common.EventService, parliamentRepository pbft.ParliamentRepository, repo pbft.StateRepository) StateApi {
	return StateApi{
		publisherID:          publisherID,
		propagateService:     propagateService,
		eventService:         eventService,
		parliamentRepository: parliamentRepository,
		repo:                 repo,
		tempPrevoteMsgPool:   pbft.NewPrevoteMsgPool(),
		tempPreCommitMsgPool: pbft.NewPreCommitMsgPool(),
	}
}

func (sApi *StateApi) StartConsensus(proposedBlock pbft.ProposedBlock) error {

	parliament := sApi.parliamentRepository.Load()
	if !parliament.IsNeedConsensus() {
		return ConsensusCreateError
	}

	createdState, err := pbft.NewState(parliament.GetRepresentatives(), proposedBlock)
	if err != nil {
		return err
	}

	createdProposeMsg := pbft.NewProposeMsg(createdState, sApi.publisherID)
	if err := sApi.propagateService.BroadcastProposeMsg(*createdProposeMsg, createdState.Representatives); err != nil {
		return err
	}

	createdState.Start()
	if err := sApi.repo.Save(*createdState); err != nil {
		return err
	}

	return nil
}

func (sApi *StateApi) HandleProposeMsg(msg pbft.ProposeMsg) error {

	parliament := sApi.parliamentRepository.Load()

	lid := parliament.GetLeader()
	if lid.GetID() != msg.SenderID {
		return pbft.InvalidLeaderIdError
	}

	builtState, err := pbft.BuildState(msg)
	if err != nil {
		return err
	}

	prevoteMsg := pbft.NewPrevoteMsg(builtState, sApi.publisherID)
	if err := sApi.propagateService.BroadcastPrevoteMsg(*prevoteMsg, builtState.Representatives); err != nil {
		return err
	}
	builtState.ToPrevoteStage()

	if err := sApi.repo.Save(*builtState); err != nil {
		return err
	}

	return nil
}

func (sApi *StateApi) HandlePrevoteMsg(msg pbft.PrevoteMsg) (returnErr error) {

	loadedState, err := sApi.repo.Load()
	if err != nil {
		sApi.tempPrevoteMsgPool.Save(&msg)
		return err
	}

	for i := 0; i < len(sApi.tempPrevoteMsgPool.Get()); i++ {
		loadedState.PrevoteMsgPool.Save(&sApi.tempPrevoteMsgPool.Get()[i])
	}
	sApi.tempPrevoteMsgPool.RemoveAllMsgs()

	if err := loadedState.SavePrevoteMsg(&msg); err != nil {
		return err
	}

	defer func() {
		if err := sApi.repo.Save(loadedState); err != nil {
			returnErr = err
		}
	}()

	if !loadedState.CheckPrevoteCondition() {
		return returnErr
	}

	newCommitMsg := pbft.NewPreCommitMsg(&loadedState, sApi.publisherID)
	if err := sApi.propagateService.BroadcastPreCommitMsg(*newCommitMsg, loadedState.Representatives); err != nil {
		return err
	}
	loadedState.ToPreCommitStage()

	return returnErr
}

func (sApi *StateApi) HandlePreCommitMsg(msg pbft.PreCommitMsg) error {

	loadedState, err := sApi.repo.Load()
	if err != nil {
		sApi.tempPreCommitMsgPool.Save(&msg)
		return err
	}

	for i := 0; i < len(sApi.tempPreCommitMsgPool.Get()); i++ {
		loadedState.PreCommitMsgPool.Save(&sApi.tempPreCommitMsgPool.Get()[i])
	}
	sApi.tempPreCommitMsgPool.RemoveAllMsgs()

	if err := loadedState.SavePreCommitMsg(&msg); err != nil {
		return err
	}

	if !loadedState.CheckPreCommitCondition() {
		if err := sApi.repo.Save(loadedState); err != nil {
			return err
		}
		return nil
	}
	//TODO ConsensusFinished 인자 추가
	if err := sApi.eventService.Publish("block.confirm", event.ConsensusFinished{}); err != nil {
		return err
	}
	sApi.repo.Remove()

	return nil
}
