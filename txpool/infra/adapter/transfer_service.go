package adapter

import (
	"errors"

	"github.com/it-chain/engine/common"
	"github.com/it-chain/engine/txpool"
	"github.com/it-chain/midgard"
	"github.com/rs/xid"
)

var ErrTxEmpty = errors.New("Empty transaction list proposed")

type Publisher func(exchange string, topic string, data interface{}) (err error) //해당 publish함수는 midgard 에서 의존성 주입을 받기 위해 interface로 작성한다.
//모든 의존성 주입은 컴포넌트.go 에서 이루어짐

type TransferService struct {
	publisher Publisher // midgard.client
}

func NewTransferService(publisher Publisher) *TransferService {
	return &TransferService{
		publisher: publisher,
	}
}

func (ts TransferService) SendLeaderTransactions(transactions []txpool.Transaction, leader txpool.Leader) error {

	if len(transactions) == 0 {
		return ErrTxEmpty
	}

	deliverCommand, err := createGrpcDeliverCommand("SendLeaderTransactionsProtocol", transactions)

	if err != nil {
		return err
	}

	deliverCommand.Recipients = append(deliverCommand.Recipients, leader.LeaderId.ToString())

	return ts.publisher("Command", "message.deliver", deliverCommand)
}

func createGrpcDeliverCommand(protocol string, body interface{}) (txpool.GrpcDeliverCommand, error) {

	data, err := common.Serialize(body)

	if err != nil {
		return txpool.GrpcDeliverCommand{}, err
	}

	return txpool.GrpcDeliverCommand{
		CommandModel: midgard.CommandModel{
			ID: xid.New().String(),
		},
		Recipients: make([]string, 0),
		Body:       data,
		Protocol:   protocol,
	}, err
}
