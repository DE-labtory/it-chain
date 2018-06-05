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

	// 최초로 서버를 생성하는 경우 connection이 형성되었음을 알리고 해당 함수를 실행시킨다.
	// bifrost onconnection 함수에 server의 onconnection 함수를 전달.
	// bifrost에서 onconnection 함수 실행하도록 위임
	s.OnConnection(server.onConnection)
	s.OnError(server.onError)

	return server
}


// gRPC connection이 생성되면 실행되는 함수를 정의한다.
// 서버 생성 시점에 bifrost에 함수가 전달되고 실행이 위임된다.
func (s Server) onConnection(connection bifrost.Connection) {

	connection.Handle(NewRequestHandler(s.publisher))
	s.connectionStore.AddConnection(connection)

	defer connection.Close()

	// connection이 생성되면 ConnectionCreatedEvent 전파!
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
