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
	"github.com/it-chain/engine/common/logger"
	"github.com/it-chain/engine/consensus/pbft"
	"github.com/it-chain/iLogger"
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
	eventService common.EventService, parliamentRepository pbft.ParliamentRepository, repo pbft.StateRepository) *StateApi {
	return &StateApi{
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

	receipients := make([]pbft.Representative, 0)

	for _, rep := range createdState.Representatives {
		if rep.ID != sApi.publisherID {
			receipients = append(receipients, rep)
		}
	}

	iLogger.Infof(nil, "[PBFT] Broadcast ProposeMsg to %v", receipients)

	if err := sApi.propagateService.BroadcastProposeMsg(*createdProposeMsg, receipients); err != nil {
		return err
	}
	logger.Infof(nil, "[PBFT] Leader broadcast ProposeMsg")
	createdState.Start()
	logger.Infof(nil, "[PBFT] Change stage to propose")
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

	builtState := pbft.BuildState(msg)

	receipients := make([]pbft.Representative, 0)
	for _, rep := range builtState.Representatives {
		if rep.ID != sApi.publisherID {
			receipients = append(receipients, rep)
		}
	}

	prevoteMsg := pbft.NewPrevoteMsg(builtState, sApi.publisherID)
	if err := sApi.propagateService.BroadcastPrevoteMsg(*prevoteMsg, receipients); err != nil {
		return err
	}
	builtState.ToPrevoteStage()
	logger.Infof(nil, "[PBFT] Change stage to Prevote")

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

	receipients := make([]pbft.Representative, 0)
	for _, rep := range loadedState.Representatives {
		if rep.ID != sApi.publisherID {
			receipients = append(receipients, rep)
		}
	}

	tempPrevoteMsgPool := sApi.tempPrevoteMsgPool.Get()
	for i := 0; i < len(tempPrevoteMsgPool); i++ {
		loadedState.PrevoteMsgPool.Save(&tempPrevoteMsgPool[i])
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

	if loadedState.CheckPrevoteCondition() {
		if loadedState.CurrentStage == pbft.PREVOTE_STAGE {
			newCommitMsg := pbft.NewPreCommitMsg(&loadedState, sApi.publisherID)
			if err := sApi.propagateService.BroadcastPreCommitMsg(*newCommitMsg, receipients); err != nil {
				return err
			}
			logger.Infof(nil, "[PBFT] Leader broadcast PreCommitMsg")
			loadedState.ToPreCommitStage()
			logger.Infof(nil, "[PBFT] Change stage to PreCommitStage")
		}
	}

	return returnErr
}

func (sApi *StateApi) HandlePreCommitMsg(msg pbft.PreCommitMsg) error {

	loadedState, err := sApi.repo.Load()
	if err != nil {
		sApi.tempPreCommitMsgPool.Save(&msg)
		return err
	}

	tempPreCommitMsgPool := sApi.tempPreCommitMsgPool.Get()
	for i := 0; i < len(tempPreCommitMsgPool); i++ {
		loadedState.PreCommitMsgPool.Save(&tempPreCommitMsgPool[i])
	}
	sApi.tempPreCommitMsgPool.RemoveAllMsgs()

	if err := loadedState.SavePreCommitMsg(&msg); err != nil {
		return err
	}

	if loadedState.CheckPreCommitCondition() {
		if err := sApi.eventService.Publish("block.confirm", event.ConsensusFinished{
			Seal: loadedState.Block.Seal,
			Body: loadedState.Block.Body,
		}); err != nil {
			return err
		}
		sApi.repo.Remove()
		logger.Infof(nil, "[PBFT] Consensus is finished.")
	}

	return nil
}
