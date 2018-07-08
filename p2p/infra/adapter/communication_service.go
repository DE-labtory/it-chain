package adapter

import (
	"github.com/it-chain/it-chain-Engine/gateway"
	"github.com/it-chain/it-chain-Engine/p2p"
	"time"
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

//request leader information in p2p network to the node specified by peerId
func (cs *CommunicationService) RequestLeaderInfo(connectionId string) error {

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

	return cs.publish("Command", "message.deliver", deliverCommand)
}

// command message which requests node list of specific node
func (cs *CommunicationService) RequestPeerList(peerId p2p.PeerId) error {

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

	return cs.publish("Command", "message.deliver", deliverCommand)
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

func (cs *CommunicationService) DeliverLeaderInfo(connectionId string, leader p2p.Leader) error {

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

	return cs.publish("Command", "message.deliver", deliverCommand)
}