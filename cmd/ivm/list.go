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

package ivm

import (
	"log"

	"fmt"

	"github.com/it-chain/engine/common/command"
	"github.com/it-chain/engine/common/rabbitmq/rpc"
	"github.com/it-chain/engine/conf"
	"github.com/urfave/cli"
)

func ListCmd() cli.Command {
	return cli.Command{
		Name:  "list",
		Usage: "it-chain ivm list",
		Action: func(c *cli.Context) error {

			list()
			return nil
		},
	}
}

func list() {

	config := conf.GetConfiguration()
	client := rpc.NewClient(config.Engine.Amqp)
	defer client.Close()

	listCommand := command.GetICodeList{}
	err := client.Call("ivm.list", listCommand, func(ICodeList command.ICodeList, err rpc.Error) {

		if !err.IsNil() {
			log.Printf("fail to get icode list err: [%s]", err.Message)
			return
		}
		
		fmt.Println("Index\t ID\t\t\t Version\t GitUrl")
		for index, iCodes := range ICodeList.ICodes {
			fmt.Printf("[%d]\t [%s]\t [%s]\t\t [%s]\n",
				index, iCodes.ID, iCodes.Version, iCodes.GitUrl)
		}
	})

	if err != nil {
		log.Fatal(err.Error())
	}
}
