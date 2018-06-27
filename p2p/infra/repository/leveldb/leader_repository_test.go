package leveldb

import (
	"testing"

	"github.com/it-chain/it-chain-Engine/p2p"
	"github.com/magiconair/properties/assert"
)

func TestLeaderRepository_GetLeader(t *testing.T) {

	leader1 := &p2p.Leader{
		LeaderId:p2p.LeaderId{
			Id:"1",
		},
	}
	leaderRepository := NewLeaderRepository(leader1)
	leader2 := &p2p.Leader{
		LeaderId:p2p.LeaderId{
			Id:"2",
		},
	}

	leaderRepository.SetLeader(leader2)
	savedLeader := leaderRepository.GetLeader()

	// Then
	assert.Equal(t, leader2, savedLeader)
}

func TestLeaderRepository_SetLeader(t *testing.T) {

	leader1 := &p2p.Leader{
		LeaderId:p2p.LeaderId{
			Id:"1",
		},
	}
	leaderRepository := NewLeaderRepository(leader1)
	leader2 := &p2p.Leader{
		LeaderId:p2p.LeaderId{
			Id:"2",
		},
	}

	leaderRepository.SetLeader(leader2)
	savedLeader := leaderRepository.GetLeader()

	// Then
	assert.Equal(t, leader2, savedLeader)
}
