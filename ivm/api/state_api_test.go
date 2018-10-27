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

package api_test

import (
	"testing"

	"github.com/it-chain/engine/ivm"
	"github.com/it-chain/engine/ivm/api"
	"github.com/it-chain/engine/ivm/infra/kvstore"
	"github.com/stretchr/testify/assert"
)

func TestStateApi_SetWriteSet(t *testing.T) {
	stateRepository := kvstore.NewLevelDBStateRepository()
	stateService := ivm.NewStateService(stateRepository)

	stateApi := api.NewStateApi(stateRepository, stateService)

	txList := make([]ivm.TransactionWriteList, 0)
	txList = append(txList, ivm.TransactionWriteList{
		Id: "1",
		WriteList: []ivm.Write{{
			Key:   []byte("1"),
			Value: []byte("1"),
		}},
	})
	stateApi.SetWriteSet(txList)

	value, _ := stateRepository.Get([]byte("1"))

	assert.Equal(t, value, []byte("1"))
	stateRepository.Close()
}
