package gateway

import (
	"github.com/it-chain/bifrost"
	"github.com/it-chain/bifrost/mux"
)

type Muxer struct {
	*mux.DefaultMux
	publisher *EventPublisher
}

func NewGatewayMux(publisher *EventPublisher) *Muxer {
	muxer := &Muxer{mux.New(), publisher}

	muxer.Handle("chat", muxer.HandleXXXProtocol)

	return muxer
}

func (m Muxer) HandleXXXProtocol(message bifrost.Message) {

}
