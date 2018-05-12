package gateway

import (
	"log"

	"github.com/it-chain/bifrost"
	"github.com/it-chain/bifrost/server"
	"github.com/it-chain/heimdall/key"
)

type Server struct {
	handler         bifrost.Handler
	server          *server.Server
	connectionStore *bifrost.ConnectionStore
	publisher       *EventPublisher
}

func NewServer(handler bifrost.Handler, publisher *EventPublisher, connectionStore *bifrost.ConnectionStore, priKey key.PriKey, pubKey key.PubKey) *Server {

	s := server.New(bifrost.KeyOpts{PriKey: priKey, PubKey: pubKey})

	server := &Server{
		handler:         handler,
		server:          s,
		connectionStore: connectionStore,
		publisher:       publisher,
	}

	s.OnConnection(server.onConnection)
	s.OnError(server.onError)

	return server
}

func (s Server) onConnection(connection bifrost.Connection) {

	connection.Handle(s.handler)
	s.connectionStore.AddConnection(connection)

	s.publisher.PublishConnCreatedEvent(connection)

	defer connection.Close()

	if err := connection.Start(); err != nil {
		connection.Close()
	}
}

func (s Server) Listen(ip string) {
	s.server.Listen(ip)
}

func (s Server) onError(err error) {
	log.Fatalln(err.Error())
}

func (s Server) Stop() {
	s.server.Stop()
}
