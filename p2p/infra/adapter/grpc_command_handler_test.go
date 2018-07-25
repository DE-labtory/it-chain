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
	"testing"

	"github.com/it-chain/engine/common/command"
	"github.com/it-chain/engine/p2p"
	"github.com/it-chain/engine/p2p/infra/adapter"
	"github.com/it-chain/engine/p2p/test/mock"
	"github.com/it-chain/midgard"
	"github.com/magiconair/properties/assert"
)

func TestGrpcCommandHandler_HandleMessageReceive(t *testing.T) {

	leader := p2p.Leader{}
	leaderByte, _ := json.Marshal(leader)

	//todo error case write!
	tests := map[string]struct {
		input struct {
			id       string
			protocol string
			body     []byte
		}
		err error
	}{
		"leader info deliver test success": {
			input: struct {
				id       string
				protocol string
				body     []byte
			}{
				id:       "1",
				protocol: string("LeaderInfoDeliverProtocol"),
				body:     leaderByte,
			},
			err: nil,
		},
		"leader info deliver test empty leader id":       {},
		"leader table deliver test success":              {},
		"peer leader table deliver test empty peer list": {},
		"peer leader table deliver test empty leader id": {},
	}

	leaderApi := &mock.MockLeaderApi{}

	electionService := p2p.ElectionService{}

	communicationApi := &mock.MockCommunicationApi{}

	pLTableService := &mock.MockPLTableService{}

	messageHandler := adapter.NewGrpcCommandHandler(leaderApi, electionService, communicationApi, pLTableService)

	for testName, test := range tests {
		grpcReceiveCommand := command.ReceiveGrpc{
			CommandModel: midgard.CommandModel{
				ID: test.input.id,
			},
			Body:     test.input.body,
			Protocol: test.input.protocol,
		}
		t.Logf("running test case %s", testName)
		err := messageHandler.HandleMessageReceive(grpcReceiveCommand)
		assert.Equal(t, err, test.err)
	}

}
