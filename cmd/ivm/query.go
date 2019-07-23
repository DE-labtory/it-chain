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
	"github.com/DE-labtory/it-chain/ivm"
	"github.com/DE-labtory/iLogger"
	"github.com/urfave/cli"
)

func QueryCmd() cli.Command {
	return cli.Command{
		Name:  "query",
		Usage: "it-chain ivm query [icode-id] [functioniname] [...args]",
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
			query(icodeId, functionName, args)

			return nil
		},
	}
}

func query(id string, functionName string, args []string) {

	config := conf.GetConfiguration()
	client := rpc.NewClient(config.Engine.Amqp)

	defer client.Close()

	queryCommand := command.ExecuteICode{
		ICodeId:  id,
		Function: functionName,
		Args:     args,
		Method:   "query",
	}

	iLogger.Infof(nil, "[Cmd] Querying icode - icodeID: [%s], func: [%s]", id, functionName)

	err := client.Call("ivm.execute", queryCommand, func(result ivm.Result, err rpc.Error) {
		if !err.IsNil() {
			iLogger.Errorf(nil, "[Cmd] Fail to query icode err: [%s]", err.Message)
			return
		}

		for key, val := range result.Data {
			iLogger.Infof(nil, "[CMD] Querying result - key: [%s], value: [%s]", key, val)
		}
	})

	if err != nil {
		iLogger.Fatal(&iLogger.Fields{"err_msg": err.Error()}, "fatal err in query cmd")
	}
}
