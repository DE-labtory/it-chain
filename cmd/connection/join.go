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

func Join() cli.Command {
	return cli.Command{
		Name:  "join",
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

	iLogger.Infof(nil, "[Cmd] Joining network - Address: [%s]", ip)
	err := client.Call("connection.join", joinNetworkCommand, func(_ struct{}, err rpc.Error) {

		if !err.IsNil() {
			iLogger.Fatalf(nil, "[Cmd] Fail to join network - Address: [%s], Err: [%s]", ip, err.Message)
			return
		}

		iLogger.Info(nil, "[Cmd] Successfully request to join network")
	})

	if err != nil {
		iLogger.Fatal(nil, err.Error())
	}

	return nil
}
