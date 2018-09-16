package adapter_test

import (
	"testing"

	"github.com/it-chain/engine/api_gateway"
	"github.com/it-chain/engine/blockchain"
	"github.com/it-chain/engine/blockchain/infra/adapter"
	"github.com/it-chain/engine/blockchain/test/mock"
	"github.com/stretchr/testify/assert"
)

func TestQuerySerivce_GetRandomPeer(t *testing.T) {

	//given

	blockAdapter := mock.BlockAdapter{}

	lastBlock := blockchain.DefaultBlock{
		Seal:   []byte("seal"),
		Height: 10,
	}

	_5thBlock := blockchain.DefaultBlock{
		Seal:   []byte("seal"),
		Height: 5,
	}

	blockAdapter.GetLastBlockFunc = func(address string) (blockchain.DefaultBlock, error) {
		return lastBlock, nil
	}

	blockAdapter.GetBlockByHeightFunc = func(address string, height blockchain.BlockHeight) (blockchain.DefaultBlock, error) {
		return _5thBlock, nil
	}

	connectionQueryApi := mock.ConnectionQueryApi{}

	connectionList := []api_gateway.Connection{
		{
			ConnectionId: "connection1",
			Address:      "address1",
		},
		{
			ConnectionId: "connection2",
			Address:      "address2",
		},
		{
			ConnectionId: "connection3",
			Address:      "address3",
		},
	}

	connectionQueryApi.GetConnectionListFunc = func() ([]api_gateway.Connection, error) {
		return connectionList, nil
	}

	queryService := adapter.NewQueryService(blockAdapter, connectionQueryApi)

	//when
	randomConnection, err := queryService.GetRandomConnection()

	//then
	assert.NoError(t, err)
	assert.Contains(t, connectionList, randomConnection)

	//when
	retrieved_lastBlock, err := queryService.GetLastBlockFromConnection(randomConnection)

	//then
	assert.NoError(t, err)
	assert.Equal(t, lastBlock, retrieved_lastBlock)

	//when
	retrieved_5thBlock, err := queryService.GetBlockByHeightFromConnection(randomConnection, 5)

	//then
	assert.NoError(t, err)
	assert.Equal(t, _5thBlock, retrieved_5thBlock)

}
