package adapter

import (
	"reflect"
	"testing"

	"github.com/it-chain/it-chain-Engine/p2p"
	"github.com/magiconair/properties/assert"
)

func TestGrpcCommandService_RequestLeaderInfo(t *testing.T) {

	tests := map[string]struct {
		input struct {
			peerId p2p.PeerId
		}
		err error
	}{
		"success": {
			input: struct {
				peerId p2p.PeerId
			}{
				peerId: p2p.PeerId{
					Id: "1",
				},
			},
			err: nil,
		},
		"empty peer id test": {
			input: struct {
				peerId p2p.PeerId
			}{
				peerId: p2p.PeerId{},
			},
			err: ErrEmptyPeerId,
		},
	}

	publish := func(exchange string, topic string, data interface{}) error {
		assert.Equal(t, exchange, "Command")
		assert.Equal(t, topic, "message.deliver")
		assert.Equal(t, reflect.TypeOf(data).String(), "p2p.MessageDeliverCommand")

		return nil
	}

	grpcCommandService := NewGrpcCommandService(publish)

	for testName, test := range tests {
		t.Logf("running test case %s", testName)
		err := grpcCommandService.RequestLeaderInfo(test.input.peerId)
		assert.Equal(t, err, test.err)
	}

}

func TestGrpcCommandService_DeliverLeaderInfo(t *testing.T) {

	tests := map[string]struct {
		input struct {
			peerId p2p.PeerId
			leader p2p.Leader
		}
		err error
	}{
		"empty peer id test": {
			input: struct {
				peerId p2p.PeerId
				leader p2p.Leader
			}{
				peerId: p2p.PeerId{},
				leader: p2p.Leader{
					LeaderId: p2p.LeaderId{
						Id: "1",
					},
				},
			},
			err: ErrEmptyPeerId,
		},
		"empty leader id test": {
			input: struct {
				peerId p2p.PeerId
				leader p2p.Leader
			}{
				peerId: p2p.PeerId{
					Id: "1",
				},
				leader: p2p.Leader{
					LeaderId: p2p.LeaderId{},
				},
			},
			err: ErrEmptyLeaderId,
		},
		"success": {
			input: struct {
				peerId p2p.PeerId
				leader p2p.Leader
			}{
				peerId: p2p.PeerId{
					Id: "1",
				},
				leader: p2p.Leader{
					LeaderId: p2p.LeaderId{
						Id: "1",
					},
				},
			},
			err: nil,
		},
	}
	publish := func(exchange string, topic string, data interface{}) error {
		{
			assert.Equal(t, exchange, "Command")
			assert.Equal(t, topic, "message.deliver")
			assert.Equal(t, reflect.TypeOf(data).String(), "p2p.MessageDeliverCommand")

			return nil
		}
	}

	grpcCommandService := NewGrpcCommandService(publish)

	for testName, test := range tests {
		t.Logf("running test case %s", testName)
		err := grpcCommandService.DeliverLeaderInfo(test.input.peerId, test.input.leader)
		assert.Equal(t, err, test.err)
	}
}

func TestGrpcCommandService_RequestPeerList(t *testing.T) {

	tests := map[string]struct {
		input p2p.PeerId
		err   error
	}{
		"empty peer id test": {
			input: p2p.PeerId{
				Id: "1",
			},
			err: nil,
		},
		"success": {
			input: p2p.PeerId{},
			err:   ErrEmptyPeerId,
		},
	}
	publish := func(exchange string, topic string, data interface{}) error {
		{
			assert.Equal(t, exchange, "Command")
			assert.Equal(t, topic, "message.deliver")
			assert.Equal(t, reflect.TypeOf(data).String(), "p2p.MessageDeliverCommand")

			return nil
		}
	}

	grpcCommandService := NewGrpcCommandService(publish)

	for testName, test := range tests {
		t.Logf("running test case %s", testName)
		err := grpcCommandService.RequestPeerList(test.input)
		assert.Equal(t, err, test.err)
	}
}

func TestGrpcCommandService_DeliverPeerList(t *testing.T) {
	tests := map[string]struct {
		input struct {
			peerId   p2p.PeerId
			peerList []p2p.Peer
		}
		err error
	}{
		"empty peer list test": {
			input: struct {
				peerId   p2p.PeerId
				peerList []p2p.Peer
			}{
				peerId: p2p.PeerId{
					Id: "1",
				},
				peerList: []p2p.Peer{},
			},
			err: ErrEmptyPeerList,
		},
		"empty peer id test": {
			input: struct {
				peerId   p2p.PeerId
				peerList []p2p.Peer
			}{
				peerId:   p2p.PeerId{},
				peerList: []p2p.Peer{},
			},
			err: ErrEmptyPeerId,
		},
		"success": {
			input: struct {
				peerId   p2p.PeerId
				peerList []p2p.Peer
			}{
				peerId: p2p.PeerId{
					Id: "1",
				},
				peerList: []p2p.Peer{
					p2p.Peer{
						PeerId: p2p.PeerId{
							Id: "1",
						},
						IpAddress: "123",
					},
				},
			},
			err: nil,
		},
	}
	publish := func(exchange string, topic string, data interface{}) error {
		{
			assert.Equal(t, exchange, "Command")
			assert.Equal(t, topic, "message.deliver")
			assert.Equal(t, reflect.TypeOf(data).String(), "p2p.MessageDeliverCommand")

			return nil
		}
	}

	grpcCommandService := NewGrpcCommandService(publish)

	for testName, test := range tests {
		t.Logf("running test case %s", testName)
		err := grpcCommandService.DeliverPeerList(test.input.peerId, test.input.peerList)
		assert.Equal(t, err, test.err)
	}

}
