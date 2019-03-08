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
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/DE-labtory/engine/common/rabbitmq/rpc"
	"github.com/DE-labtory/midgard"
	"github.com/stretchr/testify/assert"
)

type Sample struct {
	midgard.EventModel
	ID     string
	TxData TxData
}

type TxData struct {
	ID string
}

func TestClient_Call(t *testing.T) {

	wg := sync.WaitGroup{}
	wg.Add(1)

	server := rpc.NewServer("")
	defer server.Close()

	server.Register("transaction.create", func(sample Sample) (Sample, rpc.Error) {
		assert.Equal(t, sample.ID, "123")
		assert.Equal(t, sample.TxData.ID, "1234")

		return Sample{ID: "1234"}, rpc.Error{}
	})

	client := rpc.NewClient("")
	defer client.Close()

	err := client.Call("transaction.create", Sample{ID: "123", TxData: TxData{ID: "1234"}}, func(sample Sample, err rpc.Error) {

		assert.True(t, err.IsNil())

		assert.Equal(t, sample.ID, "1234")
		wg.Done()
	})

	if err != nil {
		panic(err)
	}

	wg.Wait()
}

func TestClient_InfiniteCall(t *testing.T) {

	wg := sync.WaitGroup{}
	wg.Add(1)

	server := rpc.NewServer("")
	defer server.Close()

	server.Register("transaction.create", func(sample Sample) (Sample, rpc.Error) {
		assert.Equal(t, sample.ID, "123")
		assert.Equal(t, sample.TxData.ID, "1234")

		time.Sleep(time.Second * 190)

		return Sample{ID: "1234"}, rpc.Error{}
	})

	client := rpc.NewClient("")
	defer client.Close()

	err := client.Call("transaction.create", Sample{ID: "123", TxData: TxData{ID: "1234"}}, func(sample Sample, err rpc.Error) {

		assert.True(t, err.IsNil())
		assert.Equal(t, sample.ID, "1234")
		wg.Done()
	})

	assert.Equal(t, err, rpc.ErrTimeout, "Timeout queue")
	wg.Done()

	wg.Wait()
}

func TestClient_Call_When_No_consumer(t *testing.T) {

	client := rpc.NewClient("")
	defer client.Close()
	err := client.Call("transaction.create", Sample{}, func(sample Sample) {
		fmt.Println("callbacked!")
	})

	assert.Error(t, err)
}
