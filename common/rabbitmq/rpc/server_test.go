/*
 * Copyright 2018 DE-labtory
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

package rpc_test

import (
	"testing"

	"github.com/DE-labtory/it-chain/common/rabbitmq/rpc"
	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {

	server := rpc.NewServer("")

	err := server.Register("transaction.create", func(sample Sample) Sample {

		assert.Equal(t, sample.ID, "123")
		assert.Equal(t, sample.TxData.ID, "1234")

		return Sample{ID: "1234"}
	})

	assert.NoError(t, err)
}
