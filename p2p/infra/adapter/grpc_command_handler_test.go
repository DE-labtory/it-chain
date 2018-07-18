package adapter_test

import (
	"encoding/json"
	"testing"

	"github.com/it-chain/it-chain-Engine/p2p"
	"github.com/it-chain/it-chain-Engine/p2p/infra/adapter"
	"github.com/it-chain/midgard"
	"github.com/magiconair/properties/assert"
	"github.com/it-chain/it-chain-Engine/p2p/test/mock"
)


func TestGrpcCommandHandler_HandleMessageReceive(t *testing.T) {

	leader := p2p.Leader{}
	leaderByte, _ := json.Marshal(leader)

	//todo error case write!
	tests := map[string]struct {
		input struct {
			id       string
			protocol string
			body     []byte
		}
		err error
	}{
		"leader info deliver test success": {
			input: struct {
				id       string
				protocol string
				body     []byte
			}{
				id:       "1",
				protocol: string("LeaderInfoDeliverProtocol"),
				body:     leaderByte,
			},
			err: nil,
		},
		"leader info deliver test empty leader id":       {},
		"leader table deliver test success":              {},
		"peer leader table deliver test empty peer list": {},
		"peer leader table deliver test empty leader id": {},
	}

	leaderApi := &mock.MockLeaderApi{}

	electionService := p2p.ElectionService{}

	communicationApi := &mock.MockCommunicationApi{}

	pLTableService := &mock.MockPLTableService{}

	messageHandler := adapter.NewGrpcCommandHandler(leaderApi, electionService, communicationApi, pLTableService)

	for testName, test := range tests {
		grpcReceiveCommand := p2p.GrpcReceiveCommand{
			CommandModel: midgard.CommandModel{
				ID: test.input.id,
			},
			Body:     test.input.body,
			Protocol: test.input.protocol,
		}
		t.Logf("running test case %s", testName)
		err := messageHandler.HandleMessageReceive(grpcReceiveCommand)
		assert.Equal(t, err, test.err)
	}

}
