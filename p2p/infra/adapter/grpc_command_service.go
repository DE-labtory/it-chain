package adapter

import (
	"errors"

	"github.com/it-chain/it-chain-Engine/common"
	"github.com/it-chain/it-chain-Engine/p2p"
	"github.com/it-chain/midgard"
	"github.com/rs/xid"
)

//kind of error
var ErrEmptyPeerId = errors.New("empty nodeid proposed")
var ErrEmptyLeaderId = errors.New("empty leader id proposed")
var ErrEmptyConnectionId = errors.New("empty connection ")

type Publish func(exchange string, topic string, data interface{}) (err error) // 나중에 의존성 주입을 해준다.

// message dispatcher sends messages to other nodes in p2p network
type GrpcCommandService struct {

	publish Publish // midgard.client.Publish
}

func NewGrpcCommandService(publish Publish) *GrpcCommandService {

	return &GrpcCommandService{

		publish: publish,
	}
}

func CreateGrpcDeliverCommand(protocol string, body interface{}) (p2p.GrpcDeliverCommand, error) {

	data, err := common.Serialize(body)

	if err != nil {
		return p2p.GrpcDeliverCommand{}, err
	}

	return p2p.GrpcDeliverCommand{

		CommandModel: midgard.CommandModel{
			ID: xid.New().String(),
		},
		Recipients: make([]string, 0),
		Body:       data,
		Protocol:   protocol,
	}, err
}
