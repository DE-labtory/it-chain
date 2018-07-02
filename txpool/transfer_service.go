package txpool

import (
	"errors"
	"log"
)

//todo
type TxPeriodicTransferService struct {
	txRepository     TransactionRepository
	leaderRepository LeaderRepository
	grpcService   GrpcCommandService
}

//todo (진행중) 이 함수가 call되었을 때 조건에 맞는 tx를 leader에게 전송하는 로직 추가
//todo (완료) infra의 timeout_service에 이 함수를 등록, timeout_service가 시간단위로 이 함수를 실행
func (t TxPeriodicTransferService) TransferTxToLeader() {
	transactions, err := t.txRepository.FindAll()

	if err != nil {
		log.Println(err.Error())
	}

	leader := t.leaderRepository.GetLeader()

	if leader.StringLeaderId() == "" {
		log.Println(errors.New("there is no leader"))
	}

	err = t.grpcService.SendLeaderTransactions(transactions, leader)

	if err != nil {
		log.Println(err.Error())
	}
}
