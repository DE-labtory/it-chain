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

	"github.com/it-chain/engine/common"
	"github.com/it-chain/engine/common/command"

	"time"

	"github.com/it-chain/engine/p2p"
	"github.com/it-chain/engine/p2p/infra/mem"
	test2 "github.com/it-chain/engine/p2p/test"
	"github.com/it-chain/engine/p2p/test/mock"
	"github.com/stretchr/testify/assert"
)

func TestElectionService_Vote(t *testing.T) {
	tests := map[string]struct {
		input struct {
			processList []test2.ProcessIdentity
		}
	}{
		"4 node test": {
			input: struct{ processList []test2.ProcessIdentity }{
				processList: []test2.ProcessIdentity{
					{"1", ""},
					{"2", ""},
					{"3", ""},
					{"4", ""},
				},
			},
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
			processList []test2.ProcessIdentity
		}
	}{
		"4 node test": {
			input: struct{ processList []test2.ProcessIdentity }{
				processList: []test2.ProcessIdentity{
					{"1", ""},
					{"2", ""},
					{"3", ""},
					{"4", ""},
				},
			},
		},
	}

	for testName, test := range tests {
		t.Logf("running test case %s", testName)

		processMap := test2.SetTestEnvironment(test.input.processList)

		electionServiceOf1 := (processMap["1"].Services["ElectionService"]).(*p2p.ElectionService)

		t.Logf("electionService of 1: %v", processMap["1"].Services["ElectionService"])
		t.Logf("electionService of 2: %v", processMap["2"].Services["ElectionService"])

		t.Logf("before vote check state: %v", electionServiceOf1.Election.GetState())
		t.Logf("election of 1: %v", electionServiceOf1.Election)

		electionServiceOf1.RequestVote([]string{"2"})
		t.Logf("after vote check state: %v", electionServiceOf1.Election.GetState())

		time.Sleep(4 * time.Second)

		assert.Equal(t, electionServiceOf1.Election.GetVoteCount(), 1)
	}
}

func TestElectionService_DecideToBeLeader(t *testing.T) {
	tests := map[string]struct {
		input struct {
			election p2p.Election
		}
		output struct {
			voteCount int
		}
	}{
		"when election is ticking state, vote count not reached majority": {
			input: struct{ election p2p.Election }{
				election: p2p.NewElection("this.should.not.broadcast", 30, p2p.Ticking, 0),
			},
			output: struct{ voteCount int }{
				voteCount: 0,
			},
		},
		"when election is candidate state, vote count not reached majority": {
			input: struct{ election p2p.Election }{
				election: p2p.NewElection("this.should.not.broadcast", 30, p2p.Candidate, 0),
			},
			output: struct{ voteCount int }{
				voteCount: 1,
			},
		},
		"when election is ticking state, vote count reached majority": {
			input: struct{ election p2p.Election }{
				election: p2p.NewElection("this.is.input.address", 30, p2p.Ticking, 2),
			},
			output: struct{ voteCount int }{
				voteCount: 2,
			},
		},
		"when election is candidate state, vote count reached majority": {
			input: struct{ election p2p.Election }{
				election: p2p.NewElection("this.is.input.address", 30, p2p.Candidate, 1),
			},
			output: struct{ voteCount int }{
				voteCount: 2,
			},
		},
	}

	client := mock.MockClient{}
	client.CallFunc = func(queue string, params interface{}, callback interface{}) error {
		message := p2p.UpdateLeaderMessage{}
		common.Deserialize(params.(command.DeliverGrpc).Body, &message)

		assert.Equal(t, "this.is.input.address", message.Peer.IpAddress)

		assert.Equal(t, "message.deliver", queue)
		return nil
	}
	queryService := mock.MockPeerQueryService{}
	queryService.GetPLTableFunc = func() (p2p.PLTable, error) {
		return p2p.PLTable{
			Leader: p2p.Leader{LeaderId: p2p.LeaderId{Id: "FollowMe"}},
			PeerTable: map[string]p2p.Peer{
				"1": p2p.Peer{IpAddress: "1.ipAddr", PeerId: p2p.PeerId{Id: "1"}},
				"2": p2p.Peer{IpAddress: "2.ipAddr", PeerId: p2p.PeerId{Id: "2"}},
				"3": p2p.Peer{IpAddress: "3.ipAddr", PeerId: p2p.PeerId{Id: "3"}},
			},
		}, nil
	}

	for testName, test := range tests {
		t.Logf("running test case %s", testName)

		electionService := p2p.NewElectionService(&test.input.election, queryService, client)

		// when, then
		err := electionService.DecideToBeLeader()
		assert.NoError(t, err)
		// when, then
		count := electionService.Election.GetVoteCount()
		assert.Equal(t, test.output.voteCount, count)
	}

}

