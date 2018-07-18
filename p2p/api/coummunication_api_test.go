package api_test

import (
	"testing"

	"github.com/it-chain/it-chain-Engine/p2p"
	"github.com/it-chain/it-chain-Engine/p2p/api"
	"github.com/magiconair/properties/assert"
	"github.com/it-chain/it-chain-Engine/p2p/test/mock"
)

func TestCommunicationApi_DialToUnConnectedNode(t *testing.T) {
	tests := map[string]struct {
		input struct {
			peerTable map[string]p2p.Peer
		}
		err error
	}{
		"success": {
			input: struct{ peerTable map[string]p2p.Peer }{peerTable: map[string]p2p.Peer{
				"1":{
					PeerId:p2p.PeerId{
						Id:"1",
					},
				},
			}},
			err: nil,
		},
	}

	mockPeerQueryService :=&mock.MockPeerQueryService{}
	mockPeerQueryService.FindPeerByIdFunc = func(peerId p2p.PeerId) (p2p.Peer, error) {

		peerTable := mock.MakeFakePeerTable()

		for _, peer := range peerTable {
			if peer.PeerId == peerId {
				return peer, nil
			}
		}

		return p2p.Peer{}, nil
	}

	for testName, test := range tests {

		t.Logf("running test case %s", testName)

		communicationApi := api.NewCommunicationApi(mockPeerQueryService, &mock.MockCommunicationService{})

		assert.Equal(t, communicationApi.DialToUnConnectedNode(test.input.peerTable), test.err)

	}
}


func TestCommunicationApi_DeliverPLTable(t *testing.T) {

	tests := map[string]struct {
		input struct {
			pLTable p2p.PLTable
		}
		err error
	}{
		"success": {
			input: struct{ pLTable p2p.PLTable }{pLTable: p2p.PLTable{
				Leader: p2p.Leader{
					LeaderId: p2p.LeaderId{
						Id: "1",
					},
				},
				PeerTable:map[string]p2p.Peer{
					"1":p2p.Peer{
						PeerId:p2p.PeerId{Id:"1"},
					},
					"2":p2p.Peer{
						PeerId:p2p.PeerId{Id:"2"},
					},
				},
			}},
			err: nil,
		},
	}

	communicationApi := SetupCommunicationApi()

	for testName, test := range tests {

		t.Logf("running test case %s", testName)

		assert.Equal(t, communicationApi.DeliverPLTable("1"), test.err)
	}

}

func SetupCommunicationApi() *api.CommunicationApi {

	peerQueryService := &mock.MockPeerQueryService{}

	peerQueryService.FindPeerByIdFunc = func(peerId p2p.PeerId) (p2p.Peer, error) {

		pLTable := mock.MakeFakePLTable()

		if peerId.Id ==""{
			return p2p.Peer{PeerId:p2p.PeerId{Id:""}, IpAddress:""}, p2p.ErrEmptyPeerId
		}

		for _, peer := range pLTable.PeerTable {
			if peer.PeerId == peerId {

				return peer, nil
			}
		}

		return p2p.Peer{PeerId:p2p.PeerId{Id:""}, IpAddress:""}, p2p.ErrNoMatchingPeerId
	}


	peerQueryService.GetPLTableFunc = func() (p2p.PLTable, error) {

		return mock.MakeFakePLTable(), nil
	}

	communicationService := &mock.MockCommunicationService{}

	communicationService.DeliverPLTableFunc = func(connectionId string, pLTable p2p.PLTable) error {
		return nil
	}

	communicationApi := api.NewCommunicationApi(peerQueryService, communicationService)

	return communicationApi
}