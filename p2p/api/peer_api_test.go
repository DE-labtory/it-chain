package api_test

import (
	"testing"

	"github.com/it-chain/it-chain-Engine/p2p"
	"github.com/it-chain/it-chain-Engine/p2p/api"
	"github.com/it-chain/it-chain-Engine/p2p/test/mock"
	"github.com/magiconair/properties/assert"
)

func TestPeerApi_UpdatePeerList(t *testing.T) {

	tests := map[string]struct {
		input []p2p.Peer
		err   error
	}{
		"success": {
			input: []p2p.Peer{
				{
					PeerId: p2p.PeerId{
						Id: "1",
					},
					IpAddress: "1",
				},
				{
					PeerId: p2p.PeerId{
						Id: "2",
					},
					IpAddress: "2",
				},
				{
					PeerId: p2p.PeerId{
						Id: "3",
					},
					IpAddress: "3",
				},
			},
			err: nil,
		},
	}

	peerApi := SetupPeerApi()

	for testName, test := range tests {

		t.Logf("running test case %s", testName)

		err := peerApi.UpdatePeerList(test.input)

		assert.Equal(t, err, test.err)
	}
}

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
				PeerList: []p2p.Peer{
					{
						PeerId: p2p.PeerId{
							Id: "1",
						},
						IpAddress: "1",
					},
					{
						PeerId: p2p.PeerId{
							Id: "2",
						},
						IpAddress: "2",
					},
					{
						PeerId: p2p.PeerId{
							Id: "3",
						},
						IpAddress: "3",
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

func TestPeerApi_FindById(t *testing.T) {

	tests := map[string]struct {
		input struct {
			peerId p2p.PeerId
		}
		output struct{
			peer p2p.Peer
		}
		err error
	}{
		"success": {
			input: struct{ peerId p2p.PeerId }{peerId: p2p.PeerId{Id: "1"}},
			output: struct{ peer p2p.Peer }{peer: p2p.Peer{IpAddress: "1", PeerId: p2p.PeerId{Id:"1",}},},
			err:   nil,
		},
		"no matching peer id test":{
			input: struct{ peerId p2p.PeerId }{peerId: p2p.PeerId{Id:"asdfadsf"}},
			output: struct{ peer p2p.Peer }{peer: struct {
				IpAddress string
				PeerId    p2p.PeerId
			}{IpAddress: string(""), PeerId: struct{ Id string }{Id: string("")}}},
			err:p2p.ErrNoMatchingPeerId,

		},
		"empty peer id proposed test":{
			input: struct{ peerId p2p.PeerId }{peerId: struct{ Id string }{Id: string("")}},
			output: struct{ peer p2p.Peer }{peer: struct {
				IpAddress string
				PeerId    p2p.PeerId
			}{IpAddress: string(""), PeerId: struct{ Id string }{Id: string("")}}},
			err:p2p.ErrEmptyPeerId,

		},
	}

	peerApi := SetupPeerApi()

	for testName, test := range tests{

		t.Logf("running test case %s", testName)

		peer, err := peerApi.FindById(test.input.peerId)

		assert.Equal(t, peer, test.output.peer)
		assert.Equal(t, err, test.err)
		}
}

func SetupPeerApi() *api.PeerApi {

	peerQueryService := &mock.MockPeerQueryService{}

	peerQueryService.FindAllFunc = func() ([]p2p.Peer, error) {

		return mock.MakeFakePeerList(), nil
	}

	peerQueryService.FindByIdFunc = func(peerId p2p.PeerId) (p2p.Peer, error) {

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

	pLTableQueryService := &mock.MockPLTableQueryService{}

	pLTableQueryService.GetPLTableFunc = func() (p2p.PLTable, error) {

		return mock.MakeFakePLTable(), nil
	}

	communicationService := &mock.MockCommunicationService{}

	communicationService.DeliverPLTableFunc = func(connectionId string, pLTable p2p.PLTable) error {
		return nil
	}

	peerApi := api.NewPeerApi(peerQueryService, pLTableQueryService, communicationService)

	return peerApi
}
