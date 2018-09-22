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

package adapter_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/it-chain/engine/grpc_gateway"
	"github.com/it-chain/engine/grpc_gateway/infra/adapter"
	"github.com/stretchr/testify/assert"
)

func TestHttpPeerAdapter_GetAllPeerList(t *testing.T) {
	connectionList := []grpc_gateway.Connection{grpc_gateway.Connection{GrpcGatewayAddress: "1"}}

	go createMockServer("/connections", connectionList)
	time.Sleep(3 * time.Second)

	httpPeerAdapter := adapter.HttpPeerAdapter{}
	responseConnectionList, err := httpPeerAdapter.GetAllPeerList("127.0.0.1:8080")
	assert.NoError(t, err)
	assert.Equal(t, responseConnectionList, connectionList)
}

func createMockServer(url string, responseBody interface{}) {

	http.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		b, _ := json.Marshal(responseBody)
		w.Write(b)
	})

	http.ListenAndServe(":8080", nil)
}
