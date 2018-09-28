/*
 * Copyright 2018 It-chain
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package api_gateway

import (
	"log"

	"github.com/it-chain/engine/common/command"
	"github.com/it-chain/engine/common/rabbitmq/rpc"
	"github.com/it-chain/engine/conf"
	"github.com/it-chain/engine/ivm"
	"github.com/it-chain/engine/txpool"
	"github.com/it-chain/iLogger"
	"github.com/pkg/errors"
	"github.com/rs/xid"
)

type ICodeCommandApi struct {
}

func NewICodeCommandApi() *ICodeCommandApi {
	return &ICodeCommandApi{}
}

func (i *ICodeCommandApi) deploy(amqpUrl string, gitUrl string, rawSsh string, sshPassword string) (string, error) {
	if amqpUrl == "" {
		config := conf.GetConfiguration()
		amqpUrl = config.Engine.Amqp
	}

	client := rpc.NewClient(amqpUrl)

	defer client.Close()

	deployCommand := command.Deploy{
		Url:      gitUrl,
		SshRaw:   []byte(rawSsh),
		Password: sshPassword,
	}

	iLogger.Infof(nil, "[Api_gateway] deploying icode...")
	iLogger.Infof(nil, "[Api_gateway] This may take a few minutes")

	var callBackIcodeId ivm.ID
	var callBackErr error

	err := client.Call("ivm.deploy", deployCommand, func(icode ivm.ICode, err rpc.Error) {

		if !err.IsNil() {
			iLogger.Infof(nil, "[Api_gateway] fail to deploy icode err: [%s]", err.Message)
			callBackErr = errors.New(err.Message)
			return
		}

		iLogger.Infof(nil, "[Api_gateway] icode has deployed - icodeID: [%s]", icode.ID)
		callBackErr = nil
		callBackIcodeId = icode.ID
	})

	if err != nil {
		iLogger.Fatal(&iLogger.Fields{"err_msg": err.Error()}, "[Api_gateway] fatal err in deploy cmd")
		return "", err
	}

	if callBackErr != nil {
		return "", callBackErr
	}

	return callBackIcodeId, nil
}

func (i *ICodeCommandApi) unDeploy(amqpUrl string, icodeId string) error {
	if amqpUrl == "" {
		config := conf.GetConfiguration()
		amqpUrl = config.Engine.Amqp
	}

	client := rpc.NewClient(amqpUrl)

	defer client.Close()

	undeployCommand := command.UnDeploy{
		ICodeId: icodeId,
	}

	var callBackErr error

	err := client.Call("ivm.undeploy", undeployCommand, func(empty struct{}, err rpc.Error) {

		if !err.IsNil() {
			log.Printf("[Api_gateway] fail to undeploy icode err: [%s]", err.Message)
			callBackErr = errors.New(err.Message)
			return
		}

		log.Printf("[Api_gateway] [%s] icode has undeployed", icodeId)
		callBackErr = nil
	})

	if err != nil {
		iLogger.Fatal(&iLogger.Fields{"err_msg": err.Error()}, "[Api_gateway] fatal err in unDeploy cmd")
		return err
	}

	if callBackErr != nil {
		return callBackErr
	}

	return nil
}

func (i *ICodeCommandApi) invoke(amqpUrl string, id string, functionName string, args []string) (string, error) {
	if amqpUrl == "" {
		config := conf.GetConfiguration()
		amqpUrl = config.Engine.Amqp
	}

	client := rpc.NewClient(amqpUrl)

	defer client.Close()

	invokeCommand := command.CreateTransaction{
		TransactionId: xid.New().String(),
		ICodeID:       id,
		Jsonrpc:       "2.0",
		Method:        "invoke",
		Args:          args,
		Function:      functionName,
	}

	iLogger.Infof(nil, "[Api_gateway] Invoke icode - icodeID: [%s]", id)

	var callBackTransactionId txpool.TransactionId
	var callBackErr error

	err := client.Call("transaction.create", invokeCommand, func(transaction txpool.Transaction, err rpc.Error) {

		if !err.IsNil() {
			iLogger.Errorf(nil, "[Api_gateway] Fail to invoke icode err: [%s]", err.Message)
			callBackErr = errors.New(err.Message)
			return
		}

		iLogger.Infof(nil, "[Api_gateway] Transactions are created - ID: [%s]", transaction.ID)
		callBackErr = nil
		callBackTransactionId = transaction.ID

	})

	if err != nil {
		iLogger.Fatal(&iLogger.Fields{"err_msg": err.Error()}, "[Api_gateway] fatal err in invoke cmd")
		return "", err
	}

	if callBackErr != nil {
		return "", callBackErr
	}

	return callBackTransactionId, nil
}

func (i *ICodeCommandApi) query(amqpUrl string, id string, functionName string, args []string) (map[string]string, error) {
	if amqpUrl == "" {
		config := conf.GetConfiguration()
		amqpUrl = config.Engine.Amqp
	}

	client := rpc.NewClient(amqpUrl)

	defer client.Close()

	queryCommand := command.ExecuteICode{
		ICodeId:  id,
		Function: functionName,
		Args:     args,
		Method:   "query",
	}

	iLogger.Infof(nil, "[Api_gateway] Querying icode - icodeID: [%s], func: [%s]", id, functionName)

	var callBackResult map[string]string
	var callBackErr error

	err := client.Call("ivm.execute", queryCommand, func(result ivm.Result, err rpc.Error) {
		if !err.IsNil() {
			iLogger.Errorf(nil, "[Api_gateway] Fail to query icode err: [%s]", err.Message)
			callBackErr = errors.New(err.Message)
			return
		}

		for key, val := range result.Data {
			iLogger.Infof(nil, "[Api_gateway] Querying result - key: [%s], value: [%s]", key, val)
		}

		callBackResult = result.Data
		callBackErr = nil
		if result.Err != "" {
			callBackErr = errors.New(result.Err)
		}
	})

	if err != nil {
		iLogger.Fatal(&iLogger.Fields{"err_msg": err.Error()}, "[Api_gateway] fatal err in query cmd")
		return nil, err
	}

	if callBackErr != nil {
		return nil, callBackErr
	}

	return callBackResult, nil
}
