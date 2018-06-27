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
var ErrEmptyConnectionId = errors.New("empty connection ")

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

//request leader information in p2p network to the node specified by peerId
func (md *GrpcCommandService) RequestLeaderInfo(connectionId string) error {

	if connectionId == "" {
		return ErrEmptyPeerId
	}

	body := p2p.LeaderInfoRequestMessage{
		TimeUnix: time.Now().Unix(),
	}

	//message deliver command for delivering leader info
	deliverCommand, err := CreateGrpcDeliverCommand("LeaderInfoRequestMessage", body)

	if err != nil {
		return err
	}

	deliverCommand.Recipients = append(deliverCommand.Recipients, connectionId)

	return md.publish("Command", "message.deliver", deliverCommand)
}

// command message which requests node list of specific node
func (md *GrpcCommandService) RequestPeerList(peerId p2p.PeerId) error {

	if peerId.Id == "" {
		return ErrEmptyPeerId
	}
	body := p2p.LeaderInfoRequestMessage{
		TimeUnix: time.Now().Unix(),
	}

	deliverCommand, err := CreateGrpcDeliverCommand("PeerListRequestMessage", body)

	if err != nil {
		return err
	}

	deliverCommand.Recipients = append(deliverCommand.Recipients, peerId.ToString())

	return md.publish("Command", "message.deliver", deliverCommand)
}

func (md *GrpcCommandService) DeliverLeaderInfo(connectionId string, leader p2p.Leader) error {

	if connectionId == "" {
		return ErrEmptyPeerId
	}

	if leader.LeaderId.Id == "" {
		return ErrEmptyLeaderId
	}

	deliverCommand, err := CreateGrpcDeliverCommand("UpdateLeader", leader)

	if err != nil {
		return err
	}

	deliverCommand.Recipients = append(deliverCommand.Recipients, connectionId)

	return md.publish("Command", "message.deliver", deliverCommand)
}

//deliver node list to other node specified by peerId
func (md *GrpcCommandService) DeliverPeerList(peerId p2p.PeerId, nodeList []p2p.Peer) error {

	if peerId.Id == "" {
		return ErrEmptyPeerId
	}

	if len(nodeList) == 0 {
		return ErrEmptyPeerList
	}

	messageDeliverCommand, err := CreateGrpcDeliverCommand("PeerListDeliver", nodeList)

	if err != nil {
		return err
	}

	messageDeliverCommand.Recipients = append(messageDeliverCommand.Recipients, peerId.ToString())

	return md.publish("Command", "message.deliver", messageDeliverCommand)
}

//deliver single node
func (md *GrpcCommandService) DeliverPeer(connectionId string, node p2p.Peer) error {

	messageDeliverCommand, err := CreateGrpcDeliverCommand("PeerDeliverProtocol", node)

	if err != nil {
		return err
	}

	messageDeliverCommand.Recipients = append(messageDeliverCommand.Recipients, connectionId)

	return md.publish("Command", "message.deliver", messageDeliverCommand)
}
func CreateGrpcDeliverCommand(protocol string, body interface{}) (p2p.GrpcDeliverCommand, error) {

	data, err := common.Serialize(body)

	if err != nil {
		return p2p.GrpcDeliverCommand{}, err
	}

	return p2p.GrpcDeliverCommand{
		CommandModel: midgard.CommandModel{
			ID: xid.New().String(),
		},
		Recipients: make([]string, 0),
		Body:       data,
		Protocol:   protocol,
	}, err
}
