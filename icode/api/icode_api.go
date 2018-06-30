package api

import (
	"github.com/it-chain/it-chain-Engine/icode"
	"github.com/pkg/errors"
)

type ICodeApi struct {
	ContainerService icode.ContainerService
	StoreApi         ICodeStoreApi
	MetaRepository   icode.ReadOnlyMetaRepository
}

func NewIcodeApi(containerService icode.ContainerService, storeApi ICodeStoreApi, repository icode.ReadOnlyMetaRepository) *ICodeApi {
	return &ICodeApi{
		ContainerService: nil,
		StoreApi:         nil,
		MetaRepository:   nil,
	}
}

func (iApi ICodeApi) Deploy(gitUrl string) error {
	// check for already in repository
	meta, err := iApi.MetaRepository.FindByGitURL(gitUrl)
	if meta.ICodeID != "" {
		return errors.New("Already deployed")
	}
	if err != nil {
		return err
	}

	// clone meta. in clone function, metaCreatedEvent will publish
	meta, err = iApi.StoreApi.Clone(gitUrl)
	if err != nil {
		return err
	}

	//start ICode with container
	if err = iApi.ContainerService.StartContainer(*meta); err != nil {
		return err
	}
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
