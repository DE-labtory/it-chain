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

package blockchain_test

import (
	"testing"

	"github.com/it-chain/engine/blockchain"
	"github.com/it-chain/engine/common/event"
	"github.com/it-chain/midgard"
	"github.com/stretchr/testify/assert"
)

func TestBlockSyncState(t *testing.T) {
	// when
	syncState := blockchain.NewBlockSyncState()

	// then
	assert.Equal(t, blockchain.BC_SYNC_STATE_AID, syncState.GetID())

	// When
	event1 := &event.SyncStart{
		EventModel: midgard.EventModel{
			ID: blockchain.BC_SYNC_STATE_AID,
		},
	}
	syncState.On(event1)

	// Then
	assert.Equal(t, blockchain.PROGRESSING, syncState.IsProgressing())

	// When
	event2 := &event.SyncDone{
		EventModel: midgard.EventModel{
			ID: blockchain.BC_SYNC_STATE_AID,
		},
	}
	syncState.On(event2)

	// Then
	assert.Equal(t, blockchain.DONE, syncState.IsProgressing())
}

func TestBlockSyncState_SetProgress(t *testing.T) {
	// when
	syncState := blockchain.NewBlockSyncState()

	// when
	syncState.SetProgress(blockchain.PROGRESSING)

	// then
	assert.Equal(t, blockchain.PROGRESSING, syncState.IsProgressing())

	// when
	syncState.SetProgress(blockchain.DONE)

	// then
	assert.Equal(t, blockchain.DONE, syncState.IsProgressing())
}
