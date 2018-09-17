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

package connection

import (
	"fmt"

	"github.com/it-chain/engine/common/command"
	"github.com/it-chain/engine/common/logger"
	"github.com/it-chain/engine/common/rabbitmq/rpc"
	"github.com/it-chain/engine/conf"
	"github.com/urfave/cli"
)

func List() cli.Command {
	return cli.Command{
		Name:  "list",
		Usage: "it-chain connection list",
		Action: func(c *cli.Context) error {
			return list()
		},
	}
}

func list() error {

	config := conf.GetConfiguration()
	client := rpc.NewClient(config.Engine.Amqp)
	defer client.Close()

	getConnectionList := command.GetConnectionList{}
	err := client.Call("connection.list", getConnectionList, func(getConnectionListCommand command.GetConnectionList, err rpc.Error) {

		if !err.IsNil() {
			logger.Fatalf(nil, "[Cmd] Fail to get connection list")
			return
		}

		fmt.Println("Index\t ID\t\t\t\t\t\t Address")
		for index, connection := range getConnectionListCommand.ConnectionList {
			fmt.Printf("[%d]\t [%s]\t [%s]\n",
				index, connection.ConnectionID, connection.Address)
		}
	})

	if err != nil {
		logger.Fatal(nil, err.Error())
	}

	return nil
}
