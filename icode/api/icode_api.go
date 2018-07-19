package api

import (
	"fmt"

	"github.com/it-chain/engine/icode"
)

type ICodeApi struct {
	ContainerService icode.ContainerService
	StoreApi         ICodeStoreApi
}

func NewIcodeApi(containerService icode.ContainerService, storeApi ICodeStoreApi) *ICodeApi {
	return &ICodeApi{
		ContainerService: containerService,
		StoreApi:         storeApi,
	}
}

func (iApi ICodeApi) Deploy(id string, baseSaveUrl string, gitUrl string, sshPath string) (*icode.Meta, error) {
	// check for already in repository
	/*meta, err := iApi.MetaRepository.FindByGitURL(gitUrl)
	if meta.ICodeID != "" {
		return nil, errors.New("already deployed")
	}
	if err != nil {
		return nil, err
	}
	*/

	// clone meta. in clone function, metaCreatedEvent will publish
	meta, err := iApi.StoreApi.Clone(id, baseSaveUrl, gitUrl, sshPath)
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
func (iApi ICodeApi) Invoke(tx icode.Transaction) *icode.Result {
	result, err := iApi.ContainerService.ExecuteTransaction(tx)
	if err != nil {
		fmt.Println(fmt.Sprintf("error in invoke tx, err : %s", err.Error()))
		result = &icode.Result{
			TxId:    tx.TxId,
			Data:    nil,
			Success: false,
		}
	}
	return result
}

func (iApi ICodeApi) Query(tx icode.Transaction) *icode.Result {
	result, err := iApi.ContainerService.ExecuteTransaction(tx)
	if err != nil {
		fmt.Println(fmt.Sprintf("error in invoke tx, err : %s", err.Error()))
		result = &icode.Result{
			TxId:    tx.TxId,
			Data:    nil,
			Success: false,
		}
	}
	return result
}
