package adapter

import "github.com/it-chain/it-chain-Engine/p2p"

type LeaderService struct {
	publish Publish
}

func (ls *LeaderService) DeliverLeaderInfo(connectionId string, leader p2p.Leader) error {

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

	return ls.publish("Command", "message.deliver", deliverCommand)
}