package adapter

import (
	"errors"

	"github.com/it-chain/engine/common"
	"github.com/it-chain/engine/consensus"
	"github.com/it-chain/midgard"
	"github.com/rs/xid"
)

type Publish func(exchange string, topic string, data interface{}) (err error)

type PropagateService struct {
	publish         Publish
	representatives []*consensus.Representative
}

func NewPropagateService(publish Publish, representatives []*consensus.Representative) *PropagateService {
	return &PropagateService{
		publish:         publish,
		representatives: representatives,
	}
}

func (ps PropagateService) BroadcastPrePrepareMsg(msg consensus.PrePrepareMsg) error {
	if msg.ConsensusId.Id == "" {
		return errors.New("Consensus ID is empty")
	}

	if msg.ProposedBlock.Body == nil {
		return errors.New("Block is empty")
	}

	SerializedMsg, err := common.Serialize(msg)

	if err != nil {
		return err
	}

	if err = ps.broadcastMsg(SerializedMsg, "SendPrePrepareMsgProtocol"); err != nil {
		return err
	}

	return nil
}

func (ps PropagateService) BroadcastPrepareMsg(msg consensus.PrepareMsg) error {
	if msg.ConsensusId.Id == "" {
		return errors.New("Consensus ID is empty")
	}

	if msg.BlockHash == nil {
		return errors.New("Block hash is empty")
	}

	SerializedMsg, err := common.Serialize(msg)

	if err != nil {
		return err
	}

	if err = ps.broadcastMsg(SerializedMsg, "SendPrepareMsgProtocol"); err != nil {
		return err
	}

	return nil
}

func (ps PropagateService) BroadcastCommitMsg(msg consensus.CommitMsg) error {
	if msg.ConsensusId.Id == "" {
		return errors.New("Consensus ID is empty")
	}

	SerializedMsg, err := common.Serialize(msg)

	if err != nil {
		return err
	}

	if err = ps.broadcastMsg(SerializedMsg, "SendCommitMsgProtocol"); err != nil {
		return err
	}

	return nil
}

func (ps PropagateService) broadcastMsg(SerializedMsg []byte, protocol string) error {
	if SerializedMsg == nil {
		return errors.New("Message is empty")
	}

	command, err := createSendGrpcMsgCommand(protocol, SerializedMsg)

	if err != nil {
		return err
	}

	for _, r := range ps.representatives {
		command.Recipients = append(command.Recipients, r.GetID())
	}

	return ps.publish("Command", "message.broadcast", command)
}

func createSendGrpcMsgCommand(protocol string, body interface{}) (consensus.SendGrpcMsgCommand, error) {
	data, err := common.Serialize(body)

	if err != nil {
		return consensus.SendGrpcMsgCommand{}, err
	}

	return consensus.SendGrpcMsgCommand{
		CommandModel: midgard.CommandModel{
			ID: xid.New().String(),
		},
		Recipients: make([]string, 0),
		Body:       data,
		Protocol:   protocol,
	}, nil
}
