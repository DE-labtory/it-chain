package adapter_test

import (
	"encoding/json"
	"testing"

	"github.com/it-chain/it-chain-Engine/p2p"
	"github.com/it-chain/it-chain-Engine/p2p/infra/adapter"
	"github.com/it-chain/midgard"
	"github.com/magiconair/properties/assert"
)

type MockLeaderApi struct{}

func (mla MockLeaderApi) UpdateLeader(leader p2p.Leader) error { return nil }
func (mla MockLeaderApi) DeliverLeaderInfo(connectionId string)  {}

type MockPeerApi struct{}

func (mna MockPeerApi) UpdatePeerList(peerList []p2p.Peer) error { return nil }
func (mna MockPeerApi) DeliverPeerList(connectionId string) error  { return nil }
func (mna MockPeerApi) AddPeer(peer p2p.Peer)                    {}

func TestGrpcMessageHandler_HandleMessageReceive(t *testing.T) {

	leader := p2p.Leader{}
	leaderByte, _ := json.Marshal(leader)

	peerList := make([]p2p.Peer, 0)
	newPeerList := append(peerList, p2p.Peer{PeerId: p2p.PeerId{Id: "123"}})
	peerListByte, _ := json.Marshal(newPeerList)

	//todo error case write!
	tests := map[string]struct {
		input struct {
			command p2p.GrpcReceiveCommand
		}
		err error
	}{
		"leader info deliver test success": {
			input: struct {
				command p2p.GrpcReceiveCommand
			}{
				command: p2p.GrpcReceiveCommand{
					CommandModel: midgard.CommandModel{ID: "123"},
					Body:         leaderByte,
					Protocol:     "LeaderInfoDeliverProtocol",
				},
			},
			err: nil,
		},
		"peer list deliver test success": {
			input: struct {
				command p2p.GrpcReceiveCommand
			}{
				command: p2p.GrpcReceiveCommand{
					CommandModel: midgard.CommandModel{ID: "123"},
					Body:         peerListByte,
					Protocol:     "PeerListDeliverProtocol",
				},
			},
			err: nil,
		},
	}
	leaderApi := MockLeaderApi{}
	peerApi := MockPeerApi{}
	messageHandler := adapter.NewGrpcCommandHandler(leaderApi, peerApi)

	for testName, test := range tests {
		t.Logf("running test case %s", testName)
		err := messageHandler.HandleMessageReceive(test.input.command)
		assert.Equal(t, err, test.err)
	}

}
