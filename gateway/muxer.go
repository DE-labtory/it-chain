package gateway

import (
	"github.com/it-chain/bifrost"
	"github.com/it-chain/bifrost/mux"
)

type Muxer struct {
	*mux.DefaultMux
	amqpPublisher *AMQPPublisher
}

func NewGatewayMux(amqpPublisher *AMQPPublisher) *Muxer {
	muxer := &Muxer{mux.New(), amqpPublisher}

	muxer.Handle("chat", muxer.HandleXXXProtocol)

	return muxer
}

func (m Muxer) HandleXXXProtocol(message bifrost.Message) {

}
