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
// tested
func (md *MessageDispatcher) RequestLeaderInfo(peer p2p.Node) error {

	requestBody := p2p.LeaderInfoRequestMessage{
		TimeUnix: time.Now().Unix(),
	}

	requestBodyByte, _ := common.Serialize(requestBody)

	deliverCom := &p2p.MessageDeliverCommand{
		CommandModel: midgard.CommandModel{
			ID: xid.New().String(),
		},
		Recipients: make([]string, 0),
		Body:       requestBodyByte,
		Protocol:   "LeaderInfoRequestProtocol",
	}
	deliverCom.Recipients = append(deliverCom.Recipients, peer.NodeId.ToString())

	return md.publisher.Publish("Command", "GrpcMessage", deliverCom)
}

// 단일 피어에게 새로운 리더 정보를 전달하는 메서드이다.
func (md *MessageDispatcher) DeliverLeaderInfo(toPeer p2p.Node, leader p2p.Node) error {

	// 리더 정보를 leaderInfoBody에 담아줌
	leaderInfoBody := p2p.LeaderInfoResponseMessage{
		LeaderId: leader.NodeId.ToString(),
		Address:  leader.IpAddress,
	}

	// 리더 정보 json byte 변환
	leaderInfoBodyByte, _ := common.Serialize(leaderInfoBody)

	// 메세지 전달 이벤트 구조를 담는다.
	deliverCommand := p2p.MessageDeliverCommand{
		CommandModel: midgard.CommandModel{
			ID: xid.New().String(),
		},
		Recipients: make([]string, 0),
		Body:       leaderInfoBodyByte,
		Protocol:   event.LeaderInfoDeliverProtocol,
	}

	// 메세지를 수신할 수신자들을 지정해 준다.
	deliverCommand.Recipients = append(deliverCommand.Recipients, toPeer.NodeId.ToString())

	// message를 보낸다.
	return md.publisher.Publish("Command", "message", deliverCommand)
}

// command message which requests node list of specific node
func (md *MessageDispatcher) RequestNodeList(peer p2p.Node) error {

	requestBody := p2p.LeaderInfoRequestMessage{
		TimeUnix: time.Now().Unix(),
	}

	requestBodyByte, _ := common.Serialize(requestBody)

	commandMessage := &p2p.MessageDeliverCommand{
		CommandModel: midgard.CommandModel{
			ID: xid.New().String(),
		},
		Recipients: make([]string, 0),
		Body:       requestBodyByte,
		Protocol:   "NodeListRequestMessage",
	}

	commandMessage.Recipients = append(commandMessage.Recipients, peer.NodeId.ToString())

	return md.publisher.Publish("Commnand", "GrpcMessage", commandMessage)
}

func (md *MessageDispatcher) DeliverNodeList(toNode p2p.Node, nodeList []p2p.Node) error {
	nodeListByte, _ := common.Serialize(nodeList)

	deliverCommand := p2p.MessageDeliverCommand{
		CommandModel: midgard.CommandModel{
			ID: xid.New().String(),
		},
		Recipients: make([]string, 0),
		Body:       nodeListByte,
		Protocol:   "NodeListDeliverMessage",
	}
	deliverCommand.Recipients = append(deliverCommand.Recipients, toNode.GetID())

	return md.publisher.Publish("Command", "message.*", deliverCommand)
}
