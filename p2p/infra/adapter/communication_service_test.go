package adapter_test

import (
	"reflect"
	"testing"

	"github.com/it-chain/engine/p2p"
	"github.com/it-chain/engine/p2p/infra/adapter"
	"github.com/magiconair/properties/assert"
)

func TestGrpcCommandService_DeliverPLTable(t *testing.T) {
	tests := map[string]struct {
		input struct {
			connectionId string
			pLTable      p2p.PLTable
		}
		err error
	}{
		"empty peer list test": {
			input: struct {
				connectionId string
				pLTable      p2p.PLTable
			}{
				connectionId: "1",
				pLTable:      p2p.PLTable{},
			},
			err: p2p.ErrEmptyPeerTable,
		},
		"empty connection id test": {
			input: struct {
				connectionId string
				pLTable      p2p.PLTable
			}{
				connectionId: "",
				pLTable:      p2p.PLTable{},
			},
			err: adapter.ErrEmptyConnectionId,
		},
		"success": {
			input: struct {
				connectionId string
				pLTable      p2p.PLTable
			}{
				connectionId: "1",
				pLTable: p2p.PLTable{
					Leader: p2p.Leader{},
					PeerTable: map[string]p2p.Peer{
						"123": {
							PeerId: p2p.PeerId{
								Id: "123",
							},
							IpAddress: "123",
						},
					},
				},
			},
			err: nil,
		},
	}
	publish := func(exchange string, topic string, data interface{}) error {
		{
			assert.Equal(t, exchange, "Command")
			assert.Equal(t, topic, "message.deliver")
			assert.Equal(t, reflect.TypeOf(data).String(), "p2p.GrpcDeliverCommand")

			return nil
		}
	}

	communicationService := adapter.NewCommunicationService(publish)

	for testName, test := range tests {
		t.Logf("running test case %s", testName)
		err := communicationService.DeliverPLTable(test.input.connectionId, test.input.pLTable)
		assert.Equal(t, err, test.err)
	}

}
