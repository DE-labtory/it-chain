package adapter

import (
	"errors"

	"github.com/it-chain/it-chain-Engine/blockchain"
	"github.com/it-chain/it-chain-Engine/common"
	"github.com/it-chain/midgard"
	"github.com/rs/xid"
)

//kind of error
var ErrEmptyNodeId = errors.New("empty nodeid proposed")

type Publish func(exchange string, topic string, data interface{}) (err error) // 나중에 의존성 주입을 해준다.

type MessageService struct {
	publish Publish // midgard.client.Publish
}

func NewMessageService(publish Publish) *MessageService {
	return &MessageService{
		publish: publish,
	}
}

////request leader information in p2p network to the node specified by nodeId
//func (md *MessageService) RequestBlock(nodeId p2p.NodeId) error {
//	/*
//		입력값세팅:
//		바디세팅:
//		커맨드세팅:
//		퍼블리쉬:
//	*/
//
//	// 입력값 세팅:
//	if nodeId.Id == "" {
//		return ErrEmptyNodeId
//	}
//
//	// 바디 세팅:
//	body := blockchain.BlockRequestMessage{
//		TimeUnix: time.Now().Unix(),
//	}
//
//	// 커맨드 세팅:
//
//	deliverCommand, err := CreateMessagCommand()
//
//	if err != nil {
//		return err
//	}
//
//	deliverCommand.Recipients = append(deliverCommand.Recipients, nodeId.ToString())
//
//	return md.publish()
//}
//
//func (md *MessageService) ResponseBlock(nodeId p2p.NodeId) error {
//
//	return md.publish()
//}

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
