package leveldb

import (
	"testing"
	"os"
	"github.com/magiconair/properties/assert"
	"github.com/it-chain/it-chain-Engine/p2p"
)

func TestLeaderRepository_GetLeader(t *testing.T) {
	// Given
	dbPath := "./.test"

	leaderRepository := NewLeaderRepository(dbPath)

	leader := p2p.Leader{
		LeaderId:p2p.LeaderId{
			"777",
		},
	}

	defer func() {
		leaderRepository.leveldb.Close()
		os.RemoveAll(dbPath)
	}()

	leaderRepository.SetLeader(leader)
	leader2 := leaderRepository.GetLeader()

	// Then
	assert.Equal(t, leader, leader2)
}

func TestLeaderRepository_SetLeader(t *testing.T) {
	// Given
	dbPath := "./.test"

	leaderRepository := NewLeaderRepository(dbPath)

	leader := p2p.Leader{
		LeaderId:p2p.LeaderId{
			"777",
		},
	}

	defer func() {
		leaderRepository.leveldb.Close()
		os.RemoveAll(dbPath)
	}()

	// When
	leaderRepository.SetLeader(leader)
	leader2 := *leaderRepository.GetLeader()

	// Then
	assert.Equal(t, leader, leader2)
}