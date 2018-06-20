package adapter

import (
	"errors"

	"github.com/it-chain/it-chain-Engine/p2p"
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

//request leader information in p2p network to the node specified by nodeId
func (md *MessageService) RequestBlock(nodeId p2p.NodeId) error {
	/*
		입력값세팅:
		바디세팅:
		커맨드세팅:
		퍼블리쉬:
	*/

	return md.publish()
}

func (md *MessageService) ResponseBlock(nodeId p2p.NodeId) error {

	return md.publish()
}
