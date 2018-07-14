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

func TestLeaderApi_UpdateLeaderWithLongerPeerList(t *testing.T) {

	tests := map[string]struct {
		input struct {
			peerList []p2p.Peer
			leader   p2p.Leader
		}
		output struct {
			leader p2p.Leader
		}
	}{
		"success": {
			input: struct {
				peerList []p2p.Peer
				leader   p2p.Leader
			}{peerList: []p2p.Peer{
				{
					PeerId: p2p.PeerId{
						Id: "1",
					},
				}, {
					PeerId: p2p.PeerId{
						Id: "1",
					},
				}, {
					PeerId: p2p.PeerId{
						Id: "1",
					},
				}, {
					PeerId: p2p.PeerId{
						Id: "1",
					},
				},
			}, leader: p2p.Leader{
				LeaderId: p2p.LeaderId{
					Id: "2",
				},
			}},
			output: struct{ leader p2p.Leader }{leader: p2p.Leader{
				LeaderId: p2p.LeaderId{
					Id: "2",
				},
			}},
		},
		"not updated with longer peer list case": {
			input: struct {
				peerList []p2p.Peer
				leader   p2p.Leader
			}{
				peerList: []p2p.Peer{
					{
						PeerId: p2p.PeerId{
							Id: "1",
						},
					},
				}, leader: p2p.Leader{
					LeaderId: p2p.LeaderId{
						Id: "2",
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

		leaderApi.UpdateLeaderWithLongerPeerList(test.input.leader, test.input.peerList)

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

	mockPLTableQueryService := &mock.MockPLTableQueryService{}

	mockPLTableQueryService.GetPLTableFunc = func() (p2p.PLTable, error) {


		return mock.MakeFakePLTable(), nil
	}

	leaderApi := api.NewLeaderApi(leaderService, mockPLTableQueryService)

	return leaderApi
}
