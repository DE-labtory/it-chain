package adapter_test

import (
	"testing"
	"github.com/it-chain/it-chain-Engine/blockchain"
	"fmt"
)

type CommandHandlerMockNodeApi struct {}
func (na *CommandHandlerMockNodeApi) UpdateNode(node blockchain.Node) error {return nil}

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
			}{ID: string("zf"), nodeId: string("1"), address: string("11.22.33.44")},
			err: nil,
		},
		"empty eventId test": {
			input: struct {
				ID string
				nodeId string
				address string
			}{ID: string("zf"), nodeId: string(""), address: string("11.22.33.44")}
		},
	}
	fmt.Println(tests)
}


