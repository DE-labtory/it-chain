package infra

import (
	"log"

	"github.com/it-chain/bifrost"
	"github.com/it-chain/midgard"
)

type RequestHandler struct {
	publisher midgard.Publisher
}

func (r RequestHandler) ServeRequest(msg bifrost.Message) {

	err := r.publisher.Publish("Command", "Message", GrpcRequestCommand{
		Data:         msg.Data,
		ConnectionID: msg.Conn.GetID(),
	})

	if err != nil {
		log.Println(err.Error())
	}
}

func NewRequestHandler(publisher midgard.Publisher) RequestHandler {
	return RequestHandler{
		publisher: publisher,
	}
}
