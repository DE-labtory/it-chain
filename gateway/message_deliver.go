package main

import (
	"encoding/json"

	"github.com/it-chain/bifrost/conn"
	"github.com/it-chain/bifrost/mux"
	"github.com/it-chain/bifrost/pb"
	"github.com/pkg/errors"
)

var FailToSignEnvelopeError = errors.New("Siging Failed")

type MessageDeliver struct {
	signer          Signer
	connectionStore conn.ConnectionStore
}

func NewMessageDeliver(signer Signer, connectionStore conn.ConnectionStore) *MessageDeliver {
	return &MessageDeliver{
		signer:          signer,
		connectionStore: connectionStore,
	}
}

func (m MessageDeliver) Deliver(recipients []string, protocol string, data []byte) error {

	envelope := BuildEnvelope(mux.Protocol(protocol), data)
	signedEnvelope := m.signer.SignEnvelope(envelope)

	if signedEnvelope == nil {
		return FailToSignEnvelopeError
	}

	for _, recipient := range recipients {
		connection := m.connectionStore.GetConnection(conn.ID(recipient))

		if connection != nil {
			connection.Send(signedEnvelope, nil, nil)
		}
	}

	return nil
}

func BuildEnvelope(protocol mux.Protocol, data interface{}) *pb.Envelope {

	payload, _ := json.Marshal(data)
	envelope := &pb.Envelope{}
	envelope.Protocol = string(protocol)
	envelope.Payload = payload

	return envelope
}
