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

package adapter

import (
	"github.com/it-chain/engine/blockchain"
)

type EventHandler struct {
	blockApi BlockApi
}

func NewEventHandler(api BlockApi) *EventHandler {
	return &EventHandler{
		blockApi: api,
	}
}

// TODO: write test case
func (eh *EventHandler) HandleBlockAddToPoolEvent(event blockchain.BlockAddToPoolEvent) error {
	if err := isBlockHasMissingProperty(event); err != nil {
		return err
	}
	height := event.Height
	err := eh.blockApi.CheckAndSaveBlockFromPool(height)

	if err != nil {
		return err
	}

	return nil
}

func isBlockHasMissingProperty(event blockchain.BlockAddToPoolEvent) error {
	if event.Seal == nil || event.PrevSeal == nil || event.Height == 0 ||
		event.TxList == nil || event.TxSeal == nil || event.Timestamp.IsZero() || event.Creator == nil {
		return ErrBlockMissingProperties
	}
	return nil
}
