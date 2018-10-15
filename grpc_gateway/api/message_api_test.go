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
	"sync"
	"testing"

	"github.com/it-chain/engine/grpc_gateway/api"
	"github.com/it-chain/engine/grpc_gateway/mock"
	"github.com/stretchr/testify/assert"
)

func TestMessageApi_DeliverMessage(t *testing.T) {

	wg := sync.WaitGroup{}
	wg.Add(1)

	mockService := &mock.GrpcService{}
	mockService.SendMessagesFunc = func(message []byte, protocol string, connIDs ...string) error {
		assert.Equal(t, message, []byte("hello"))
		assert.Equal(t, protocol, "123")
		assert.Equal(t, connIDs, []string([]string{"peer1", "peer2"}))
		wg.Done()
		return nil
	}

	mApi := api.NewMessageApi(mockService)

	err := mApi.DeliverMessage([]byte("hello"), "123", "peer1", "peer2")
	assert.NoError(t, err)

	wg.Wait()
}
