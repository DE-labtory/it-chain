package adapter

import (
	"github.com/it-chain/it-chain-Engine/p2p"
)

type CommunicationService struct {
	publish Publish
}

func NewCommunicationService(publish Publish) *CommunicationService{

	return &CommunicationService{
		publish:publish,
	}
}

func (cs *CommunicationService) Dial(ipAddress string) error {

	command := p2p.ConnectionCreateCommand{
		Address: ipAddress,
	}

	cs.publish("Command", "connection.create", command)

	return nil
}

func (cs *CommunicationService) DeliverPLTable(connectionId string, peerLeaderTable p2p.PLTable) error {

	if connectionId == "" {
		return ErrEmptyConnectionId
	}

	if len(peerLeaderTable.PeerTable) == 0 {
		return p2p.ErrEmptyPeerTable
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