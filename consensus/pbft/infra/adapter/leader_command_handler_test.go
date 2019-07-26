package adapter_test

import (
	"reflect"
	"testing"

	"github.com/DE-labtory/it-chain/common"
	"github.com/DE-labtory/it-chain/common/command"
	"github.com/DE-labtory/it-chain/consensus/pbft"
	"github.com/DE-labtory/it-chain/consensus/pbft/infra/adapter"
	"github.com/DE-labtory/it-chain/consensus/pbft/test/mock"
	"github.com/stretchr/testify/assert"
)

func TestNewLeaderCommandHandler(t *testing.T) {
	parliamentApi := &mock.ParliamentApi{}
	leaderCommandHandler := adapter.NewLeaderCommandHandler(parliamentApi)

	assert.Equal(t, "*adapter.LeaderCommandHandler", reflect.TypeOf(leaderCommandHandler).String())
}

func TestLeaderCommandHandler_HandleMessageReceive_RequestLeaderProtocol(t *testing.T) {
	// given
	Command := command.ReceiveGrpc{
		MessageId:    "WhoAreYou",
		Body:         []byte{1, 2, 3},
		ConnectionID: "IWantToConnectWithYou",
		Protocol:     "RequestLeaderProtocol",
	}
	// given
	parliamentApi := &mock.ParliamentApi{}
	parliamentApi.DeliverLeaderFunc = func(connectionId string) {
		// then
		assert.Equal(t, Command.ConnectionID, connectionId)
	}

	// when
	leaderCommandHandler := adapter.NewLeaderCommandHandler(parliamentApi)
	// when
	err := leaderCommandHandler.HandleMessageReceive(Command)
	// then
	assert.NoError(t, err)
}

func TestLeaderCommandHandler_HandleMessageReceive_LeaderDeliveryProtocol(t *testing.T) {
	// given
	leaderDeliveryMsg := &pbft.LeaderDeliveryMessage{
		Leader: pbft.Leader{
			LeaderId: "IAmCreativeLeader",
		},
	}
	data, _ := common.Serialize(leaderDeliveryMsg)

	// given
	Command := command.ReceiveGrpc{
		MessageId:    "WhoAreYou",
		Body:         data,
		ConnectionID: "IWantToConnectWithYou",
		Protocol:     "LeaderDeliveryProtocol",
	}

	// given
	parliamentApi := &mock.ParliamentApi{}
	parliamentApi.UpdateLeaderFunc = func(nodeId common.NodeID) error {
		// then
		assert.Equal(t, leaderDeliveryMsg.Leader.LeaderId, nodeId)
		return nil
	}

	// when
	leaderCommandHandler := adapter.NewLeaderCommandHandler(parliamentApi)
	// when
	err := leaderCommandHandler.HandleMessageReceive(Command)
	// then
	assert.NoError(t, err)
}

func TestLeaderCommandHandler_HandleMessageReceive_LeaderDeliveryProtocol_DeserializeErr(t *testing.T) {
	// given
	Command := command.ReceiveGrpc{
		MessageId:    "WhoAreYou",
		Body:         []byte("StupidBytes"),
		ConnectionID: "IWantToConnectWithYou",
		Protocol:     "LeaderDeliveryProtocol",
	}

	// given
	parliamentApi := &mock.ParliamentApi{}

	// when
	leaderCommandHandler := adapter.NewLeaderCommandHandler(parliamentApi)
	// when
	assert.Panics(t, func() { leaderCommandHandler.HandleMessageReceive(Command) })
}
