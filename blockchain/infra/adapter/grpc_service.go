package adapter

import (
	"errors"

	"github.com/it-chain/it-chain-Engine/blockchain"
	"github.com/it-chain/it-chain-Engine/common"
	"github.com/it-chain/it-chain-Engine/p2p"
	"github.com/it-chain/midgard"
	"github.com/it-chain/yggdrasill/impl"
	"github.com/rs/xid"
)

//kind of error
var ErrEmptyNodeId = errors.New("empty nodeid proposed")

// ToDo: 구현.(gitId:junk-sound)
type Publish func(exchange string, topic string, data interface{}) (err error)

type MessageService struct {
	publish Publish // midgard.client.Publish
}

func NewMessageService(publish Publish) *MessageService {
	return &MessageService{
		publish: publish,
	}
}

func (md *MessageService) RequestBlock(nodeId p2p.NodeId, height uint64) error {

	if nodeId.Id == "" {
		return ErrEmptyNodeId
	}

	body := blockchain.BlockRequestMessage{
		Height: height,
	}

	deliverCommand, err := createMessageDeliverCommand("BlockRequestProtocol", body)
	if err != nil {
		return err
	}

	deliverCommand.Recipients = append(deliverCommand.Recipients, nodeId.ToString())

	return md.publish("Command", "message.deliver", deliverCommand)
}

func (md *MessageService) ResponseBlock(nodeId p2p.NodeId, block impl.DefaultBlock) error {

	if nodeId.Id == "" {
		return ErrEmptyNodeId
	}

	body := block

	deliverCommand, err := createMessageDeliverCommand("BlockResponseProtocol", body)
	if err != nil {
		return err
	}

	deliverCommand.Recipients = append(deliverCommand.Recipients, nodeId.ToString())

	return md.publish("Command", "message.deliver", deliverCommand)
}

func createMessageDeliverCommand(protocol string, body interface{}) (blockchain.MessageDeliverCommand, error) {

	data, err := common.Serialize(body)
	if err != nil {
		return blockchain.MessageDeliverCommand{}, err
	}

	return blockchain.MessageDeliverCommand{
		CommandModel: midgard.CommandModel{
			ID: xid.New().String(),
		},
		Recipients: make([]string, 0),
		Body:       data,
		Protocol:   protocol,
	}, err
}
