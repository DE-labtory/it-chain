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
	"sync"

	"github.com/it-chain/engine/common/event"
	"github.com/it-chain/engine/ivm"
	"github.com/it-chain/engine/ivm/api"
)

type BlockCommittedEventHandler struct {
	icodeApi api.ICodeApi
	stateApi *api.StateApi
	mutex    *sync.Mutex
}

func NewBlockCommittedEventHandler(icodeApi api.ICodeApi, stateApi *api.StateApi) *BlockCommittedEventHandler {
	return &BlockCommittedEventHandler{
		icodeApi: icodeApi,
		stateApi: stateApi,
		mutex:    &sync.Mutex{},
	}
}

func (b *BlockCommittedEventHandler) HandleBlockCommittedEventHandler(blockCommittedEvent event.BlockCommitted) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.stateApi.SetWriteSet(blockCommittedEvent.WriteSet)
}

func createRequestList(transactionList []event.Tx) []ivm.Request {

	requestList := make([]ivm.Request, 0)

	for _, transaction := range transactionList {
		requestList = append(requestList, ivm.Request{
			Function: transaction.Function,
			Args:     transaction.Args,
			ICodeID:  transaction.ICodeID,
			Type:     "invoke",
		})
	}

	return requestList
}
