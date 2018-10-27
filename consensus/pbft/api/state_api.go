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
	"github.com/it-chain/iLogger"
	"github.com/it-chain/sdk/logger"
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

func (s *StateApi) StartConsensus(proposedBlock pbft.ProposedBlock) error {

	parliament := s.parliamentRepository.Load()
	if !parliament.IsNeedConsensus() {
		return ConsensusCreateError
	}

	createdState, err := pbft.NewState(parliament.GetRepresentatives(), proposedBlock)
	if err != nil {
		return err
	}

	createdProposeMsg := pbft.NewProposeMsg(createdState, s.publisherID)

	receipients := createdState.GetReceipients(s.publisherID)

	iLogger.Infof(nil, "[PBFT] Leader broadcasts ProposeMsg to %v", receipients)
	if err := s.propagateService.BroadcastProposeMsg(*createdProposeMsg, receipients); err != nil {
		return err
	}

	createdState.Start()
	iLogger.Infof(nil, "[PBFT] Consensus starts - Stage: [%s]", createdState.CurrentStage)

	if err := s.stateRepository.Save(*createdState); err != nil {
		return err
	}

	return nil
}

func (s *StateApi) AcceptProposal(proposal pbft.ProposeMsg) (returnErr error) {

	parliament := s.parliamentRepository.Load()

	lid := parliament.GetLeader()
	if lid.GetID() != proposal.SenderID {
		return pbft.InvalidLeaderIdError
	}

	builtState := pbft.BuildState(proposal)

	receipients := builtState.GetReceipients(s.publisherID)

	iLogger.Debugf(nil, "[PBFT] Representative broadcasts PreVoteMsg to %v", receipients)
	prevoteMsg := pbft.NewPrevoteMsg(builtState, s.publisherID)
	if err := s.propagateService.BroadcastPrevoteMsg(*prevoteMsg, receipients); err != nil {
		return err
	}

	builtState.ToPrevoteStage()
	iLogger.Infof(nil, "[PBFT] Prevoted - Stage: [%s]", builtState.CurrentStage)

	defer func() {
		if err := s.stateRepository.Save(*builtState); err != nil {
			returnErr = err
		}
	}()

	return returnErr
}

func (s *StateApi) ReceivePrevote(msg pbft.PrevoteMsg) (returnErr error) {

	if err := s.tempPrevoteMsgPool.Save(&msg); err != nil {
		return err
	}

	loadedState, err := s.stateRepository.Load()
	if err != nil {
		iLogger.Debugf(nil, "[PBFT] Cannot receive Prevote message - Error: [%s]", err)
		return nil
	}

	// todo : refact to new function (UpdateMsgPool)
	for _, poolMsg := range s.tempPrevoteMsgPool.Get() {
		if err := loadedState.PrevoteMsgPool.Save(&poolMsg); err != nil {
			iLogger.Debugf(nil, "[PBFT] Saving prevote msg is failed - Error: [%s]", err)
		}
	}
	s.tempPrevoteMsgPool.RemoveAllMsgs()

	loadedState = s.checkPrevote(loadedState)
	loadedState = s.checkPreCommit(loadedState)

	defer func() {
		if loadedState.StateID.ID == "" {
			returnErr = err
		} else {
			if err := s.stateRepository.Save(loadedState); err != nil {
				returnErr = err
			}
		}
	}()

	return returnErr
}

func (s *StateApi) ReceivePreCommit(msg pbft.PreCommitMsg) error {

	if err := s.tempPreCommitMsgPool.Save(&msg); err != nil {
		return err
	}

	loadedState, err := s.stateRepository.Load()
	if err != nil {
		iLogger.Debugf(nil, "[PBFT] Cannot receive PreCommit message - Error: [%s]", err)
		return nil
	}

	// todo : refact to new function (UpdateMsgPool)
	for _, poolMsg := range s.tempPreCommitMsgPool.Get() {
		if err := loadedState.SavePreCommitMsg(&poolMsg); err != nil {
			iLogger.Debugf(nil, "[PBFT] Saving precommit msg is failed - Error: [%s]", err)
		}
	}
	s.tempPreCommitMsgPool.RemoveAllMsgs()

	loadedState = s.checkPreCommit(loadedState)

	if loadedState.StateID.ID == "" {
		s.stateRepository.Remove()
		logger.Infof(nil, "[PBFT] Consensus is finished.")

		return nil
	}

	if s.stateRepository.Save(loadedState); err != nil {
		return err
	}

	return nil
}

func (s *StateApi) checkPrevote(state pbft.State) pbft.State {

	if state.CurrentStage != pbft.PREVOTE_STAGE {
		iLogger.Debugf(nil, "[PBFT] Already prevoted!")
		return state
	}

	receipients := state.GetReceipients(s.publisherID)

	if state.CheckPrevoteCondition() {
		commitMsg := pbft.NewPreCommitMsg(&state, s.publisherID)
		if err := s.propagateService.BroadcastPreCommitMsg(*commitMsg, receipients); err != nil {
			iLogger.Errorf(nil, "[PBFT] Broadcating precommit is failed - Error: [%s]", err)
			return state
		}

		state.ToPreCommitStage()
		iLogger.Infof(nil, "[PBFT] PreCommitted - Stage: [%s]", state.CurrentStage)
	}

	return state
}

func (s *StateApi) checkPreCommit(state pbft.State) pbft.State {

	if state.CurrentStage != pbft.PRECOMMIT_STAGE {
		iLogger.Debugf(nil, "[PBFT] Not Complete Prevote")
		return state
	}

	if state.CheckPreCommitCondition() {
		e := event.ConsensusFinished{
			Seal: state.Block.Seal,
			Body: state.Block.Body,
		}

		if err := s.eventService.Publish("block.confirm", e); err != nil {
			iLogger.Errorf(nil, "[PBFT] Publishing block confirm event is failed - Error: [%s]", err)
			return state
		}
		iLogger.Debug(nil, "[PBFT] Published block confirm event")

		return pbft.State{}
	}

	return state
}
