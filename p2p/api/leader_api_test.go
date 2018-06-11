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
	"errors"
)

func TestLeaderApi_UpdateLeader(t *testing.T) {
	path := "test"
	url := "test"
	event := p2p.LeaderUpdatedEvent{}
	store := leveldb.NewEventStore(path, leveldb.NewSerializer(event))
	client := rabbitmq.Connect(url)

	eventRepository := midgard.NewRepo(store, client)

	publisher := messaging.Publisher(func(exchange string, topic string, data interface{}) error{
		assert.Equal(t, exchange, "Event")
		assert.Equal(t, topic, "leader.update")
		assert.Equal(t, reflect.TypeOf(data).String(), "LeaderUpdatedEvent")
	})
	messageDispatcher := messaging.NewMessageDispatcher(publisher)

	nodeId := p2p.NodeId{
		Id: "777",
	}

	myInfo := p2p.NewNode(conf.GetConfiguration().Common.NodeIp, nodeId)
	leaderApi := NewLeaderApi(*eventRepository, messageDispatcher, myInfo)

	leader := p2p.Leader{
		LeaderId:p2p.LeaderId{
			Id:"777",
		},
	}
	err := leaderApi.UpdateLeader(leader)
	assert.Equal(t, err, errors.New("empty leader id purposed"))
}

func TestLeaderApi_DeliverLeaderInfo(t *testing.T) {

	path := "test"
	url := "test"
	event := p2p.LeaderUpdatedEvent{}
	store := leveldb.NewEventStore(path, leveldb.NewSerializer(event))
	client := rabbitmq.Connect(url)

	eventRepository := midgard.NewRepo(store, client)

	publisher := messaging.Publisher(func(exchange string, topic string, data interface{}) error{
		assert.Equal(t, exchange, "Command")
		assert.Equal(t, topic, "message.deliver")
		assert.Equal(t, reflect.TypeOf(data).String(), "LeaderInfoDelivery")
	})
	messageDispatcher := messaging.NewMessageDispatcher(publisher)

	nodeId := p2p.NodeId{
		Id: "777",
	}

	myInfo := p2p.NewNode(conf.GetConfiguration().Common.NodeIp, nodeId)
	leaderApi := NewLeaderApi(*eventRepository, messageDispatcher, myInfo)

	err := leaderApi.DeliverLeaderInfo(p2p.NodeId{Id:"777"})
	assert.Equal(t, err, errors.New("empty node id purposed"))
}
