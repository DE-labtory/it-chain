package infra

import (
	"log"
	"sync"

	"errors"

	"github.com/it-chain/bifrost"
	"github.com/it-chain/bifrost/client"
	"github.com/it-chain/heimdall/key"
	"github.com/it-chain/it-chain-Engine/gateway"
	"github.com/it-chain/midgard"
)

var ErrConnAlreadyExist = errors.New("connection is already exist")

type GrpcService struct {
	connStore ConnectionStore
	publisher midgard.Publisher
	priKey    key.PriKey
	pubKey    key.PubKey
}

func NewGrpcDialService(priKey key.PriKey, pubKey key.PubKey) *GrpcService {
	return &GrpcService{}
}

func (g GrpcService) Dial(address string) (gateway.Connection, error) {

	connection, err := client.Dial(g.buildDialOption(address))

	if err != nil {
		return gateway.Connection{}, err
	}

	if g.connStore.Exist(connection.GetID()) {
		connection.Close()
		return gateway.Connection{}, ErrConnAlreadyExist
	}

	connection.Handle(NewMessageHandler(g.publisher))

	go func() {
		defer connection.Close()

		if err := connection.Start(); err != nil {
			connection.Close()
			log.Printf("connections are closing [%s]", err)
		}

		g.connStore.Delete(connection.GetID())
		return
	}()

	g.connStore.Add(connection)

	return gateway.Connection{
		ID:      connection.GetID(),
		Address: connection.GetIP(),
	}, nil
}

func (g GrpcService) CloseConnection(connID string) {
	connection := g.connStore.Find(connID)

	if connection == nil {
		return
	}

	connection.Close()
	g.connStore.Delete(connection.GetID())
}

func (g GrpcService) SendMessages(message []byte, protocol string, connIDs ...string) {

	for _, connID := range connIDs {
		connection := g.connStore.Find(connID)

		if connection != nil {
			connection.Send(message, protocol, nil, nil)
		}
	}
}

func (g GrpcService) buildDialOption(address string) (string, client.ClientOpts, client.GrpcOpts) {

	clientOpt := client.ClientOpts{
		Ip:     address,
		PriKey: g.priKey,
		PubKey: g.pubKey,
	}

	grpcOpt := client.GrpcOpts{
		TlsEnabled: false,
		Creds:      nil,
	}

	return address, clientOpt, grpcOpt
}

type ConnectionStore interface {
	Exist(connID bifrost.ConnID) bool
	Add(conn bifrost.Connection)
	Delete(connID bifrost.ConnID)
	Find(connID bifrost.ConnID) bifrost.Connection
}

type MemConnectionStore struct {
	sync.RWMutex
	connMap map[bifrost.ConnID]bifrost.Connection
}

func NewMemConnectionStore() MemConnectionStore {
	return MemConnectionStore{
		connMap: make(map[bifrost.ConnID]bifrost.Connection),
	}
}

func (connStore MemConnectionStore) Exist(connID bifrost.ConnID) bool {

	_, ok := connStore.connMap[connID]

	//exist
	if ok {
		return true
	}

	return false
}

func (connStore MemConnectionStore) Add(conn bifrost.Connection) {

	connStore.Lock()
	defer connStore.Unlock()

	connID := conn.GetID()

	if connStore.Exist(connID) {
		return
	}
	connStore.connMap[connID] = conn
}

func (connStore MemConnectionStore) Delete(connID bifrost.ConnID) {

	connStore.Lock()
	defer connStore.Unlock()

	delete(connStore.connMap, connID)
}

func (connStore MemConnectionStore) Find(connID bifrost.ConnID) bifrost.Connection {

	connStore.Lock()
	conn, ok := connStore.connMap[connID]

	connStore.Unlock()

	//exist
	if ok {
		return conn
	}

	return nil
}

type MessageHandler struct {
	publisher midgard.Publisher
}

func (r MessageHandler) ServeRequest(msg bifrost.Message) {

	err := r.publisher.Publish("Command", "Message", gateway.MessageReceiveCommand{
		Data:         msg.Data,
		ConnectionID: msg.Conn.GetID(),
	})

	if err != nil {
		log.Println(err.Error())
	}
}

func NewMessageHandler(publisher midgard.Publisher) MessageHandler {
	return MessageHandler{
		publisher: publisher,
	}
}
