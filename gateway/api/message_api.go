package api

import (
	"github.com/it-chain/bifrost"
	"github.com/it-chain/it-chain-Engine/gateway"
	"github.com/it-chain/midgard"
)

type MessageApi struct {
	store     *bifrost.ConnectionStore
	publisher midgard.Publisher
}

func NewMessageApi(store *bifrost.ConnectionStore, publisher midgard.Publisher) *MessageApi {
	return &MessageApi{
		store:     store,
		publisher: publisher,
	}
}

func (m MessageApi) DeliverMessage(command gateway.MessageDeliverCommand) {

	for _, recipient := range command.Recipients {
		connection := m.store.GetConnection(bifrost.ConnID(recipient))

		if connection != nil {
			connection.Send(command.Body, command.Protocol, nil, nil)
		}
	}
}
