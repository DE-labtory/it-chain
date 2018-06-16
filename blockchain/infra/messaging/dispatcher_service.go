package messaging

import (
	"errors"

	"github.com/it-chain/it-chain-Engine/blockchain"
	"github.com/it-chain/midgard"
	"github.com/rs/xid"
)

var ErrEmptyBlock = errors.New("block is nil")

type Publisher func(exchange string, topic string, data interface{}) (err error) //해당 publish함수는 midgard 에서 의존성 주입을 받기 위해 interface로 작성한다.
//모든 의존성 주입은 컴포넌트.go 에서 이루어짐

type MessageDispatcher struct {
	publisher Publisher // midgard.client
}

func NewDispatcher(publisher Publisher) *MessageDispatcher {
	return &MessageDispatcher{
		publisher: publisher,
	}
}

func (m *MessageDispatcher) SendBlockCreatedEvent(block blockchain.Block) error {
	if block == nil {
		return ErrEmptyBlock
	}

	event := blockchain.BlockCreatedEvent{
		EventModel: midgard.EventModel{
			ID: xid.New().String(),
		},
		Block: block,
	}

	return m.publisher("Event", "Block", event)
}
