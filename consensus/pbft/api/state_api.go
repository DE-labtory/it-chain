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
	propagateService     *pbft.PropagateService
	eventService         common.EventService
	parliamentRepository pbft.ParliamentRepository
	stateRepository      pbft.StateRepository
	tempPrevoteMsgPool   pbft.PrevoteMsgPool
	tempPreCommitMsgPool pbft.PreCommitMsgPool
}

var ConsensusCreateError = errors.New("Consensus can't be created")

func NewStateApi(publisherID string, propagateService *pbft.PropagateService,
	eventService common.EventService, parliamentRepository pbft.ParliamentRepository, repo pbft.StateRepository) *StateApi {
	return &StateApi{
		publisherID:          publisherID,
		propagateService:     propagateService,
		eventService:         eventService,
		parliamentRepository: parliamentRepository,
		stateRepository:      repo,
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

	receipients := createdState.GetReceipients(sApi.publisherID)

	iLogger.Infof(nil, "[PBFT] Leader broadcasts ProposeMsg to %v", receipients)
	if err := sApi.propagateService.BroadcastProposeMsg(*createdProposeMsg, receipients); err != nil {
		return err
	}

	createdState.Start()
	iLogger.Infof(nil, "[PBFT] Consensus starts - Stage: [%s]", createdState.CurrentStage)

	if err := sApi.stateRepository.Save(*createdState); err != nil {
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

	receipients := builtState.GetReceipients(sApi.publisherID)

	iLogger.Debugf(nil, "[PBFT] Representative broadcasts PreVoteMsg to %v", receipients)
	prevoteMsg := pbft.NewPrevoteMsg(builtState, sApi.publisherID)
	if err := sApi.propagateService.BroadcastPrevoteMsg(*prevoteMsg, receipients); err != nil {
		return err
	}

	builtState.ToPrevoteStage()
	logger.Infof(nil, "[PBFT] Prevoted - Stage: [%s]", builtState.CurrentStage)

	if err := sApi.stateRepository.Save(*builtState); err != nil {
		return err
	}

	return nil
}

func (sApi *StateApi) HandlePrevoteMsg(msg pbft.PrevoteMsg) (returnErr error) {

	loadedState, err := sApi.stateRepository.Load()
	if err != nil {
		sApi.tempPrevoteMsgPool.Save(&msg)
		iLogger.Debugf(nil, "[PBFT] %s while handling PreVote message", err)
		return err
	}

	receipients := loadedState.GetReceipients(sApi.publisherID)

	tempPrevoteMsgPool := sApi.tempPrevoteMsgPool.Get()
	for i := 0; i < len(tempPrevoteMsgPool); i++ {
		loadedState.PrevoteMsgPool.Save(&tempPrevoteMsgPool[i])
	}
	sApi.tempPrevoteMsgPool.RemoveAllMsgs()

	if err := loadedState.SavePrevoteMsg(&msg); err != nil {
		return err
	}

	defer func() {
		if err := sApi.stateRepository.Save(loadedState); err != nil {
			returnErr = err
		}
	}()

	if loadedState.CheckPrevoteCondition() {
		if loadedState.CurrentStage == pbft.PREVOTE_STAGE {
			iLogger.Infof(nil, "[PBFT] Representative broadcasts PreCommitMsg to %v", receipients)
			newCommitMsg := pbft.NewPreCommitMsg(&loadedState, sApi.publisherID)
			if err := sApi.propagateService.BroadcastPreCommitMsg(*newCommitMsg, receipients); err != nil {
				return err
			}

			loadedState.ToPreCommitStage()
			iLogger.Infof(nil, "[PBFT] PreCommitted - Stage: [%s]", loadedState.CurrentStage)
		}
	}

	return returnErr
}

func (sApi *StateApi) HandlePreCommitMsg(msg pbft.PreCommitMsg) error {

	loadedState, err := sApi.stateRepository.Load()
	if err != nil {
		sApi.tempPreCommitMsgPool.Save(&msg)
		iLogger.Debug(nil, "[PBFT] State repository is empty when handling PreCommit message")
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
		e := event.ConsensusFinished{
			Seal: loadedState.Block.Seal,
			Body: loadedState.Block.Body,
		}

		if err := sApi.eventService.Publish("block.confirm", e); err != nil {
			return err
		}
		iLogger.Debug(nil, "[PBFT] Published block confirm event")

		sApi.stateRepository.Remove()
		logger.Infof(nil, "[PBFT] Consensus is finished.")
		return nil
	}

	if err := sApi.stateRepository.Save(loadedState); err != nil {
		return err
	}

	return nil
}
