package service

import "it-chain/domain"

//peer 최상위 service
type SmartContractService interface{
	Invoke()
	Query(transaction domain.Transaction) (error)
	Deploy(ReposPath string) (string, error)
	PullAllSmartContracts(authenticatedGit string, errorHandler func(error), completionHandler func())
	ValidateTransactionsInBlock(block *domain.Block) error
}
