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

package infra_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/it-chain/bifrost"
	"github.com/it-chain/bifrost/pb"
	"github.com/it-chain/engine/common/command"
	"github.com/it-chain/engine/grpc_gateway"
	"github.com/it-chain/engine/grpc_gateway/infra"
	"github.com/it-chain/heimdall"
	"github.com/it-chain/heimdall/hashing"
	"github.com/it-chain/heimdall/hecdsa"
	"github.com/it-chain/heimdall/keystore"
	"github.com/stretchr/testify/assert"
)

type MockConn struct {
	ID string
}

func (MockConn) Close() {
	panic("implement me")
}

func (m MockConn) GetID() bifrost.ConnID {
	return m.ID
}

func (MockConn) GetIP() bifrost.Address {
	return bifrost.Address{IP: "1"}
}

func (MockConn) GetPeerKey() bifrost.Key {
	panic("implement me")
}

func (MockConn) Handle(handler bifrost.Handler) {
	panic("implement me")
}

func (MockConn) Send(data []byte, protocol string, successCallBack func(interface{}), errCallBack func(error)) {
	panic("implement me")
}

func (MockConn) Start() error {
	panic("implement me")
}

func (MockConn) GetMetaData() map[string]string {
	return nil
}

//MessageHandler
func TestMessageHandler_ServeRequest(t *testing.T) {

	//given
	tests := map[string]struct {
		input  bifrost.Message
		output command.ReceiveGrpc
		err    error
	}{
		"success": {
			input: bifrost.Message{
				Data: []byte("hello world"),
				Conn: MockConn{
					ID: "123",
				},
				Envelope: &pb.Envelope{
					Protocol: "protocol",
				},
			},
			output: command.ReceiveGrpc{
				Body:         []byte("hello world"),
				ConnectionID: "123",
				Protocol:     "protocol",
			},
			err: nil,
		},
	}

	var publish = func(topic string, data interface{}) (err error) {

		//then
		assert.Equal(t, topic, "message.receive")
		assert.Equal(t, data, command.ReceiveGrpc{
			Body:         []byte("hello world"),
			ConnectionID: "123",
			Protocol:     "protocol",
		})
		return nil
	}

	messageHandler := infra.NewMessageHandler(publish)

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)

		//when
		messageHandler.ServeRequest(test.input)
	}
}

func TestMemConnectionStore_Add(t *testing.T) {

	//given
	tests := map[string]struct {
		input  bifrost.Connection
		output bifrost.Connection
		err    error
	}{
		"add success": {
			input:  MockConn{ID: "123"},
			output: MockConn{ID: "123"},
			err:    nil,
		},
	}

	connectionStore := infra.NewMemConnectionStore()

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)

		//when
		err := connectionStore.Add(test.input)

		//then
		assert.Equal(t, connectionStore.Find(test.input.GetID()), test.output)
		assert.Equal(t, err, test.err)
	}
}

func TestMemConnectionStore_Find(t *testing.T) {

	//given
	tests := map[string]struct {
		input  string
		output bifrost.Connection
		err    error
	}{
		"find success": {
			input:  "123",
			output: MockConn{ID: "123"},
			err:    nil,
		},
		"find nil": {
			input:  "124",
			output: nil,
			err:    nil,
		},
	}

	connectionStore := infra.NewMemConnectionStore()
	connectionStore.Add(MockConn{ID: "123"})

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)

		//when
		conn := connectionStore.Find(test.input)

		//then
		assert.Equal(t, conn, test.output)
	}
}

func TestMemConnectionStore_Delete(t *testing.T) {

	//given
	tests := map[string]struct {
		input  string
		output bifrost.Connection
		err    error
	}{
		"delete success": {
			input:  "123",
			output: nil,
			err:    nil,
		},
		"delete fail": {
			input:  "124",
			output: nil,
			err:    nil,
		},
	}

	connectionStore := infra.NewMemConnectionStore()
	connectionStore.Add(MockConn{ID: "123"})

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)

		//when
		connectionStore.Delete(test.input)

		//then
		assert.Equal(t, connectionStore.Find(test.input), test.output)
	}
}

type MockHandler struct {
	OnConnectionFunc    func(connection grpc_gateway.Connection)
	OnDisconnectionFunc func(connection grpc_gateway.Connection)
}

func (m *MockHandler) OnConnection(connection grpc_gateway.Connection) {
	m.OnConnectionFunc(connection)
}

func (m *MockHandler) OnDisconnection(connection grpc_gateway.Connection) {
	m.OnDisconnectionFunc(connection)
}

type mockECDSASigner struct {
	keyDirPath string
	hashOpt    *hashing.HashOpt
}

func (signer *mockECDSASigner) Sign(message []byte) ([]byte, error) {
	return hecdsa.SignWithKeyInLocal(signer.keyDirPath, message, signer.hashOpt)
}

type mockECDSAVerifier struct {
	signerOpt heimdall.SignerOpts
}

