package infra

import (
	"log"
	"sync"

	"errors"

	"github.com/it-chain/bifrost"
	"github.com/it-chain/bifrost/client"
	"github.com/it-chain/bifrost/server"
	"github.com/it-chain/heimdall/key"
	"github.com/it-chain/it-chain-Engine/gateway"
	"github.com/it-chain/midgard"
)

var ErrConnAlreadyExist = errors.New("connection is already exist")

type ConnectionHandler interface {
	OnConnection(connection gateway.Connection)
	OnDisconnection(connection gateway.Connection)
}

type GrpcHostService struct {
	connStore         ConnectionStore
	bifrostServer     *server.Server
	publisher         midgard.Publisher
	priKey            key.PriKey
	pubKey            key.PubKey
	connectionHandler ConnectionHandler
}

func NewGrpcHostService(priKey key.PriKey, pubKey key.PubKey, publisher midgard.Publisher) *GrpcHostService {

	s := server.New(bifrost.KeyOpts{PriKey: priKey, PubKey: pubKey})

	grpcHostService := &GrpcHostService{
		connStore:     NewMemConnectionStore(),
		bifrostServer: s,
		publisher:     publisher,
		priKey:        priKey,
		pubKey:        pubKey,
	}

	s.OnConnection(grpcHostService.onConnection)
	s.OnError(grpcHostService.onError)

	return grpcHostService
}

func (g GrpcHostService) SetHandler(connectionHandler ConnectionHandler) {
	g.connectionHandler = connectionHandler
}

func (g GrpcHostService) Dial(address string) (gateway.Connection, error) {

	connection, err := client.Dial(g.buildDialOption(address))

	if err != nil {
		return gateway.Connection{}, err
	}

	if g.connStore.Exist(connection.GetID()) {
		connection.Close()
		return gateway.Connection{}, ErrConnAlreadyExist
	}

	g.connStore.Add(connection)

	go g.startConnectionUntilClose(connection)

	return toGatewayConnectionModel(connection), nil
}

func toGatewayConnectionModel(connection bifrost.Connection) gateway.Connection {
	return gateway.Connection{
		AggregateModel: midgard.AggregateModel{
			ID: connection.GetID(),
		},
		Address: connection.GetIP(),
	}
}

// connection이 형성되는 경우 실행하는 코드이다.
func (g GrpcHostService) onConnection(connection bifrost.Connection) {

	if g.connStore.Exist(connection.GetID()) {
		connection.Close()
		return
	}

	g.connStore.Add(connection)
	g.connectionHandler.OnConnection(toGatewayConnectionModel(connection))

	g.startConnectionUntilClose(connection)
}

func (g GrpcHostService) startConnectionUntilClose(connection bifrost.Connection) {

	connection.Handle(NewMessageHandler(g.publisher))

	if err := connection.Start(); err != nil {
		connection.Close()
		g.connStore.Delete(connection.GetID())
		g.connectionHandler.OnDisconnection(toGatewayConnectionModel(connection))
	}
}

func (g GrpcHostService) CloseConnection(connID string) {

	connection := g.connStore.Find(connID)

	if connection == nil {
		return
	}

	connection.Close()
	g.connStore.Delete(connection.GetID())
}

func (g GrpcHostService) SendMessages(message []byte, protocol string, connIDs ...string) {

	for _, connID := range connIDs {
		connection := g.connStore.Find(connID)

		if connection != nil {
			connection.Send(message, protocol, nil, nil)
		}
	}
}

func (g GrpcHostService) buildDialOption(address string) (string, client.ClientOpts, client.GrpcOpts) {

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

func (s GrpcHostService) Listen(ip string) {
	s.bifrostServer.Listen(ip)
}

func (s GrpcHostService) onError(err error) {
	log.Fatalln(err.Error())
}

func (s GrpcHostService) Stop() {
	s.bifrostServer.Stop()
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
