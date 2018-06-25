package adapter_test

import (
	"testing"
	"github.com/it-chain/it-chain-Engine/blockchain"
	"github.com/it-chain/it-chain-Engine/blockchain/infra/adapter"
	"github.com/it-chain/midgard"
	"github.com/magiconair/properties/assert"
)

type CommandHandlerMockNodeApi struct {
	UpdateNodeFunc func(node blockchain.Node) error
}
func (na CommandHandlerMockNodeApi) UpdateNode(node blockchain.Node) error {
	return na.UpdateNodeFunc(node)
}

type CommandHandlerMockBlockApi struct {}
func (ba CommandHandlerMockBlockApi) CreateGenesisBlock(genesisConfFilePath string) (blockchain.Block, error) {return nil, nil}
func (ba CommandHandlerMockBlockApi) CreateBlock(txList []blockchain.Transaction) (blockchain.Block, error) {return nil, nil}


func TestBlockchainCommandHandler_HandleUpdateNodeCommand(t *testing.T) {
	tests := map[string] struct {
		input struct {
			ID string
			nodeId string
			address string
		}
		err error
	}{
		"success": {
			input: struct {
				ID string
				nodeId string
				address string
			}{ID: string("zf"), nodeId: string("zf2"), address: string("11.22.33.44")},
			err: nil,
		},
		"empty eventId test": {
			input: struct {
				ID string
				nodeId string
				address string
			}{ID: string(""), nodeId: string("zf2"), address: string("11.22.33.44")},
			err: nil,
		},
		"empty nodeId test": {
			input: struct {
				ID string
				nodeId string
				address string
			}{ID: string("zf"), nodeId: string(""), address: string("11.22.33.44")},
			err: nil,
		},
		"empty ip address test": {
			input: struct {
				ID string
				nodeId string
				address string
			}{ID: string("zf"), nodeId: string("zf2"), address: string("")},
			err: nil,
		},
	}

	mockNodeApi := CommandHandlerMockNodeApi{}
	mockNodeApi.UpdateNodeFunc = func(node blockchain.Node) error {
		assert.Equal(t, node.NodeId.Id, string("zf2"))
		assert.Equal(t, node.IpAddress, string("11.22.33.44"))
		return nil
	}

	commandHandler := adapter.NewCommandHandler(CommandHandlerMockBlockApi{}, mockNodeApi)

	for testName, test := range tests {
		t.Logf("running test case %s", testName)

		command := blockchain.NodeUpdateCommand{
			EventModel: midgard.EventModel{
				ID: test.input.ID,
			},
			Node: blockchain.Node{
				NodeId: blockchain.NodeId{
					test.input.nodeId,
				},
				IpAddress: test.input.address,
			},
		}
		commandHandler.HandleUpdateNodeCommand(command)
	}
}


