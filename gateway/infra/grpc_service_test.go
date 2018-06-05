package infra_test

import (
	"testing"

	"github.com/it-chain/bifrost"
	"github.com/it-chain/heimdall/key"
	"github.com/it-chain/it-chain-Engine/gateway"
	"github.com/it-chain/it-chain-Engine/gateway/infra"
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

func (MockConn) GetIP() string {
	return "1"
}

func (MockConn) GetPeerKey() key.PubKey {
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

//MessageHandler
func TestMessageHandler_ServeRequest(t *testing.T) {

	//given
	tests := map[string]struct {
		input  bifrost.Message
		output gateway.MessageReceiveCommand
		err    error
	}{
		"success": {
			input: bifrost.Message{
				Data: []byte("hello world"),
				Conn: MockConn{
					ID: "123",
				},
			},
			output: gateway.MessageReceiveCommand{
				Data:         []byte("hello world"),
				ConnectionID: "123",
			},
			err: nil,
		},
	}

	var publish = func(exchange string, topic string, data interface{}) (err error) {

		//then
		assert.Equal(t, exchange, "Command")
		assert.Equal(t, topic, "message.receive")
		assert.Equal(t, data, gateway.MessageReceiveCommand{
			Data:         []byte("hello world"),
			ConnectionID: "1",
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
		"add same id": {
			input:  MockConn{ID: "123"},
			output: MockConn{},
			err:    infra.ErrConnAlreadyExist,
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
