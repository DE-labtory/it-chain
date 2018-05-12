package gateway

import (
	"log"

	"github.com/it-chain/bifrost"
	"github.com/it-chain/bifrost/server"
	"github.com/it-chain/heimdall/key"
)

type Server struct {
	mux             *Muxer
	server          *server.Server
	ConnectionStore *bifrost.ConnectionStore
	publisher       *EventPublisher
}

func NewServer(mux *Muxer, priKey key.PriKey, pubKey key.PubKey) *Server {

	s := server.New(bifrost.KeyOpts{PriKey: priKey, PubKey: pubKey})

	server := &Server{
		mux:    mux,
		server: s,
	}

	s.OnConnection(server.onConnection)
	s.OnError(server.onError)

	return server
}

func (s Server) onConnection(connection bifrost.Connection) {

	connection.Handle(s.mux)
	s.ConnectionStore.AddConnection(connection)

	s.publisher.ConnCreatedEvent(connection)

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
