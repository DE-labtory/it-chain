/*
 * Copyright 2018 DE-labtory
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

package ivm

import (
	"errors"

	"github.com/DE-labtory/it-chain/common/command"
	"github.com/DE-labtory/it-chain/common/rabbitmq/rpc"
	"github.com/DE-labtory/it-chain/conf"
	"github.com/DE-labtory/it-chain/txpool"
	"github.com/DE-labtory/iLogger"
	"github.com/rs/xid"
	"github.com/urfave/cli"
)

func InvokeCmd() cli.Command {
	return cli.Command{
		Name:  "invoke",
		Usage: "it-chain ivm invoke [icode-id] [function-name] [...args]",
		Action: func(c *cli.Context) error {
			if c.NArg() < 2 {
				return errors.New("not enough args")
			}

			icodeId := c.Args().Get(0)
			functionName := c.Args().Get(1)
			args := make([]string, 0)

			for i := 2; i < c.NArg(); i++ {
				args = append(args, c.Args().Get(i))
			}

			invoke(icodeId, functionName, args)

			return nil
		},
	}
}

func invoke(id string, functionName string, args []string) {

	config := conf.GetConfiguration()
	client := rpc.NewClient(config.Engine.Amqp)

	defer client.Close()

	invokeCommand := command.CreateTransaction{
		TransactionId: xid.New().String(),
		ICodeID:       id,
		Jsonrpc:       "2.0",
		Method:        "invoke",
		Args:          args,
		Function:      functionName,
	}

	iLogger.Infof(nil, "[Cmd] Invoke icode - icodeID: [%s]", id)

	err := client.Call("transaction.create", invokeCommand, func(transaction txpool.Transaction, err rpc.Error) {

		if !err.IsNil() {
			iLogger.Errorf(nil, "[Cmd] Fail to invoke icode err: [%s]", err.Message)
			return
		}

		iLogger.Infof(nil, "[Cmd] Transactions are created - ID: [%s]", transaction.ID)
	})

	if err != nil {
		iLogger.Fatal(&iLogger.Fields{"err_msg": err.Error()}, "fatal err in query cmd")
	}
}
