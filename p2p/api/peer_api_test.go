package api_test

import (
	"testing"

	"github.com/it-chain/it-chain-Engine/p2p"
	"github.com/it-chain/it-chain-Engine/p2p/api"
	"github.com/it-chain/it-chain-Engine/p2p/test/mock"
	"github.com/magiconair/properties/assert"
)

func TestPeerApi_DeliverPLTable(t *testing.T) {

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

	peerApi := SetupPeerApi()

	for testName, test := range tests {

		t.Logf("running test case %s", testName)

		assert.Equal(t, peerApi.DeliverPLTable("1"), test.err)
	}

}

func SetupPeerApi() *api.PeerApi {

	pLTableQueryService := &mock.MockPLTableQueryService{}

	pLTableQueryService.FindPeerByIdFunc = func(peerId p2p.PeerId) (p2p.Peer, error) {

		peerList := mock.MakeFakePeerList()

		if peerId.Id ==""{
			return p2p.Peer{PeerId:p2p.PeerId{Id:""}, IpAddress:""}, p2p.ErrEmptyPeerId
		}

		for _, peer := range peerList {
			if peer.PeerId == peerId {

				return peer, nil
			}
		}

		return p2p.Peer{PeerId:p2p.PeerId{Id:""}, IpAddress:""}, p2p.ErrNoMatchingPeerId
	}


	pLTableQueryService.GetPLTableFunc = func() (p2p.PLTable, error) {

		return mock.MakeFakePLTable(), nil
	}

	communicationService := &mock.MockCommunicationService{}

	communicationService.DeliverPLTableFunc = func(connectionId string, pLTable p2p.PLTable) error {
		return nil
	}

	peerApi := api.NewPeerApi(pLTableQueryService, communicationService)

	return peerApi
}
