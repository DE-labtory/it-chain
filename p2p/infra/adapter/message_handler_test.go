package adapter_test

import (
	"testing"
	"github.com/it-chain/it-chain-Engine/p2p/infra/adapter"
	"github.com/it-chain/it-chain-Engine/p2p"
	"github.com/it-chain/midgard"
	"github.com/magiconair/properties/assert"
	"encoding/json"
)
type MockLeaderApi struct {}
func (mla MockLeaderApi) UpdateLeader(leader p2p.Leader) error{return nil}
func (mla MockLeaderApi) DeliverLeaderInfo(nodeId p2p.NodeId){}

type MockNodeApi struct{}
func (mna MockNodeApi) UpdateNodeList(nodeList []p2p.Node) error{return nil}
func (mna MockNodeApi) DeliverNodeList(nodeId p2p.NodeId) error{return nil}
func (mna MockNodeApi) AddNode(node p2p.Node){}

func TestGrpcMessageHandler_HandleMessageReceive(t *testing.T) {
	leader := p2p.Leader{}
	leaderByte, _ := json.Marshal(leader)

	nodeList := make([]p2p.Node, 0)
	newNodeList := append(nodeList, p2p.Node{NodeId:p2p.NodeId{Id:"123"}})
	nodeListByte, _ := json.Marshal(newNodeList)

	tests := map[string]struct{
		input struct{
			command p2p.GrpcRequestCommand
		}
		err error
	}{
		"leader info deliver test success": {
			input: struct{
				command p2p.GrpcRequestCommand
			}{
				command: p2p.GrpcRequestCommand{
					CommandModel: midgard.CommandModel{ID:"123"},
					Data:leaderByte,
					Protocol:"LeaderInfoDeliverProtocol",
				},
			},
			err: nil,
		},
		"node list deliver test success": {
			input: struct{
				command p2p.GrpcRequestCommand
			}{
				command: p2p.GrpcRequestCommand{
					CommandModel: midgard.CommandModel{ID:"123"},
					Data:nodeListByte,
					Protocol:"NodeListDeliverProtocol",
				},
			},
			err: nil,
		},
	}
	leaderApi := MockLeaderApi{}
	nodeApi := MockNodeApi{}
	messageHandler := adapter.NewMessageHandler(leaderApi, nodeApi)

	for testName, test := range tests{
		t.Logf("running test case %s", testName)
		err := messageHandler.HandleMessageReceive(test.input.command)
		assert.Equal(t, err, test.err)
	}

}
