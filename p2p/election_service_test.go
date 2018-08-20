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

package p2p_test

import (
	"testing"

	"time"

	"github.com/it-chain/engine/p2p"
	test2 "github.com/it-chain/engine/p2p/test"
	"github.com/magiconair/properties/assert"
)

func TestElectionService_Vote(t *testing.T) {
	tests := map[string]struct {
		input struct {
			processList []string
		}
	}{
		"4 node test": {
			input: struct{ processList []string }{processList: []string{"1", "2", "3", "4"}},
		},
	}

	for testName, test := range tests {
		t.Logf("running test case %s", testName)

		processMap := test2.SetTestEnvironment(test.input.processList)
		electionServiceOf2 := (*processMap["2"].Services["ElectionService"]).(p2p.ElectionService)

		electionServiceOf2.Election.SetState("candidate")

		electionServiceOf1 := (*processMap["1"].Services["ElectionService"]).(p2p.ElectionService)

		t.Logf("electionService: %v", processMap["1"].Services["ElectionService"])

		t.Logf("before vote check state: %v", electionServiceOf2.Election.GetState())
		t.Logf("election of 2: %v", electionServiceOf2.Election)
		electionServiceOf1.Vote("2")
		t.Logf("after vote check state: %v", electionServiceOf2.Election.GetState())

		time.Sleep(5 * time.Second)

		assert.Equal(t, electionServiceOf2.Election.GetVoteCount(), 1)
	}

}

func TestElectionService_RequestVote(t *testing.T) {
	tests := map[string]struct {
		input struct {
			processList []string
		}
	}{
		"4 node test": {
			input: struct{ processList []string }{processList: []string{"1", "2", "3", "4"}},
		},
	}

	for testName, test := range tests {
		t.Logf("running test case %s", testName)

		processMap := test2.SetTestEnvironment(test.input.processList)
		electionServiceOf2 := (*processMap["2"].Services["ElectionService"]).(p2p.ElectionService)

		electionServiceOf1 := (*processMap["1"].Services["ElectionService"]).(p2p.ElectionService)

		t.Logf("electionService: %v", processMap["1"].Services["ElectionService"])

		t.Logf("before vote check state: %v", electionServiceOf2.Election.GetState())
		t.Logf("election of 2: %v", electionServiceOf2.Election)
		electionServiceOf1.RequestVote([]string{"2"})
		t.Logf("after vote check state: %v", electionServiceOf2.Election.GetState())

		time.Sleep(5 * time.Second)

		assert.Equal(t, electionServiceOf1.Election.GetVoteCount(), 1)
	}
}

func TestElectionService_DecideToBeLeader(t *testing.T) {

}

func TestElectionService_BroadcastLeader(t *testing.T) {

}

func TestElectionService_ElectLeaderWithRaft(t *testing.T) {

}

func TestNewElectionService(t *testing.T) {

}
