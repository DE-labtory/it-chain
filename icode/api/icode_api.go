package api

import (
	"github.com/it-chain/it-chain-Engine/icode/domain/model"
	"github.com/it-chain/it-chain-Engine/icode/domain/repository"
	"github.com/it-chain/it-chain-Engine/icode/domain/service"
)

//ICode의 Invoke, Query, 검증 수행
type ICodeApi struct {
	iCodeMetaRepository repository.ICodeMetaRepository
	containerService    service.ContainerService
	itCodeStoreApi      ItCodeStoreApi
}

//Get ICode from git and start ICode container, push to backup server
func (iApi ICodeApi) Deploy(gitUrl string) error {

	//clone from git
	iCodeMeta, err := iApi.itCodeStoreApi.Clone(gitUrl)

	if err != nil {
		return err
	}

	//start ICode with container
	if err = iApi.containerService.Start(*iCodeMeta); err != nil {
		return err
	}

	//save ICode meta
	if err = iApi.iCodeMetaRepository.Save(*iCodeMeta); err != nil {
		return err
	}

	//push to backup server
	err = iApi.itCodeStoreApi.Push(iCodeMeta)

	if err != nil {
		return err
	}

	return nil
}

//UnDeploy ICode
func (iApi ICodeApi) UnDeploy(id model.ICodeID) error {

	err := iApi.containerService.Stop(id)

	if err != nil {
		return err
	}

	err = iApi.iCodeMetaRepository.Remove(id)

	if err != nil {
		return err
	}

	return nil
}

//Invoke transactions on ICode
func (iApi ICodeApi) Invoke(txs []model.Transaction) {

}

//Query transactions on ICode (Read Only transaction request on ICode)
func (iApi ICodeApi) Query(tx model.Transaction) {

}