func TestElectionService_DecideToBeLeader_WhenElectionCandiateState(t *testing.T) {
	election := p2p.NewElection("ip.address", 30, p2p.Candidate, 0)

	queryService := mock.MockPeerQueryService{}
	queryService.GetPLTableFunc = func() (p2p.PLTable, error) {
		return p2p.PLTable{
			Leader: p2p.Leader{LeaderId: p2p.LeaderId{Id: "FollowMe"}},
			PeerTable: map[string]p2p.Peer{
				"1": p2p.Peer{IpAddress: "1.ipAddr", PeerId: p2p.PeerId{Id: "1"}},
				"2": p2p.Peer{IpAddress: "2.ipAddr", PeerId: p2p.PeerId{Id: "2"}},
				"3": p2p.Peer{IpAddress: "3.ipAddr", PeerId: p2p.PeerId{Id: "3"}},
			},
		}, nil
	}

	client := mock.MockClient{}

	electionService := p2p.NewElectionService(&election, queryService, client)

	// when, then
	err := electionService.DecideToBeLeader()
	assert.NoError(t, err)
	// when, then
	count := electionService.Election.GetVoteCount()
	assert.Equal(t, 1, count)
}

func TestElectionService_BroadcastLeader(t *testing.T) {
	tests := map[string]struct {
		input struct {
			processList []test2.ProcessIdentity
		}
	}{
		"4 node test": {
			input: struct{ processList []test2.ProcessIdentity }{
				processList: []test2.ProcessIdentity{
					{"1", "1.ipAddress"},
					{"2", "2.ipAddress"},
					{"3", "3.ipAddress"},
					{"4", "4.ipAddress"},
				},
			},
		},
	}

	for testName, test := range tests {
		t.Logf("running test case %s", testName)

		processMap := test2.SetTestEnvironment(test.input.processList)

		newLeader := p2p.Peer{IpAddress: "1.ipAddress", PeerId: p2p.PeerId{Id: "1"}}

		electionServiceOf2 := (processMap["2"].Services["ElectionService"]).(*p2p.ElectionService)
		electionServiceOf2.Election.SetCandidate(&newLeader)

		electionServiceOf1 := (processMap["1"].Services["ElectionService"]).(*p2p.ElectionService)
		electionServiceOf1.Election.SetState(p2p.Candidate)
		electionServiceOf1.BroadcastLeader(newLeader)

		time.Sleep(4 * time.Second)

		peerRepositoryOf2 := (processMap["2"].Services["PeerRepository"]).(*mem.PeerRepository)
		broadcastedLeaderOf2, _ := peerRepositoryOf2.GetLeader()

		t.Logf("broadcasted leader %v", broadcastedLeaderOf2)
		assert.Equal(t, "1", broadcastedLeaderOf2.GetID())
	}

}

func TestElectionService_ElectLeaderWithRaft(t *testing.T) {
	tests := map[string]struct {
		input struct {
			processList []test2.ProcessIdentity
		}
	}{
		"8 node test": {
			input: struct{ processList []test2.ProcessIdentity }{
				processList: []test2.ProcessIdentity{
					{"1", ""},
					{"2", ""},
					{"3", ""},
					{"4", ""},
					{"5", ""},
					{"6", ""},
					{"7", ""},
					{"8", ""},
				},
			},
		},
	}

	for testName, test := range tests {
		t.Logf("running test case %s", testName)

		peerRepoList := make([]*mem.PeerRepository, 0)

		processMap := test2.SetTestEnvironment(test.input.processList)

		for _, p := range processMap {
			process := *p
			peerRepository := process.Services["PeerRepository"]
			electionService := process.Services["ElectionService"]

			for _, peer := range test.input.processList {
				peerRepository.(*mem.PeerRepository).Save(
					p2p.Peer{
						PeerId:    struct{ Id string }{Id: peer.Id},
						IpAddress: peer.Id,
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
	election := p2p.NewElection("this.is.ip.addres", 30, p2p.Ticking, 123)
	queryService := mock.MockPeerQueryService{}
	client := mock.MockClient{}

	electionService := p2p.NewElectionService(&election, queryService, client)

	assert.Equal(t, 123, electionService.Election.GetVoteCount())
	assert.Equal(t, 30, electionService.Election.GetLeftTime())
	assert.Equal(t, p2p.Ticking, electionService.Election.GetState())
}

func TestGenRandomInRange(t *testing.T) {
	v1 := p2p.GenRandomInRange(0, 10)
	v2 := p2p.GenRandomInRange(0, 10)
	v3 := p2p.GenRandomInRange(0, 10)

	t.Logf("%v", v1)
	t.Logf("%v", v2)
	t.Logf("%v", v3)
}
