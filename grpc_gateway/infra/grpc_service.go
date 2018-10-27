/*
 * Copyright 2018 It-chain
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package infra

import (
	"errors"
	"sync"

	"github.com/it-chain/bifrost"
	"github.com/it-chain/bifrost/client"
	"github.com/it-chain/bifrost/server"
	"github.com/it-chain/engine/common/command"
	"github.com/it-chain/engine/grpc_gateway"
	"github.com/it-chain/heimdall"
	"github.com/it-chain/iLogger"
)

var ErrConnAlreadyExist = errors.New("connection is already exist")

type Publish func(topic string, data interface{}) (err error)

type ConnectionHandler interface {
	OnConnection(connection grpc_gateway.Connection)
	OnDisconnection(connection grpc_gateway.Connection)
}

type HostID string

type HostInfo struct {
	GrpcGatewayAddress string
	ApiGatewayAddress  string
}

type GrpcHostService struct {
	HostID            HostID
	connStore         ConnectionStore
	metaData          map[string]string
	bifrostServer     *server.Server
	publish           Publish
	priKey            heimdall.PriKey
	pubKey            heimdall.PubKey
	connectionHandler ConnectionHandler
	crypto            bifrost.Crypto
}

func NewGrpcHostService(priKey heimdall.PriKey, pubKey heimdall.PubKey, publish Publish, hostInfo HostInfo, crypto bifrost.Crypto) *GrpcHostService {

	metaData := make(map[string]string, 0)
	metaData["GrpcGatewayAddress"] = hostInfo.GrpcGatewayAddress
	metaData["ApiGatewayAddress"] = hostInfo.ApiGatewayAddress

	s := server.New(bifrost.KeyOpts{PriKey: priKey.(bifrost.Key), PubKey: pubKey.(bifrost.Key)}, crypto, metaData)

	grpcHostService := &GrpcHostService{
		connStore:     NewMemConnectionStore(),
		bifrostServer: s,
		publish:       publish,
		metaData:      metaData,
		priKey:        priKey,
		pubKey:        pubKey,
		crypto:        crypto,
	}

	s.OnConnection(grpcHostService.onConnection)
	s.OnError(grpcHostService.onError)

	return grpcHostService
}

func (g *GrpcHostService) GetHostID() string {
	return g.priKey.ID()
}

func (g *GrpcHostService) SetHandler(connectionHandler ConnectionHandler) {
	g.connectionHandler = connectionHandler
}

func (g *GrpcHostService) Dial(address string) (grpc_gateway.Connection, error) {
	connection, err := client.Dial(g.buildDialOption(address))

	if err != nil {
		return grpc_gateway.Connection{}, err
	}

	if g.connStore.Exist(connection.GetID()) {
		connection.Close()
		g.connStore.Find(connection.GetID())
		return grpc_gateway.Connection{}, ErrConnAlreadyExist
	}

	g.connStore.Add(connection)

	go g.startConnectionUntilClose(connection)

	return toGatewayConnectionModel(connection), nil
}

func toGatewayConnectionModel(connection bifrost.Connection) grpc_gateway.Connection {

	grpcConnection := grpc_gateway.Connection{
		ConnectionID: connection.GetID(),
	}

	metaData := connection.GetMetaData()

	if metaData != nil {
		grpcConnection.GrpcGatewayAddress = metaData["GrpcGatewayAddress"]
		grpcConnection.ApiGatewayAddress = metaData["ApiGatewayAddress"]
	}

	return grpcConnection
}

// connection이 형성되는 경우 실행하는 코드이다.
func (g *GrpcHostService) onConnection(connection bifrost.Connection) {

	if g.connStore.Exist(connection.GetID()) {
		connection.Close()
		return
	}

	g.connStore.Add(connection)
	g.connectionHandler.OnConnection(toGatewayConnectionModel(connection))

	g.startConnectionUntilClose(connection)
}

func (g *GrpcHostService) startConnectionUntilClose(connection bifrost.Connection) {

	iLogger.Infof(nil, "[gRPC-Gateway] Handling connection - ConnectionID: [%s]", connection.GetID())
	connection.Handle(NewMessageHandler(g.publish))

	if err := connection.Start(); err != nil {
		connection.Close()
		g.connStore.Delete(connection.GetID())
		g.connectionHandler.OnDisconnection(toGatewayConnectionModel(connection))
	}
}

func (g *GrpcHostService) GetAllConnections() ([]grpc_gateway.Connection, error) {
	connectionList := g.connStore.FindAll()
	grpcConnectionList := make([]grpc_gateway.Connection, 0)

	for _, connection := range connectionList {
		grpcConnectionList = append(grpcConnectionList, toGatewayConnectionModel(connection))
	}

	return grpcConnectionList, nil
}

func (g *GrpcHostService) CloseAllConnections() error {
	for _, connection := range g.connStore.FindAll() {
		connection.Close()
	}

	return nil
}

func (g *GrpcHostService) CloseConnection(connID string) error {

	connection := g.connStore.Find(connID)

	err := errors.New("invalid nodeID")

	if connection == nil {
		return err
	}

	connection.Close()
	g.connStore.Delete(connection.GetID())

	return nil
}

func (g *GrpcHostService) SendMessages(message []byte, protocol string, connIDs ...string) error {

	//wg := sync.WaitGroup{}
	//wg.Add(len(connIDs))

	var err error
	for _, connID := range connIDs {
		connection := g.connStore.Find(connID)

		if connection != nil {
			connection.Send(message, protocol, nil, nil)
		}
	}

	//wg.Wait()
	return err
}

func (g *GrpcHostService) IsConnectionExist(connectionID string) bool {
	if g.connStore.Find(connectionID) == nil {
		return false
	}

	return true
}

func (g GrpcHostService) buildDialOption(address string) (string, map[string]string, client.ClientOpts, client.GrpcOpts, bifrost.Crypto) {

	clientOpt := client.ClientOpts{
		Ip:     address,
		PubKey: g.pubKey.(bifrost.Key),
	}

	grpcOpt := client.GrpcOpts{
		TlsEnabled: false,
		Creds:      nil,
	}

	return address, g.metaData, clientOpt, grpcOpt, g.crypto
}

func (s *GrpcHostService) Listen(ip string) {
	s.bifrostServer.Listen(ip)
}

func (s *GrpcHostService) onError(err error) {
	iLogger.Fatalf(nil, "[gRPC-Gateway] Connection error - [Err]: [%s]", err.Error())
}

func (s *GrpcHostService) Stop() {
	s.bifrostServer.Stop()
}

type ConnectionStore interface {
	Exist(connID bifrost.ConnID) bool
	Add(conn bifrost.Connection) error
	Delete(connID bifrost.ConnID)
	Find(connID bifrost.ConnID) bifrost.Connection
	FindAll() []bifrost.Connection
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
	connStore.Lock()
	defer connStore.Unlock()

	_, ok := connStore.connMap[connID]

	//exist
	if ok {
		return true
	}

	return false
}

func (connStore MemConnectionStore) exist(connID bifrost.ConnID) bool {

	_, ok := connStore.connMap[connID]

	//exist
	if ok {
		return true
	}

	return false
}

func (connStore MemConnectionStore) Add(conn bifrost.Connection) error {

	connStore.Lock()
	defer connStore.Unlock()

	connID := conn.GetID()

	if connStore.exist(connID) {
		return ErrConnAlreadyExist
	}

	connStore.connMap[connID] = conn

	return nil
}

func (connStore MemConnectionStore) Delete(connID bifrost.ConnID) {

	connStore.Lock()
	defer connStore.Unlock()

	if connStore.exist(connID) {
		delete(connStore.connMap, connID)
	}
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

func (connStore MemConnectionStore) FindAll() []bifrost.Connection {
	connStore.Lock()
	defer connStore.Unlock()

	connectionList := make([]bifrost.Connection, 0)

	for _, connection := range connStore.connMap {
		connectionList = append(connectionList, connection)
	}

	return connectionList
}

type MessageHandler struct {
	publish Publish
}

func (r MessageHandler) ServeRequest(msg bifrost.Message) {

	if msg.Envelope == nil {
		return
	}

	err := r.publish("message.receive", command.ReceiveGrpc{
		Body:         msg.Data,
		ConnectionID: msg.Conn.GetID(),
		Protocol:     msg.Envelope.Protocol,
	})

	if err != nil {
		iLogger.Errorf(nil, "[gRPC-Gateway] Fail to publish message received error - [Err]: [%s]", err.Error())
	}
}

func NewMessageHandler(publish Publish) MessageHandler {
	return MessageHandler{
		publish: publish,
	}
}
