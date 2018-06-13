package api

import (
	"testing"

	"github.com/it-chain/it-chain-Engine/p2p/infra/messaging"
	"github.com/it-chain/it-chain-Engine/p2p"
	"github.com/it-chain/it-chain-Engine/conf"
	"github.com/magiconair/properties/assert"
	"reflect"
	"github.com/it-chain/midgard/bus/rabbitmq"
	"github.com/it-chain/midgard/store/leveldb"
	"github.com/it-chain/midgard"
	leveldb2 "github.com/it-chain/it-chain-Engine/p2p/infra/repository/leveldb"
)

// make leader api for use
// make fake dependencies 1. eventRepository 2. messageDispatcher 3. myInfo
// make test map
// test continue
func TestLeaderApi_UpdateLeader(t *testing.T) {

	tests := map[string]struct {
		input  p2p.Leader
		err    error
	}{
		"empty leader id test":{
			input:p2p.Leader{
				LeaderId:p2p.LeaderId{
					Id:"",
				},
			},
			err:ErrEmptyLeaderId,
		},
		"first leader update test": {
			input:  p2p.Leader{
				LeaderId:p2p.LeaderId{
					Id: "1",
				},
			},
			err:    nil,
		},
		"update with same leader": {
			input:  p2p.Leader{
				LeaderId:p2p.LeaderId{
					Id: "1",
				},
			},
			err:    ErrSameLeader,
		},
	}

	publisher := messaging.Publisher(func(exchange string, topic string, data interface{}) error{})
	myInfo := p2p.NewNode(conf.GetConfiguration().Common.NodeIp, p2p.NodeId{Id:"1"})
	event := p2p.LeaderUpdatedEvent{}
	leaderApi := SetupLeaderApi("path", "path", "path", event, publisher, myInfo)

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)
		err := leaderApi.UpdateLeader(test.input)
		assert.Equal(t, err, test.err)
	}
}

func TestLeaderApi_DeliverLeaderInfo(t *testing.T) {
	tests := map[string]struct{
		input p2p.NodeId
		err error
	}{
		"proper node id test":{
			input:p2p.NodeId{
				Id:"",
			},
			err:ErrEmptyNodeId,
		},
	}
	publisher := messaging.Publisher(func(exchange string, topic string, data interface{}) error{
		assert.Equal(t, exchange, "Command")
		assert.Equal(t, topic, "message.deliver")
		assert.Equal(t, reflect.TypeOf(data).String(), "LeaderInfoDelivery")

		return nil
	})
	myInfo := p2p.NewNode(conf.GetConfiguration().Common.NodeIp, p2p.NodeId{Id:"1"})
	event := p2p.LeaderDeliveredEvent{}
	leaderApi := SetupLeaderApi("test", "test", "", event, publisher, myInfo)

	for testName, test := range tests{
		t.Logf("running test case %s", testName)
		err := leaderApi.DeliverLeaderInfo(test.input)
		assert.Equal(t, err, test.err)
	}
}

func SetupLeaderApi(leader_repo_path string, event_repo_path string, url string, event midgard.Event, publisher messaging.Publisher, myInfo *p2p.Node) *LeaderApi{

	store := leveldb.NewEventStore(event_repo_path, leveldb.NewSerializer(event))
	client := rabbitmq.Connect(url)
	leaderRepository := leveldb2.NewLeaderRepository(leader_repo_path)
	eventRepository := midgard.NewRepo(store, client)


	messageDispatcher := messaging.NewMessageDispatcher(publisher)

	leaderApi := NewLeaderApi(leaderRepository, *eventRepository, messageDispatcher, myInfo)

	return leaderApi
}