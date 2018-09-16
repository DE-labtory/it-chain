package adapter

import (
	"math/rand"

	"time"

	"github.com/it-chain/engine/api_gateway"
	"github.com/it-chain/engine/blockchain"
)

type ConnectionQueryApi interface {
	GetConnectionList() ([]api_gateway.Connection, error)
	FindConnectionById(connectionId string) (api_gateway.Connection, error)
	FindConnectionByAddress(address string) (api_gateway.Connection, error)
}

type BlockAdapter interface {
	GetLastBlock(address string) (blockchain.DefaultBlock, error)
	GetBlockByHeight(address string, height blockchain.BlockHeight) (blockchain.DefaultBlock, error)
}

type QuerySerivce struct {
	blockAdapter       BlockAdapter
	connectionQueryApi ConnectionQueryApi
}

func NewQueryService(blockAdapter BlockAdapter, connectionQueryApi ConnectionQueryApi) QuerySerivce {
	return QuerySerivce{
		blockAdapter:       blockAdapter,
		connectionQueryApi: connectionQueryApi,
	}
}

func (s QuerySerivce) GetRandomConnection() (api_gateway.Connection, error) {
	connectionList, err := s.connectionQueryApi.GetConnectionList()
	if err != nil {
		return api_gateway.Connection{}, err
	}

	//ToDo: 자기 자신 빼고 아무도 없는 것이 connectionList == 0 인지, connectionList == 1 인지 확인(주석작성자: junk-sound)
	if len(connectionList) == 0 {
		return api_gateway.Connection{}, nil
	}

	randSource := rand.NewSource(time.Now().UnixNano())

	r := rand.New(randSource)

	randomIndex := r.Intn(len(connectionList))

	randomConnection := connectionList[randomIndex]

	return randomConnection, nil
}

func (s QuerySerivce) GetLastBlockFromConnection(connection api_gateway.Connection) (blockchain.DefaultBlock, error) {

	var block blockchain.DefaultBlock

	block, err := s.blockAdapter.GetLastBlock(connection.Address)

	if err != nil {
		return blockchain.DefaultBlock{}, err
	}

	return block, nil
}

func (s QuerySerivce) GetBlockByHeightFromConnection(connection api_gateway.Connection, height blockchain.BlockHeight) (blockchain.DefaultBlock, error) {

	var block blockchain.DefaultBlock

	block, err := s.blockAdapter.GetBlockByHeight(connection.Address, height)
	if err != nil {
		return blockchain.DefaultBlock{}, err
	}

	return block, nil
}
