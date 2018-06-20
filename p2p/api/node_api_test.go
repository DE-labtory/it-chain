package api_test

import (
	"testing"
	"github.com/it-chain/it-chain-Engine/p2p"
	"errors"
	"github.com/it-chain/it-chain-Engine/p2p/api"
	"github.com/magiconair/properties/assert"
)

var ErrEmptyNodeList = errors.New("empty node list proposed")

//todo make node api test
//todo make fake dependencies 1. eventRepository 2. messageDispatcher 3. nodeRepository
//todo make test map
//todo test continue

type MockNodeRepository struct {}
func (mnr MockNodeRepository)FindById(id p2p.NodeId) (*p2p.Node, error){return nil, nil}
func (mnr MockNodeRepository)FindAll() ([]p2p.Node, error){return nil, nil}

type MockNodeMessageService struct {}
func (mnms MockNodeMessageService) DeliverNodeList(nodeId p2p.NodeId, nodeList []p2p.Node) error{return nil}

func TestNodeApi_UpdateNodeList(t *testing.T) {

	tests := map[string]struct {
		input []p2p.Node
		err   error
	}{
		"success":{
			input: []p2p.Node{
				p2p.Node{
					NodeId:p2p.NodeId{
						Id:"1",
					},
				},
			},
			err: nil,
		},
	}

	nodeApi := SetupNodeApi()

	for testName, test := range tests {
		t.Logf("running test case %s", testName)
		err := nodeApi.UpdateNodeList(test.input)
		assert.Equal(t, err, test.err)
	}
}

func SetupNodeApi() *api.NodeApi {
	nodeRepository := MockNodeRepository{}
	eventRepository := MockEventRepository{}
	messageService := MockNodeMessageService{}

	nodeApi := api.NewNodeApi(nodeRepository, eventRepository, messageService)

	return nodeApi
}
