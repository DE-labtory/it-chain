package adapter

import (
	"time"

	"errors"

	"github.com/it-chain/it-chain-Engine/common"
	"github.com/it-chain/it-chain-Engine/p2p"
	"github.com/it-chain/midgard"
	"github.com/rs/xid"
)

//kind of error
var ErrEmptyPeerId = errors.New("empty nodeid proposed")
var ErrEmptyLeaderId = errors.New("empty leader id proposed")
var ErrEmptyPeerList = errors.New("empty node list proposed")

type Publish func(exchange string, topic string, data interface{}) (err error) // 나중에 의존성 주입을 해준다.

// message dispatcher sends messages to other nodes in p2p network
type GrpcCommandService struct {
	publish Publish // midgard.client.Publish
}

func NewGrpcCommandService(publish Publish) *GrpcCommandService {
	return &GrpcCommandService{
		publish: publish,
	}
}

//request leader information in p2p network to the node specified by nodeId
func (md *GrpcCommandService) RequestLeaderInfo(nodeId p2p.PeerId) error {

	if nodeId.Id == "" {
		return ErrEmptyPeerId
	}

	body := p2p.LeaderInfoRequestMessage{
		TimeUnix: time.Now().Unix(),
	}

	//message deliver command for delivering leader info
	deliverCommand, err := CreateMessageDeliverCommand("LeaderInfoRequestMessage", body)

	if err != nil {
		return err
	}

	deliverCommand.Recipients = append(deliverCommand.Recipients, nodeId.ToString())

	return md.publish("Command", "message.deliver", deliverCommand)
}

// command message which requests node list of specific node
func (md *GrpcCommandService) RequestPeerList(nodeId p2p.PeerId) error {

	if nodeId.Id == "" {
		return ErrEmptyPeerId
	}
	body := p2p.LeaderInfoRequestMessage{
		TimeUnix: time.Now().Unix(),
	}

	deliverCommand, err := CreateMessageDeliverCommand("PeerListRequestMessage", body)

	if err != nil {
		return err
	}

	deliverCommand.Recipients = append(deliverCommand.Recipients, nodeId.ToString())

	return md.publish("Command", "message.deliver", deliverCommand)
}

func (md *GrpcCommandService) DeliverLeaderInfo(nodeId p2p.PeerId, leader p2p.Leader) error {

	if nodeId.Id == "" {
		return ErrEmptyPeerId
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
func (md *GrpcCommandService) DeliverPeerList(nodeId p2p.PeerId, nodeList []p2p.Peer) error {

	if nodeId.Id == "" {
		return ErrEmptyPeerId
	}

	if len(nodeList) == 0 {
		return ErrEmptyPeerList
	}

	messageDeliverCommand, err := CreateMessageDeliverCommand("PeerListDeliver", nodeList)

	if err != nil {
		return err
	}

	messageDeliverCommand.Recipients = append(messageDeliverCommand.Recipients, nodeId.ToString())

	return md.publish("Command", "message.deliver", messageDeliverCommand)
}

//deliver single node
func (md *GrpcCommandService) DeliverPeer(nodeId p2p.PeerId, node p2p.Peer) error {

	messageDeliverCommand, err := CreateMessageDeliverCommand("PeerDeliverProtocol", node)

	if err != nil {
		return err
	}

	messageDeliverCommand.Recipients = append(messageDeliverCommand.Recipients, node.PeerId.ToString())

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
