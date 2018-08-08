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

package blockchain

import (
	"errors"
	"fmt"

	"github.com/it-chain/engine/common/event"
	"github.com/it-chain/midgard"
)

// BlockSyncState Aggregate ID
var BC_SYNC_STATE_AID = "BC_SYNC_STATE_AID"

type ProgressState bool

const (
	PROGRESSING ProgressState = true
	DONE        ProgressState = false
)

type SyncState interface {
	SetProgress(state ProgressState)
}

// 현재 블록 동기화가 진행 중인지 정보를 가진다.
type BlockSyncState struct {
	midgard.AggregateModel
	isProgress ProgressState
}

func NewBlockSyncState() *BlockSyncState {
	return &BlockSyncState{
		AggregateModel: midgard.AggregateModel{
			ID: BC_SYNC_STATE_AID,
		},
		isProgress: DONE,
	}
}

func (bss *BlockSyncState) GetID() string {
	return BC_SYNC_STATE_AID
}

func (bss *BlockSyncState) SetProgress(state ProgressState) {
	var event midgard.Event
	if state == PROGRESSING {
		event = createSyncStartEvent()
	} else { // state == DONE
		event = createSyncDoneEvent()
	}
	bss.On(event)
}

func createSyncStartEvent() *event.SyncStart {
	return &event.SyncStart{
		EventModel: midgard.EventModel{
			ID: BC_SYNC_STATE_AID,
		},
	}
}

func createSyncDoneEvent() *event.SyncDone {
	return &event.SyncDone{
		EventModel: midgard.EventModel{
			ID: BC_SYNC_STATE_AID,
		},
	}
}

func (bss *BlockSyncState) IsProgressing() ProgressState {
	return bss.isProgress
}

func (bss *BlockSyncState) On(e midgard.Event) error {
	switch v := e.(type) {

	case *(event.SyncStart):
		bss.isProgress = PROGRESSING

	case *(event.SyncDone):
		bss.isProgress = DONE

	default:
		return errors.New(fmt.Sprintf("unhandled event [%s]", v))
	}

	return nil
}
