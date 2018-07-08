package adapter

import (
	"github.com/it-chain/it-chain-Engine/gateway"
	"github.com/it-chain/it-chain-Engine/p2p"
)

type CommunicationService struct {
	publish Publish
}

func (cs *CommunicationService) Dial(ipAddress string) error {

	command := gateway.ConnectionCreateCommand{
		Address: ipAddress,
	}
	cs.publish("Command", "connection.create", command)
	return nil
}

func (cs *CommunicationService) DeliverPLTable(connectionId string, peerLeaderTable p2p.PLTable) error {

	if connectionId == "" {
		return ErrEmptyPeerId
	}

	if len(peerLeaderTable.PeerList) == 0 {
		return p2p.ErrEmptyPeerList
	}

	//create peer table message
	peerLeaderTableMessage := p2p.PLTableMessage{
		PLTable: peerLeaderTable,
	}

	grpcDeliverCommand, err := CreateGrpcDeliverCommand("PeerTableDeliver", peerLeaderTableMessage)

	if err != nil {
		return err
	}

	grpcDeliverCommand.Recipients = append(grpcDeliverCommand.Recipients, connectionId)

	return cs.publish("Command", "message.deliver", grpcDeliverCommand)
}