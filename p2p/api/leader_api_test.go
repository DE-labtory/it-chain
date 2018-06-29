package api_test

import (
	"testing"

	"github.com/it-chain/it-chain-Engine/p2p"
	"github.com/it-chain/it-chain-Engine/p2p/api"
	"github.com/it-chain/midgard"
	"github.com/magiconair/properties/assert"
)

type MockGrpcCommandService struct{}

func (mms MockGrpcCommandService) DeliverLeaderInfo(connectionId string, leader p2p.Leader) error {
	return nil
}

type MockEventRepository struct{}

func (mer MockEventRepository) Save(aggregateID string, events ...midgard.Event) error { return nil }

type MockReadOnlyLeaderRepository struct{}

func (mrolr MockReadOnlyLeaderRepository) GetLeader() p2p.Leader { return p2p.Leader{} }

func TestLeaderApi_UpdateLeader(t *testing.T) {

	tests := map[string]struct {
		input p2p.Leader
		err   error
	}{
		"empty leader id test": {
			input: p2p.Leader{
				LeaderId: p2p.LeaderId{
					Id: "",
				},
			},
			err: api.ErrEmptyLeaderId,
		},
		"first leader update test": {
			input: p2p.Leader{
				LeaderId: p2p.LeaderId{
					Id: "1",
				},
			},
			err: nil,
		},
	}

	leaderApi := SetupLeaderApi()

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)
		err := leaderApi.UpdateLeader(test.input)
		assert.Equal(t, err, test.err)
	}
}

func TestLeaderApi_DeliverLeaderInfo(t *testing.T) {
	tests := map[string]struct {
		input string
		err   error
	}{
		"proper node id test": {
			input: "",
			err: api.ErrEmptyConnectionId,
		},
	}
	leaderApi := SetupLeaderApi()

	for testName, test := range tests {
		t.Logf("running test case %s", testName)
		err := leaderApi.DeliverLeaderInfo(test.input)
		assert.Equal(t, err, test.err)
	}
}

func SetupLeaderApi() *api.LeaderApi {

	leaderRepository := MockReadOnlyLeaderRepository{}
	eventRepository := MockEventRepository{}

	grpcCommandService := MockGrpcCommandService{}
	leaderApi := api.NewLeaderApi(leaderRepository, eventRepository, grpcCommandService, &p2p.Peer{PeerId: p2p.PeerId{Id: "123"}})

	return leaderApi
}
