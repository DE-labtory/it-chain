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

package rpc_t_test

import (
	"testing"

	"github.com/DE-labtory/it-chain/common/rabbitmq/rpc_t"
	"github.com/stretchr/testify/assert"
)

func TestClient_Call(t *testing.T) {

	server := rpc_t.NewServer("")
	defer server.Close()

	server.Register("transaction.create", func(data []byte) ([]byte, error) {
		return data, nil
	})

	client := rpc_t.NewClient("")
	defer client.Close()

	result, err := client.Call("transaction.create", []byte("hello world"))

	if err != nil {
		panic(err)
	}

	assert.Equal(t, result.Data, []byte("hello world"))
}

func TestClient_Call_When_No_consumer(t *testing.T) {

	client := rpc_t.NewClient("")
	defer client.Close()

	_, err := client.Call("transaction.create", []byte("hello world"))

	assert.Error(t, err)
}
