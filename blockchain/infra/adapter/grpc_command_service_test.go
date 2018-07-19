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
	"testing"

	"reflect"

	"github.com/it-chain/engine/blockchain"
	"github.com/it-chain/engine/blockchain/infra/adapter"
	"github.com/stretchr/testify/assert"
)

func TestGrpcCommandService_RequestBlock(t *testing.T) {

	tests := map[string]struct {
		input struct {
			peerId blockchain.PeerId
			height uint64
		}
		err error
	}{
		"success: request block": {
			input: struct {
				peerId blockchain.PeerId
				height uint64
			}{
				peerId: blockchain.PeerId{
					Id: "1",
				},
				height: uint64(0),
			},
			err: nil,
		},
		"fail: empty node id": {
			input: struct {
				peerId blockchain.PeerId
				height uint64
			}{
				peerId: blockchain.PeerId{},
				height: uint64(0),
			},
			err: adapter.ErrEmptyNodeId,
		},
	}

	publish := func(exchange string, topic string, data interface{}) error {
		assert.Equal(t, exchange, "Command")
		assert.Equal(t, topic, "message.deliver")
		assert.Equal(t, reflect.TypeOf(data).String(), "blockchain.GrpcDeliverCommand")
		return nil
	}

	GrpcCommandService := adapter.NewGrpcCommandService(publish)

	for testName, test := range tests {
		t.Logf("running test case %s", testName)
		err := GrpcCommandService.RequestBlock(test.input.peerId, test.input.height)
		assert.Equal(t, err, test.err)
	}
}

func TestGrpcCommandService_ResponseBlock(t *testing.T) {

	tests := map[string]struct {
		input struct {
			peerId blockchain.PeerId
			block  blockchain.DefaultBlock
		}
		err error
	}{
		"success: request block": {
			input: struct {
				peerId blockchain.PeerId
				block  blockchain.DefaultBlock
			}{
				peerId: blockchain.PeerId{
					Id: "1",
				},
				block: blockchain.DefaultBlock{
					Seal: []byte("seal"),
				},
			},
			err: nil,
		},
		"fail: empty node id": {
			input: struct {
				peerId blockchain.PeerId
				block  blockchain.DefaultBlock
			}{
				peerId: blockchain.PeerId{},
				block: blockchain.DefaultBlock{
					Seal: []byte("seal"),
				},
			},
			err: adapter.ErrEmptyNodeId,
		},
		"fail: empty block seal": {
			input: struct {
				peerId blockchain.PeerId
				block  blockchain.DefaultBlock
			}{
				peerId: blockchain.PeerId{
					"1",
				},
				block: blockchain.DefaultBlock{
					Seal: nil,
				},
			},
			err: adapter.ErrEmptyBlockSeal,
		},
	}

	publish := func(exchange string, topic string, data interface{}) error {
		assert.Equal(t, exchange, "Command")
		assert.Equal(t, topic, "message.deliver")
		assert.Equal(t, reflect.TypeOf(data).String(), "blockchain.GrpcDeliverCommand")

		return nil
	}

	GrpcCommandService := adapter.NewGrpcCommandService(publish)

	for testName, test := range tests {
		t.Logf("running test case %s", testName)
		err := GrpcCommandService.ResponseBlock(test.input.peerId, &test.input.block)
		assert.Equal(t, err, test.err)
	}

}
