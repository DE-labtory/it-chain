package main

import (
	"encoding/json"

	"github.com/it-chain/bifrost/mux"
	"github.com/it-chain/bifrost/pb"
)

func BuildEnvelope(protocol mux.Protocol, data interface{}) *pb.Envelope {
	payload, _ := json.Marshal(data)
	envelope := &pb.Envelope{}
	envelope.Protocol = string(protocol)
	envelope.Payload = payload

	return envelope
}