func (verifier *mockECDSAVerifier) Verify(peerKey bifrost.Key, signature, message []byte) (bool, error) {
	return hecdsa.Verify(peerKey.(heimdall.PubKey), signature, message, verifier.signerOpt)
}

type mockECDSAKeyRecoverer struct {
}

func (rec *mockECDSAKeyRecoverer) RecoverKeyFromByte(keyBytes []byte, isPrivate bool) (bifrost.Key, error) {
	recoverer := &hecdsa.KeyRecoverer{}
	key, err := recoverer.RecoverKeyFromByte(keyBytes, isPrivate)
	return key.(bifrost.Key), err
}

var setupGrpcHostService = func(t *testing.T, ip string, keyDirPath string, publish func(topic string, data interface{}) error) (*infra.GrpcHostService, func()) {
	keyGenOpt, err := hecdsa.NewKeyGenOpt("P-384")
	assert.NoError(t, err)
	pri, err := hecdsa.GenerateKey(keyGenOpt)
	assert.NoError(t, err)
	pub := pri.PublicKey()

	err = keystore.StorePriKeyWithoutPwd(pri, keyDirPath)
	assert.NoError(t, err)

	hashOpt, err := hashing.NewHashOpt("SHA384")
	assert.NoError(t, err)

	signer := &mockECDSASigner{
		keyDirPath: keyDirPath,
		hashOpt:    hashOpt,
	}

	signerOpt := hecdsa.NewSignerOpts(hashOpt)
	verifier := &mockECDSAVerifier{
		signerOpt: signerOpt,
	}

	keyRecoverer := &mockECDSAKeyRecoverer{}

	crypto := bifrost.Crypto{
		Signer:       signer,
		Verifier:     verifier,
		KeyRecoverer: keyRecoverer,
	}

	hostService := infra.NewGrpcHostService(pri, pub, publish, infra.HostInfo{
		GrpcGatewayAddress: ip,
	}, crypto)

	go hostService.Listen(ip)

	return hostService, func() {
		hostService.Stop()
		time.Sleep(3 * time.Second)
		os.RemoveAll(keyDirPath)
	}
}

func TestGrpcHostService_Dial(t *testing.T) {

	//given
	var publish = func(topic string, data interface{}) (err error) {
		return nil
	}

	serverHostService, tearDown1 := setupGrpcHostService(t, "127.0.0.1:7777", "server", publish)
	clientHostService, tearDown2 := setupGrpcHostService(t, "127.0.0.1:8888", "client", publish)

	//times to need to setup server
	time.Sleep(3 * time.Second)

	tests := map[string]struct {
		input  string
		output string
		err    error
	}{
		"dial success": {
			input:  "127.0.0.1:7777",
			output: serverHostService.GetHostID(),
			err:    nil,
		},
	}

	handler := &MockHandler{}
	handler.OnConnectionFunc = func(connection grpc_gateway.Connection) {
		fmt.Println(connection)
	}

	handler.OnDisconnectionFunc = func(connection grpc_gateway.Connection) {
		fmt.Println("connection is closing", connection)
	}

	defer tearDown2()
	defer tearDown1()

	serverHostService.SetHandler(handler)
	clientHostService.SetHandler(handler)

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)

		conn, err := clientHostService.Dial(test.input)
		assert.Equal(t, err, test.err)
		assert.Equal(t, conn.ConnectionID, test.output)
	}
}

func TestGrpcHostService_Dial_When_connection_exist(t *testing.T) {

	//given
	tests := map[string]struct {
		input  string
		output string
		err    error
	}{
		"dial exist connection": {
			input:  "127.0.0.1:7777",
			output: "",
			err:    infra.ErrConnAlreadyExist,
		},
	}

	var publish = func(topic string, data interface{}) (err error) {

		return nil
	}

	serverHostService, tearDown1 := setupGrpcHostService(t, "127.0.0.1:7777", "server", publish)
	clientHostService, tearDown2 := setupGrpcHostService(t, "127.0.0.1:8888", "client", publish)

	//times to need to setup server
	time.Sleep(3 * time.Second)

	handler := &MockHandler{}
	handler.OnConnectionFunc = func(connection grpc_gateway.Connection) {
		fmt.Println(connection)
	}

	handler.OnDisconnectionFunc = func(connection grpc_gateway.Connection) {
		fmt.Println("connection is closing", connection)
	}

	defer tearDown1()
	defer tearDown2()

	serverHostService.SetHandler(handler)
	clientHostService.SetHandler(handler)

	_, err := clientHostService.Dial("127.0.0.1:7777")
	assert.NoError(t, err)

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)

		conn, err := clientHostService.Dial(test.input)
		assert.Equal(t, err, test.err)
		assert.Equal(t, conn.GrpcGatewayAddress, test.output)

	}
}

