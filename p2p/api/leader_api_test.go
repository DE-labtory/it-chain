package api_test

import (
	"testing"

	"github.com/it-chain/it-chain-Engine/p2p"
	"github.com/it-chain/it-chain-Engine/p2p/api"
	"github.com/it-chain/it-chain-Engine/p2p/test/mock"
	"github.com/magiconair/properties/assert"
)

var MockPLTable = mock.MakeFakePLTable()

func TestLeaderApi_UpdateLeaderWithAddress(t *testing.T) {

	tests := map[string]struct {
		input struct {
			ipAddress string
		}
		output struct {
			leader p2p.Leader
		}
		err error
	}{
		"success": {
			input: struct{ ipAddress string }{ipAddress: "2"},
			output: struct{ leader p2p.Leader }{leader: p2p.Leader{
				LeaderId: p2p.LeaderId{
					Id: "2",
				},
			}},
			err: nil,
		},
		"no matching ipAddress test": {
			input: struct{ ipAddress string }{ipAddress: "234"},
			output: struct{ leader p2p.Leader }{
				leader: p2p.Leader{
					LeaderId: p2p.LeaderId{
						Id: "1",
					},
				},
			},
			err: api.ErrNoMatchingPeerWithIpAddress,
		},
	}

	leaderApi := SetupLeaderApi()

	for testName, test := range tests {

		t.Logf("running test case %s", testName)

		err := leaderApi.UpdateLeaderWithAddress(test.input.ipAddress)

		//assert.Equal(t, mock.MakeFakePLTable().Leader, test.output.leader)

		assert.Equal(t, err, test.err)
	}
}

func TestLeaderApi_UpdateLeaderWithLargePeerTable(t *testing.T) {

	tests := map[string]struct {
		input struct {
			pLTable p2p.PLTable
		}
		output struct {
			leader p2p.Leader
		}
	}{
		"success": {
			input: struct{ pLTable p2p.PLTable }{pLTable: p2p.PLTable{
				Leader:p2p.Leader{
					LeaderId:p2p.LeaderId{
						Id:"2",
					},
				},
				PeerTable:map[string]p2p.Peer{
					"1":{
						PeerId: p2p.PeerId{
							Id: "1",
						},
					},
					"2":{
						PeerId: p2p.PeerId{
							Id: "1",
						},
					},
					"3":{
						PeerId: p2p.PeerId{
							Id: "1",
						},
					},
					"4":{
						PeerId: p2p.PeerId{
							Id: "1",
						},
					},
				},
			}},
			output: struct{ leader p2p.Leader }{leader: p2p.Leader{
				LeaderId: p2p.LeaderId{
					Id: "2",
				},
			}},
		},
		"not updated with longer peer list case": {
			input: struct{ pLTable p2p.PLTable }{pLTable: p2p.PLTable{
				Leader:p2p.Leader{
					LeaderId:p2p.LeaderId{
						Id:"1",
					},
				},
				PeerTable:map[string]p2p.Peer{
					"1":p2p.Peer{
						PeerId:p2p.PeerId{
							Id:"1",
						},
					},
					"2":p2p.Peer{
						PeerId:p2p.PeerId{
							Id:"1",
						},
					},
				},
			}},
			output: struct {
				leader p2p.Leader
			}{
				leader: p2p.Leader{
					LeaderId: p2p.LeaderId{
						Id: "1",
					},
				},
			},
		},
	}

	leaderApi := SetupLeaderApi()

	for testName, test := range tests {
		t.Logf("running test case %s ", testName)

		leaderApi.UpdateLeaderWithLargePeerTable(test.input.pLTable)

		t.Logf("%s", MockPLTable.Leader.LeaderId.Id)

		assert.Equal(t, MockPLTable.Leader, test.output.leader)
	}

}

func SetupLeaderApi() api.LeaderApi {

	leaderService := &mock.MockLeaderService{}

	leaderService.SetFunc = func(leader p2p.Leader) error {

		MockPLTable.Leader = leader

		return nil
	}

	mockPeerQueryService := &mock.MockPeerQueryService{}

	mockPeerQueryService.GetPLTableFunc = func() (p2p.PLTable, error) {


		return mock.MakeFakePLTable(), nil
	}

	leaderApi := api.NewLeaderApi(leaderService, mockPeerQueryService)

	return leaderApi
}
