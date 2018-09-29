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

	"time"

	"github.com/it-chain/engine/consensus/pbft"
	"github.com/it-chain/engine/consensus/pbft/api"
	"github.com/it-chain/engine/consensus/pbft/infra/mem"
	test2 "github.com/it-chain/engine/consensus/pbft/test"
	"github.com/it-chain/engine/p2p/test/mock"
	"github.com/stretchr/testify/assert"
)

func TestElectionApi_Vote(t *testing.T) {
	tests := map[string]struct {
		input struct {
			processList []string
		}
	}{
		"4 node test": {
			input: struct{ processList []string }{
				processList: []string{"1", "2", "3", "4"},
			},
		},
	}

	for testName, test := range tests {
		t.Logf("running test case %s", testName)
		env := test2.SetTestEnvironment(test.input.processList)

		electionServiceOf2 := (env.ProcessMap["2"].Services["ElectionService"]).(*pbft.ElectionService)

		electionServiceOf2.SetState(pbft.CANDIDATE)

		electionApi1 := (env.ProcessMap["1"].Services["ElectionApi"]).(*api.ElectionApi)
		t.Logf("electionService of process 1: %v", env.ProcessMap["1"].Services["ElectionService"])

		t.Logf("before vote check state: %v", electionServiceOf2.GetState())
		t.Logf("election of 2: %v", electionServiceOf2)
		electionApi1.Vote("2")
		t.Logf("after vote check state: %v", electionServiceOf2.GetState())

		time.Sleep(2 * time.Second)

		assert.Equal(t, 1, electionServiceOf2.GetVoteCount())
	}
}

func TestElectionApi_DecideToBeLeader(t *testing.T) {
	tests := map[string]struct {
		input struct {
			election *pbft.ElectionService
		}
		output struct {
			voteCount int
		}
	}{
		"when election is ticking state, vote count not reached majority": {
			input: struct{ election *pbft.ElectionService }{
				election: pbft.NewElectionService("this.should.not.broadcast", 30, pbft.TICKING, 0),
			},
			output: struct{ voteCount int }{
				voteCount: 0,
			},
		},
		"when election is candidate state, vote count not reached majority": {
			input: struct{ election *pbft.ElectionService }{
				election: pbft.NewElectionService("this.should.not.broadcast", 30, pbft.CANDIDATE, 0),
			},
			output: struct{ voteCount int }{
				voteCount: 1,
			},
		},
		"when election is ticking state, vote count reached majority": {
			input: struct{ election *pbft.ElectionService }{
				election: pbft.NewElectionService("this.is.input.address", 30, pbft.TICKING, 2),
			},
			output: struct{ voteCount int }{
				voteCount: 2,
			},
		},
		"when election is candidate state, vote count reached majority": {
			input: struct{ election *pbft.ElectionService }{
				election: pbft.NewElectionService("this.is.input.address", 30, pbft.CANDIDATE, 1),
			},
			output: struct{ voteCount int }{
				voteCount: 2,
			},
		},
	}

	for testName, test := range tests {
		t.Logf("running test case %s", testName)

		networkManager := test2.SetTestEnvironment([]string{"1", "2", "3"})
		networkManager.ProcessMap["1"].Services["ElectionApi"].(*api.ElectionApi).ElectionService.SetLeftTime(test.input.election.GetLeftTime())
		networkManager.ProcessMap["1"].Services["ElectionApi"].(*api.ElectionApi).ElectionService.SetState(test.input.election.GetState())
		networkManager.ProcessMap["1"].Services["ElectionApi"].(*api.ElectionApi).ElectionService.SetVoteCount(test.input.election.GetVoteCount())

		electionApi1 := networkManager.ProcessMap["1"].Services["ElectionApi"].(*api.ElectionApi)

		// when, then
		err := electionApi1.DecideToBeLeader()
		assert.NoError(t, err)
		// when, then
		count := electionApi1.GetVoteCount()
		assert.Equal(t, test.output.voteCount, count)
	}

}

