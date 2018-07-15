package adapter

import (
	"time"

	"github.com/it-chain/it-chain-Engine/grpc_gateway"
	"github.com/it-chain/it-chain-Engine/p2p"
)

type PeerService struct {
	publish Publish
}

func (ps *PeerService) Dial(ipAddress string) error {
	command := grpc_gateway.ConnectionCreateCommand{
		Address: ipAddress,
	}
	ps.publish("Command", "connection.create", command)
	return nil
}

//request leader information in p2p network to the node specified by peerId
func (ps *PeerService) RequestLeaderInfo(connectionId string) error {

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

	return ps.publish("Command", "message.deliver", deliverCommand)
}

// command message which requests node list of specific node
func (ps *PeerService) RequestPeerList(peerId p2p.PeerId) error {

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

	return ps.publish("Command", "message.deliver", deliverCommand)
}

