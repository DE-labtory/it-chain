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
	"github.com/DE-labtory/engine/common/command"
	"github.com/DE-labtory/engine/common/rabbitmq/rpc"
	"github.com/DE-labtory/engine/conf"
	"github.com/DE-labtory/engine/grpc_gateway"
	"github.com/DE-labtory/iLogger"
	"github.com/urfave/cli"
)

func Dial() cli.Command {
	return cli.Command{
		Name:  "dial",
		Usage: "it-chain connection dial [node-ip]",
		Action: func(c *cli.Context) error {
			nodeIp := c.Args().Get(0)
			return dial(nodeIp)
		},
	}
}

func dial(ip string) error {

	config := conf.GetConfiguration()
	client := rpc.NewClient(config.Engine.Amqp)
	defer client.Close()

	createConnectionCommand := command.CreateConnection{
		Address: ip,
	}

	iLogger.Infof(nil, "[Cmd] Creating connection - Address: [%s]", ip)
	err := client.Call("connection.create", createConnectionCommand, func(connection grpc_gateway.Connection, err rpc.Error) {

		if !err.IsNil() {
			iLogger.Fatalf(nil, "[Cmd] Fail to create connection - Address: [%s]", ip)
			return
		}

		iLogger.Infof(nil, "[Cmd] Connection created - gRPC-Address: [%s], Id:[%s]", connection.GrpcGatewayAddress, connection.ConnectionID)
	})

	if err != nil {
		iLogger.Fatal(nil, err.Error())
	}

	return nil
}
