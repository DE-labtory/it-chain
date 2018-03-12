package service

import "github.com/it-chain/it-chain-Engine/legacy/domain"

//peer 최상위 service
type SmartContractService interface{

	Invoke(transaction *domain.Transaction) (*domain.SmartContractResponse, error)
	Query()
	Deploy(ReposPath string) (string, error)
	PullAllSmartContracts(errorHandler func(error), completionHandler func())
	ValidateTransactionsOfBlock(block *domain.Block) (error)
	ValidateTransaction(transaction *domain.Transaction)
}
