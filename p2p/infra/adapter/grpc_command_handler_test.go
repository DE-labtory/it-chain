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
func (mna MockPeerApi) GetPeerLeaderTable() (p2p.PeerLeaderTable){
	peerTable := p2p.PeerLeaderTable{
		Leader:p2p.Leader{LeaderId:p2p.LeaderId{Id:"1"}},
		PeerList:[]p2p.Peer{p2p.Peer{PeerId:p2p.PeerId{Id:"2"}}},
	}
	return peerTable
}
func (mna MockPeerApi) FindById(peerId p2p.PeerId) (p2p.Peer, error){
	peer := p2p.Peer{PeerId:peerId}
	return peer, nil
}
func (mna MockPeerApi) GetPeerList() []p2p.Peer{
	peerList := []p2p.Peer{{PeerId:p2p.PeerId{Id:"2"}}}
	return peerList
}
func (mna MockPeerApi) UpdatePeerList(peerList []p2p.Peer) error { return nil }
func (mna MockPeerApi) DeliverPeerLeaderTable(connectionId string) error  { return nil }
func (mna MockPeerApi) AddPeer(peer p2p.Peer)                    {}

type MockCommandService struct{}
func (mcs MockCommandService) Dial(ipAddress string) error {return nil}

func TestGrpcCommandHandler_HandleMessageReceive(t *testing.T) {

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
		"peer table deliver test success": {
			input: struct {
				command p2p.GrpcReceiveCommand
			}{
				command: p2p.GrpcReceiveCommand{
					CommandModel: midgard.CommandModel{ID: "123"},
					Body:         peerListByte,
					Protocol:     "PeerLeaderTableDeliverProtocol",
				},
			},
			err: nil,
		},
	}
	leaderApi := MockLeaderApi{}
	peerApi := MockPeerApi{}
	commandService := MockCommandService{}

	messageHandler := adapter.NewGrpcCommandHandler(leaderApi, peerApi, commandService)

	for testName, test := range tests {
		t.Logf("running test case %s", testName)
		err := messageHandler.HandleMessageReceive(test.input.command)
		assert.Equal(t, err, test.err)
	}

}

//todo
func TestReceiverPeerLeaderTable(t *testing.T) {

}

//todo
func TestUpdateWithLongerPeerList(t *testing.T) {

}