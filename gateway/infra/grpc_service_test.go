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
}

func (MockConn) Close() {
	panic("implement me")
}

func (MockConn) GetID() bifrost.ConnID {
	return "1"
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
				Conn: MockConn{},
			},
			output: gateway.MessageReceiveCommand{
				Data:         []byte("hello world"),
				ConnectionID: "1",
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
