package api_test

import (
	"testing"
	"github.com/it-chain/it-chain-Engine/p2p"
	"github.com/it-chain/it-chain-Engine/p2p/api"
	"github.com/magiconair/properties/assert"
	"github.com/it-chain/midgard"
)

//todo make leader api test
//todo make fake dependencies 1. eventRepository 2. messageDispatcher 3. readonlyleaderRepository
//todo make test map
//todo test continue
type MockMessageService struct{}
func (mms MockMessageService) DeliverLeaderInfo(nodeId p2p.NodeId, leader p2p.Leader) error{return nil}

type MockEventRepository struct {}
func (mer MockEventRepository) Save(aggregateID string, events ...midgard.Event) error{return nil}

type MockReadOnlyLeaderRepository struct {}
func (mrolr MockReadOnlyLeaderRepository) GetLeader() p2p.Leader{return p2p.Leader{}}

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
		input p2p.NodeId
		err   error
	}{
		"proper node id test": {
			input: p2p.NodeId{
				Id: "",
			},
			err: api.ErrEmptyNodeId,
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

	messageService := MockMessageService{}
	leaderApi := api.NewLeaderApi(leaderRepository, eventRepository, messageService, &p2p.Node{NodeId:p2p.NodeId{Id:"123"}})

	return leaderApi
}
