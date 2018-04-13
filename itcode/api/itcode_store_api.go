package api

import "github.com/it-chain/it-chain-Engine/smartcontract/domain/itcode"

//Api to import or store itcode from outside
type ItCodeStoreApi interface {

	//get itcode from outside
	Clone(repositoryUrl string) (*itcode.ItCode, error)

	//push code to auth repo
	Push(itcode.ItCode) error
}
