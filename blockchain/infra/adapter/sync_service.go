package adapter

import (
	"github.com/it-chain/it-chain-Engine/blockchain"
	"github.com/it-chain/it-chain-Engine/common"
	"github.com/it-chain/midgard"
	"github.com/rs/xid"
)

// ToDo: 구현.(gitId:junk-sound)
type Publish func(exchange string, topic string, data interface{}) (err error)

type SyncService struct {
	publish Publish // midgard.client.Publish
}

func NewGrpcCommandService(publish Publish) *SyncService {
	return &SyncService{
		publish: publish,
	}
}

func (ss *SyncService) RequestBlock(peerId blockchain.PeerId, height uint64) error {
	if peerId.Id == "" {
		return ErrEmptyNodeId
	}

	body := height

	deliverCommand, err := createGrpcDeliverCommand("BlockRequestProtocol", body)
	if err != nil {
		return err
	}

	deliverCommand.Recipients = append(deliverCommand.Recipients, peerId.ToString())

	return ss.publish("Command", "message.deliver", deliverCommand)
}

func (ss *SyncService) ResponseBlock(peerId blockchain.PeerId, block blockchain.Block) error {
	if peerId.Id == "" {
		return ErrEmptyNodeId
	}

	if block.GetSeal() == nil {
		return ErrEmptyBlockSeal
	}

	body := block

	deliverCommand, err := createGrpcDeliverCommand("BlockResponseProtocol", body)
	if err != nil {
		return err
	}

	deliverCommand.Recipients = append(deliverCommand.Recipients, peerId.ToString())

	return ss.publish("Command", "message.deliver", deliverCommand)
}

// TODO "SyncCheckResponseProtocol"을 통해서 last block을 전달한다.
func (ss *SyncService) SyncCheckResponse(block blockchain.Block) error {
	return nil
}

func createGrpcDeliverCommand(protocol string, body interface{}) (blockchain.GrpcDeliverCommand, error) {

	data, err := common.Serialize(body)
	if err != nil {
		return blockchain.GrpcDeliverCommand{}, err
	}

	return blockchain.GrpcDeliverCommand{
		CommandModel: midgard.CommandModel{
			ID: xid.New().String(),
		},
		Recipients: make([]string, 0),
		Body:       data,
		Protocol:   protocol,
	}, err
}
