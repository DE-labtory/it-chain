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

package connection

import (
	"github.com/DE-labtory/it-chain/common/command"
	"github.com/DE-labtory/it-chain/common/rabbitmq/rpc"
	"github.com/DE-labtory/it-chain/conf"
	"github.com/DE-labtory/iLogger"
	"github.com/urfave/cli"
)

func Disconnect() cli.Command {
	return cli.Command{
		Name:  "disconnect",
		Usage: "it-chain connection disconnect [node-id]",
		Action: func(c *cli.Context) error {
			nodeId := c.Args().Get(0)
			return disconnect(nodeId)
		},
	}
}

func disconnect(nodeId string) error {

	config := conf.GetConfiguration()
	client := rpc.NewClient(config.Engine.Amqp)
	defer client.Close()

	disconnectNetworkCommand := command.CloseConnection{
		ConnectionID: nodeId,
	}

	iLogger.Infof(nil, "[Cmd] Disconnecting network - Address: [%s]", nodeId)
	err := client.Call("connection.close", disconnectNetworkCommand, func(connection command.CloseConnection, err rpc.Error) {

		if !err.IsNil() {
			iLogger.Fatalf(nil, "[Cmd] Fail to disconnect to network - Address: [%s], Err: [%s]", connection, err)
			return
		}

		iLogger.Info(nil, "[Cmd] Successfully disconnect to network")
	})

	if err != nil {
		iLogger.Fatal(nil, err.Error())
	}

	return nil
}
