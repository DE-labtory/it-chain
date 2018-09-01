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
	"github.com/it-chain/engine/common/event"
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
	Id         string
	IsProgress ProgressState
}

func NewBlockSyncState() *BlockSyncState {
	return &BlockSyncState{
		Id:         BC_SYNC_STATE_AID,
		IsProgress: DONE,
	}
}

func (bss *BlockSyncState) GetID() string {
	return BC_SYNC_STATE_AID
}

func (bss *BlockSyncState) SetProgress(state ProgressState) {
	if state == PROGRESSING {
		bss.IsProgress = PROGRESSING
	} else { // state == DONE
		bss.IsProgress = DONE
	}
}

func createSyncStartEvent() *event.SyncStart {
	return &event.SyncStart{
		EventId: BC_SYNC_STATE_AID,
	}
}

func createSyncDoneEvent() *event.SyncDone {
	return &event.SyncDone{
		EventId: BC_SYNC_STATE_AID,
	}
}

func (bss *BlockSyncState) IsProgressing() ProgressState {
	return bss.IsProgress
}
