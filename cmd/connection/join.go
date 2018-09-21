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
	"github.com/it-chain/engine/common/command"
	"github.com/it-chain/engine/common/logger"
	"github.com/it-chain/engine/common/rabbitmq/rpc"
	"github.com/it-chain/engine/conf"
	"github.com/urfave/cli"
)

func Join() cli.Command {
	return cli.Command{
		Name:  "dial",
		Usage: "it-chain connection join [node-ip-of-the-network]",
		Action: func(c *cli.Context) error {
			nodeIp := c.Args().Get(0)
			return join(nodeIp)
		},
	}
}

func join(ip string) error {

	config := conf.GetConfiguration()
	client := rpc.NewClient(config.Engine.Amqp)
	defer client.Close()

	joinNetworkCommand := command.JoinNetwork{
		Address: ip,
	}

	logger.Infof(nil, "[Cmd] Joining network - Address: [%s]", ip)
	err := client.Call("connection.join", joinNetworkCommand, func(getConnectionListCommand command.GetConnectionList, err rpc.Error) {

		if !err.IsNil() {
			logger.Fatalf(nil, "[Cmd] Fail to join network - Address: [%s]", ip)
			return
		}

		logger.Info(nil, "[Cmd] Successfully joined to network")
		for index, connection := range getConnectionListCommand.ConnectionList {
			logger.Infof(nil, "[Cmd] Connected connection - Index: [%s], ConnectionId: [%s], Address: [%s]", index, connection.ConnectionID, connection.Address)
		}
	})

	if err != nil {
		logger.Fatal(nil, err.Error())
	}

	return nil
}