func TestElectionApi_RequestVote(t *testing.T) {
	tests := map[string]struct {
		input struct {
			processList []string
		}
	}{
		"4 node test": {
			input: struct{ processList []string }{
				processList: []string{"1", "2", "3", "4"},
			},
		},
	}

	for testName, test := range tests {
		t.Logf("running test case %s", testName)

		env := test2.SetTestEnvironment(test.input.processList)

		electionApi := (env.ProcessMap["1"].Services["ElectionApi"]).(*api.ElectionApi)

		t.Logf("electionService of 1: %v", env.ProcessMap["1"].Services["ElectionApi"])
		t.Logf("electionService of 2: %v", env.ProcessMap["2"].Services["ElectionApi"])

		t.Logf("before vote check state: %v", electionApi.GetState())

		electionApi.SetState(pbft.CANDIDATE)
		electionApi.RequestVote([]string{"2"})
		t.Logf("after vote check state: %v", electionApi.GetState())

		time.Sleep(4 * time.Second)

		assert.Equal(t, electionApi.GetVoteCount(), 1)
	}
}

func TestElectionApi_ElectLeaderWithRaft(t *testing.T) {
	tests := map[string]struct {
		input struct {
			processList []string
		}
	}{
		"8 node test": {
			input: struct{ processList []string }{
				processList: []string{"1", "2", "3", "4", "5", "6", "7", "8"},
			},
		},
	}

	for testName, test := range tests {
		t.Logf("running test case %s", testName)

		env := test2.SetTestEnvironment(test.input.processList)

		for _, p := range env.ProcessMap {
			process := *p
			electionApi := process.Services["ElectionApi"].(*api.ElectionApi)
			go electionApi.ElectLeaderWithRaft()
		}

		time.Sleep(5 * time.Second)

		leader1 := env.ProcessMap["1"].Services["ParliamentApi"].(*api.ParliamentApi).GetLeader()
		leader2 := env.ProcessMap["2"].Services["ParliamentApi"].(*api.ParliamentApi).GetLeader()
		leader3 := env.ProcessMap["3"].Services["ParliamentApi"].(*api.ParliamentApi).GetLeader()
		leader4 := env.ProcessMap["4"].Services["ParliamentApi"].(*api.ParliamentApi).GetLeader()

		t.Logf("leader1: %v", leader1)
		t.Logf("leader2: %v", leader2)
		t.Logf("leader3: %v", leader3)
		t.Logf("leader4: %v", leader4)

		assert.Equal(t, leader2, leader1)
		assert.Equal(t, leader3, leader1)
		assert.Equal(t, leader4, leader1)
	}
}

func TestGenRandomInRange(t *testing.T) {
	v1 := pbft.GenRandomInRange(0, 10)
	v2 := pbft.GenRandomInRange(0, 10)
	v3 := pbft.GenRandomInRange(0, 10)

	t.Logf("%v", v1)
	t.Logf("%v", v2)
	t.Logf("%v", v3)
}

func TestElectionApi_GetCandidate(t *testing.T) {
	api := setElectionApi()
	api.ElectionService.SetCandidate(pbft.Representative{
		ID: "1",
	})

	assert.Equal(t, api.GetCandidate().ID, "1")
}

func TestElectionApi_GetState(t *testing.T) {
	api := setElectionApi()

	assert.Equal(t, api.GetState(), pbft.CANDIDATE)
}

func TestElectionApi_GetVoteCount(t *testing.T) {
	api := setElectionApi()

	assert.Equal(t, api.GetVoteCount(), 0)
}

func setElectionApi() *api.ElectionApi {

	electionService := pbft.NewElectionService("1", 30, pbft.CANDIDATE, 0)
	parliament := pbft.NewParliament()
	parliamentRepository := mem.NewParliamentRepository()

	parliament.AddRepresentative(pbft.NewRepresentative("1"))
	parliament.AddRepresentative(pbft.NewRepresentative("2"))
	parliament.AddRepresentative(pbft.NewRepresentative("3"))
	parliamentRepository.Save(parliament)

	eventService := &mock.MockEventService{}
	api := api.NewElectionApi(electionService, parliamentRepository, eventService)
	return api
}