func TestGrpcHostService_SendMessages(t *testing.T) {

	//given
	tests := map[string]struct {
		input struct {
			Ip       string
			Message  []byte
			Protocol string
		}
		output string
		err    error
	}{
		"send message success": {
			input: struct {
				Ip       string
				Message  []byte
				Protocol string
			}{
				Ip:       "127.0.0.1:7777",
				Message:  []byte("hello"),
				Protocol: "testProtocol",
			},
			output: "127.0.0.1:7777",
			err:    nil,
		},
	}

	var publishedData []byte
	var connID string

	var publish = func(topic string, data interface{}) (err error) {

		assert.Equal(t, topic, "message.receive")
		assert.Equal(t, data, command.ReceiveGrpc{
			Body:         publishedData,
			ConnectionID: connID,
			Protocol:     "testProtocol",
		})
		return nil
	}

	serverHostService, tearDown1 := setupGrpcHostService(t, "127.0.0.1:7777", "server", publish)
	clientHostService, tearDown2 := setupGrpcHostService(t, "127.0.0.1:8888", "client", publish)

	//times to need to setup server
	time.Sleep(3 * time.Second)

	defer tearDown1()
	defer tearDown2()

	handler := &MockHandler{}
	handler.OnConnectionFunc = func(connection grpc_gateway.Connection) {
		connID = connection.ConnectionID
	}

	handler.OnDisconnectionFunc = func(connection grpc_gateway.Connection) {
		fmt.Println("connection is closing", connection)
	}

	serverHostService.SetHandler(handler)
	clientHostService.SetHandler(handler)

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)

		conn, err := clientHostService.Dial(test.input.Ip)
		assert.NoError(t, err)

		publishedData = test.input.Message

		err = clientHostService.SendMessages(test.input.Message, test.input.Protocol, conn.ConnectionID)
		assert.NoError(t, err)
	}
}

func TestGrpcHostService_Close(t *testing.T) {

	//given
	tests := map[string]struct {
		input  string
		output string
		err    error
	}{
		"dial success": {
			input:  "127.0.0.1:7777",
			output: "127.0.0.1:7777",
			err:    nil,
		},
		"dial exist connection": {
			input:  "127.0.0.1:7777",
			output: "",
			err:    infra.ErrConnAlreadyExist,
		},
	}

	var publish = func(topic string, data interface{}) (err error) {

		return nil
	}

	serverHostService, tearDown1 := setupGrpcHostService(t, "127.0.0.1:7777", "server", publish)
	clientHostService, tearDown2 := setupGrpcHostService(t, "127.0.0.1:8888", "client", publish)

	//times to need to setup server
	time.Sleep(3 * time.Second)

	defer tearDown1()
	defer tearDown2()

	var connID string

	handler := &MockHandler{}
	handler.OnConnectionFunc = func(connection grpc_gateway.Connection) {
		fmt.Println(connection)
		connID = connection.ConnectionID
	}

	handler.OnDisconnectionFunc = func(connection grpc_gateway.Connection) {
		fmt.Println("connection is closing", connection)
		assert.Equal(t, connID, connection.ConnectionID)
	}

	serverHostService.SetHandler(handler)
	clientHostService.SetHandler(handler)

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)

		conn, err := clientHostService.Dial(test.input)
		assert.NoError(t, err)
		clientHostService.CloseConnection(conn.ConnectionID)
	}
}

func TestMemConnectionStore_FindAll(t *testing.T) {

	connStore := infra.NewMemConnectionStore()
	connectionList := []bifrost.Connection{&bifrost.GrpcConnection{ID: "123"}, &bifrost.GrpcConnection{ID: "124"}}
	for _, connection := range connectionList {
		connStore.Add(connection)
	}

	foundedConnectionList := connStore.FindAll()
	assert.Equal(t, 2, len(foundedConnectionList))
	assert.Contains(t, foundedConnectionList, &bifrost.GrpcConnection{ID: "123"})
	assert.Contains(t, foundedConnectionList, &bifrost.GrpcConnection{ID: "124"})
}

func TestGrpcHostService_GetAllConnections(t *testing.T) {

	var publish = func(topic string, data interface{}) (err error) {
		return nil
	}

	serverHostService, tearDown := setupGrpcHostService(t, "127.0.0.1:7777", "server", publish)
	clientHostService, tearDown2 := setupGrpcHostService(t, "127.0.0.1:8888", "client", publish)

	//times to need to setup server
	time.Sleep(3 * time.Second)
	defer tearDown()
	defer tearDown2()

	var connID string
	handler := &MockHandler{}
	handler.OnConnectionFunc = func(connection grpc_gateway.Connection) {
		connID = connection.ConnectionID
	}

	handler.OnDisconnectionFunc = func(connection grpc_gateway.Connection) {

	}

	serverHostService.SetHandler(handler)
	clientHostService.SetHandler(handler)

	_, err := clientHostService.Dial("127.0.0.1:7777")
	assert.NoError(t, err)

	connections, err := serverHostService.GetAllConnections()
	assert.NoError(t, err)
	assert.Equal(t, connections[0].ConnectionID, connID)
}
