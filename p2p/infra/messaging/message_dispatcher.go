package messaging

import (
	"time"

	"github.com/it-chain/it-chain-Engine/common"
	"github.com/it-chain/it-chain-Engine/messaging/rabbitmq/event"
	"github.com/it-chain/it-chain-Engine/p2p"
	"github.com/it-chain/midgard"
	"github.com/rs/xid"
)

type MessageDispatcher struct {
	publisher midgard.Publisher
}

func NewMessageDispatcher(publisher midgard.Publisher) *MessageDispatcher {
	return &MessageDispatcher{
		publisher: publisher,
	}
}

func (md *MessageDispatcher) RequestLeaderInfo(nodeId p2p.NodeId) error {
	body := p2p.LeaderInfoRequestMessage{
		TimeUnix: time.Now().Unix(),
	}

	deliverCommand, err := CreateMessageDeliverCommand(event.LeaderInfoDeliverProtocol, body)


	if err != nil {
		return err
	}

	deliverCommand.Recipients = append(deliverCommand.Recipients, nodeId.ToString())

	return md.publisher.Publish("Command", "message.deliver", deliverCommand)
}

// command message which requests node list of specific node
func (md *MessageDispatcher) RequestNodeList(nodeId p2p.NodeId) error {

	body := p2p.LeaderInfoRequestMessage{
		TimeUnix: time.Now().Unix(),
	}

	deliverCommand, err := CreateMessageDeliverCommand("NodeListRequestMessage", body)

	if err != nil {
		return err
	}

	deliverCommand.Recipients = append(deliverCommand.Recipients, nodeId.ToString())

	return md.publisher.Publish("Commnand", "message.deliver", deliverCommand)
}

func (md *MessageDispatcher) DeliverLeaderInfo(nodeId p2p.NodeId, leader p2p.Node) error {

	deliverCommand, err := CreateMessageDeliverCommand("UpdateLeader", leader)

	if err != nil {
		return err
	}

	deliverCommand.Recipients = append(deliverCommand.Recipients, nodeId.ToString())

	return md.publisher.Publish("Command", "message.deliver", deliverCommand)
}

func (md *MessageDispatcher) DeliverNodeList(nodeId p2p.NodeId, nodeList []p2p.Node) error {
	deliverCommand, err := CreateMessageDeliverCommand("NodeListDeliver", nodeList)

	if err != nil {
		return err
	}

	deliverCommand.Recipients = append(deliverCommand.Recipients, nodeId.ToString())


	return md.publisher.Publish("Command", "message.deliver", deliverCommand)
}

func CreateMessageDeliverCommand(protocol string, body interface{}) (p2p.MessageDeliverCommand, error) {

	data, err := common.Serialize(body)

	if err != nil {
		return p2p.MessageDeliverCommand{}, err
	}

	return p2p.MessageDeliverCommand{
		CommandModel: midgard.CommandModel{
			ID: xid.New().String(),
		},
		Recipients: make([]string, 0),
		Body:       data,
		Protocol:   protocol,
	}, err
}