package api

import (
	"errors"

	"fmt"

	"github.com/it-chain/it-chain-Engine/icode"
)

type ICodeApi struct {
	ContainerService icode.ContainerService
	StoreApi         ICodeStoreApi
	MetaRepository   icode.ReadOnlyMetaRepository
}

func NewIcodeApi(containerService icode.ContainerService, storeApi ICodeStoreApi, repository icode.ReadOnlyMetaRepository) *ICodeApi {
	return &ICodeApi{
		ContainerService: containerService,
		StoreApi:         storeApi,
		MetaRepository:   repository,
	}
}

func (iApi ICodeApi) Deploy(baseSaveUrl string, gitUrl string) (*icode.Meta, error) {
	// check for already in repository
	meta, err := iApi.MetaRepository.FindByGitURL(gitUrl)
	if meta.ICodeID != "" {
		return nil, errors.New("already deployed")
	}
	if err != nil {
		return nil, err
	}

	// clone meta. in clone function, metaCreatedEvent will publish
	meta, err = iApi.StoreApi.Clone(baseSaveUrl, gitUrl)
	if err != nil {
		return nil, err
	}

	//start ICode with container
	if err = iApi.ContainerService.StartContainer(*meta); err != nil {
		return nil, err
	}
	return meta, nil
}

func (iApi ICodeApi) UnDeploy(id icode.ID) error {
	// stop iCode container
	err := iApi.ContainerService.StopContainer(id)
	if err != nil {
		return err
	}

	// publish meta delete event
	err = icode.DeleteMeta(id)
	if err != nil {
		return err
	}

	return nil
}

//todo need asnyc process
func (iApi ICodeApi) Invoke(txs []icode.Transaction) []icode.Result {
	resultData := make([]icode.Result, 0)
	for _, tx := range txs {
		result, err := iApi.ContainerService.ExecuteTransaction(tx)
		if err != nil {
			fmt.Println(fmt.Sprintf("error in invoke tx, err : %s", err.Error()))
			result = &icode.Result{
				TxId:    tx.TxId,
				Data:    nil,
				Success: false,
			}
		}
		resultData = append(resultData, *result)
	}
	return resultData
}

func (iApi ICodeApi) Query(tx icode.Transaction) (*icode.Result, error) {
	result, err := iApi.ContainerService.ExecuteTransaction(tx)

	if err != nil {
		return nil, err
	}

	return result, nil
}
