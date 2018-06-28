package memory_test

import (
	"testing"

	"github.com/it-chain/it-chain-Engine/p2p"
	"github.com/magiconair/properties/assert"
	"github.com/it-chain/it-chain-Engine/p2p/infra/repository/memory"
)

func TestLeaderRepository_GetLeader(t *testing.T) {
	tests := map[string]struct{
		input p2p.Leader
		outPut p2p.Leader
		err error
	}{
		"success":{
			input:p2p.Leader{
				LeaderId:p2p.LeaderId{
					Id:"2",
				},
			},
			outPut: p2p.Leader{
				LeaderId:p2p.LeaderId{
					Id:"2",
				},
			},
			err:nil,
		},
	}

	firstLeader := p2p.Leader{
		LeaderId:p2p.LeaderId{
			Id:"1",
		},
	}

	leaderRepository := memory.NewLeaderRepository(firstLeader)

	for testName, test := range tests{
		t.Logf("running test case %s", testName)
		leaderRepository.SetLeader(test.input)
		gotLeader := leaderRepository.GetLeader()
		assert.Equal(t, gotLeader, test.outPut)
	}
}

func TestLeaderRepository_SetLeader(t *testing.T) {
	tests := map[string]struct{
		input p2p.Leader
		outPut p2p.Leader
		err error
	}{
		"success":{
			input:p2p.Leader{
				LeaderId:p2p.LeaderId{
					Id:"2",
				},
			},
			outPut: p2p.Leader{
				LeaderId:p2p.LeaderId{
					Id:"2",
				},
			},
			err:nil,
		},
	}

	firstLeader := p2p.Leader{
		LeaderId:p2p.LeaderId{
			Id:"1",
		},
	}

	leaderRepository := memory.NewLeaderRepository(firstLeader)

	for testName, test := range tests{
		t.Logf("running test case %s", testName)
		leaderRepository.SetLeader(test.input)
		gotLeader := leaderRepository.GetLeader()
		assert.Equal(t, gotLeader, test.outPut)
	}

}
