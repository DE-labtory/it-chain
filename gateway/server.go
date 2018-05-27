package gateway

import (
	"log"

	"github.com/it-chain/bifrost"
	"github.com/it-chain/bifrost/server"
	"github.com/it-chain/heimdall/key"
	"github.com/it-chain/midgard"
)

type Server struct {
	server          *server.Server
	connectionStore *bifrost.ConnectionStore
	publisher       midgard.Publisher
}

func NewServer(publisher midgard.Publisher, connectionStore *bifrost.ConnectionStore, priKey key.PriKey, pubKey key.PubKey) *Server {

	s := server.New(bifrost.KeyOpts{PriKey: priKey, PubKey: pubKey})

	server := &Server{
		server:          s,
		connectionStore: connectionStore,
		publisher:       publisher,
	}

	s.OnConnection(server.onConnection)
	s.OnError(server.onError)

	return server
}

func (s Server) onConnection(connection bifrost.Connection) {

	connection.Handle(NewRequestHandler(s.publisher))
	s.connectionStore.AddConnection(connection)

	defer connection.Close()

	err := s.publisher.Publish("Event", "Connection", ConnectionCreatedEvent{
		Address: connection.GetIP(),
		EventModel: midgard.EventModel{
			ID: connection.GetID(),
		},
	})

	if err != nil {
		log.Println(err.Error())
		return
	}

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
