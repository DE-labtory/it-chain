package adapter

import (
	"time"

	"errors"

	"github.com/it-chain/it-chain-Engine/common"
	"github.com/it-chain/it-chain-Engine/messaging/rabbitmq/event"
	"github.com/it-chain/it-chain-Engine/p2p"
	"github.com/it-chain/midgard"
	"github.com/rs/xid"
)

//kind of error
var ErrEmptyNodeId = errors.New("empty nodeid proposed")
var ErrEmptyLeaderId = errors.New("empty leader id proposed")
var ErrEmptyNodeList = errors.New("empty node list proposed")


type Publish func(exchange string, topic string, data interface{}) (err error) // 나중에 의존성 주입을 해준다.

// message dispatcher sends messages to other nodes in p2p network
type MessageService struct {
	publish Publish // midgard.client.Publish
}

func NewMessageService(publish Publish) *MessageService {
	return &MessageService{
		publish: publish,
	}
}

//request leader information in p2p network to the node specified by nodeId
func (md *MessageService) RequestLeaderInfo(nodeId p2p.NodeId) error {

	if nodeId.Id == "" {
		return ErrEmptyNodeId
	}

	body := p2p.LeaderInfoRequestMessage{
		TimeUnix: time.Now().Unix(),
	}

	//message deliver command for delivering leader info
	deliverCommand, err := CreateMessageDeliverCommand(event.LeaderInfoDeliverProtocol, body)

	if err != nil {
		return err
	}

	deliverCommand.Recipients = append(deliverCommand.Recipients, nodeId.ToString())

	return md.publish("Command", "message.deliver", deliverCommand)
}

// command message which requests node list of specific node
func (md *MessageService) RequestNodeList(nodeId p2p.NodeId) error {

	if nodeId.Id == "" {
		return ErrEmptyNodeId
	}
	body := p2p.LeaderInfoRequestMessage{
		TimeUnix: time.Now().Unix(),
	}

	deliverCommand, err := CreateMessageDeliverCommand("NodeListRequestMessage", body)

	if err != nil {
		return err
	}

	deliverCommand.Recipients = append(deliverCommand.Recipients, nodeId.ToString())

	return md.publish("Command", "message.deliver", deliverCommand)
}

func (md *MessageService) DeliverLeaderInfo(nodeId p2p.NodeId, leader p2p.Leader) error {

	if nodeId.Id == "" {
		return ErrEmptyNodeId
	}

	if leader.LeaderId.Id == "" {
		return ErrEmptyLeaderId
	}

	deliverCommand, err := CreateMessageDeliverCommand("UpdateLeader", leader)

	if err != nil {
		return err
	}

	deliverCommand.Recipients = append(deliverCommand.Recipients, nodeId.ToString())

	return md.publish("Command", "message.deliver", deliverCommand)
}

//deliver node list to other node specified by nodeId
func (md *MessageService) DeliverNodeList(nodeId p2p.NodeId, nodeList []p2p.Node) error {

	if nodeId.Id == "" {
		return ErrEmptyNodeId
	}

	if len(nodeList) == 0 {
		return ErrEmptyNodeList
	}

	messageDeliverCommand, err := CreateMessageDeliverCommand("NodeListDeliver", nodeList)

	if err != nil {
		return err
	}

	messageDeliverCommand.Recipients = append(messageDeliverCommand.Recipients, nodeId.ToString())

	return md.publish("Command", "message.deliver", messageDeliverCommand)
}

//deliver single node
func (md *MessageService) DeliverNode(nodeId p2p.NodeId, node p2p.Node) error {

	messageDeliverCommand, err := CreateMessageDeliverCommand("NodeDeliverProtocol", node)

	if err != nil {
		return err
	}

	messageDeliverCommand.Recipients = append(messageDeliverCommand.Recipients, node.NodeId.ToString())

	return md.publish("Command", "message.deliver", messageDeliverCommand)
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
