package message

import (
	"github.com/golang/protobuf/proto"
)

func (envelope *Envelope) GetMessage() (*StreamMessage, error){

	m := &StreamMessage{}

	err := proto.Unmarshal(envelope.Payload,m)

	if err != nil{
		return nil, err
	}

	return m, nil
}