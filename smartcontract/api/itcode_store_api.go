package api

import "github.com/it-chain/it-chain-Engine/smartcontract/domain/smartContract"

//Api to import or store itcode from outside
type ItCodeStoreApi interface {

	//get all itcode from repostories
	Init(repositorysUrl string) ([]*smartContract.SmartContract, error)

	//get itcode from outside
	Clone(repositoryUrl string) (*smartContract.SmartContract, error)

	//push code to auth repo
	Push(smartContract.SmartContract) error
}
