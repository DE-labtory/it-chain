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
	"github.com/it-chain/engine/p2p/infra/mem"
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
		electionServiceOf2 := (processMap["2"].Services["ElectionService"]).(*p2p.ElectionService)

		electionServiceOf2.Election.SetState(p2p.Candidate)

		electionServiceOf1 := (processMap["1"].Services["ElectionService"]).(*p2p.ElectionService)

		t.Logf("electionService of process 1: %v", processMap["1"].Services["ElectionService"])

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

		electionServiceOf2 := (processMap["2"].Services["ElectionService"]).(*p2p.ElectionService)

		electionServiceOf1 := (processMap["1"].Services["ElectionService"]).(*p2p.ElectionService)

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
	tests := map[string]struct {
		input struct {
			processList []string
		}
	}{
		"8 node test": {
			input: struct{ processList []string }{processList: []string{"1", "2", "3", "4", "5", "6", "7", "8"}},
		},
	}

	for testName, test := range tests {
		t.Logf("running test case %s", testName)

		peerRepoList := make([]*mem.PeerRepository, 0)

		processMap := test2.SetTestEnvironment(test.input.processList)

		for _, p := range processMap {
			process := *p
			t.Logf("proces!!!!! %v", process)
			peerRepository := process.Services["PeerRepository"]
			electionService := process.Services["ElectionService"]

			for _, peerId := range test.input.processList {
				peerRepository.(*mem.PeerRepository).Save(
					p2p.Peer{
						PeerId:    struct{ Id string }{Id: peerId},
						IpAddress: peerId,
					},
				)
			}
			electionService.(*p2p.ElectionService).ElectLeaderWithRaft()
			peerRepoList = append(peerRepoList, peerRepository.(*mem.PeerRepository))
		}
		time.Sleep(5 * time.Second)

		t.Logf("peerRepo 1: %v", peerRepoList[0])
		leader1, _ := peerRepoList[0].GetLeader()
		leader2, _ := peerRepoList[1].GetLeader()
		leader3, _ := peerRepoList[2].GetLeader()
		leader4, _ := peerRepoList[3].GetLeader()
		leader5, _ := peerRepoList[4].GetLeader()
		t.Logf("leader1: %v", leader1)
		t.Logf("leader2: %v", leader2)
		t.Logf("leader3: %v", leader3)
		t.Logf("leader4: %v", leader4)
		t.Logf("leader5: %v", leader5)
		assert.Equal(t, leader2, leader1)
		assert.Equal(t, leader3, leader1)
		assert.Equal(t, leader4, leader1)
		assert.Equal(t, leader5, leader1)
	}
}

func TestNewElectionService(t *testing.T) {

}

func TestGenRandomInRange(t *testing.T) {
	v1 := p2p.GenRandomInRange(0, 10)
	v2 := p2p.GenRandomInRange(0, 10)
	v3 := p2p.GenRandomInRange(0, 10)

	t.Logf("%v", v1)
	t.Logf("%v", v2)
	t.Logf("%v", v3)
}
