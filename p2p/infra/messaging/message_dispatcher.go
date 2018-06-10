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

// publish command to amqp to get leader info from other node
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

// 단일 피어에게 새로운 리더 정보를 전달하는 메서드이다.
func (md *MessageDispatcher) DeliverLeaderInfo(nodeId p2p.NodeId, leader p2p.Node) error {

	// 리더 정보를 leaderInfoBody에 담아줌
	body := p2p.LeaderInfoResponseMessage{
		LeaderId: leader.NodeId.ToString(),
		Address:  leader.IpAddress,
	}

	deliverCommand, err := CreateMessageDeliverCommand(event.LeaderInfoDeliverProtocol, body)

	if err != nil {
		return err
	}

	// 메세지를 수신할 수신자들을 지정해 준다.
	deliverCommand.Recipients = append(deliverCommand.Recipients, nodeId.ToString())

	// topic 과 serilized data를 받아 publisher 한다.
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

func (md *MessageDispatcher) ResponseTable(toNode p2p.Node, nodes []p2p.Node) error {
	panic("implement me")
}

// 새로운 리더를 업데이트하는 메서드이다.
// todo 추후 리더 변경 알고리즘이 구상된다면 해당 내용을 반영하여 복수의 수신자가 되어야 한다.
//todo fix path
func (md *MessageDispatcher) SendLeaderUpdateMessage(nodeId p2p.NodeId, leader p2p.Node) error {

	deliverCommand, err := CreateMessageDeliverCommand("UpdateLeader", leader)

	if err != nil {
		return err
	}

	deliverCommand.Recipients = append(deliverCommand.Recipients, nodeId.ToString())

	return md.publisher.Publish("Command", "message.deliver", deliverCommand)
}

// deliver content of node repository to new node
func (md *MessageDispatcher) SendDeliverNodeListMessage(nodeId p2p.NodeId, nodeList []p2p.Node) error {

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
