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

	"github.com/it-chain/engine/common/command"
	"github.com/it-chain/engine/common/rabbitmq/rpc"
	"github.com/it-chain/engine/conf"
	"github.com/it-chain/engine/ivm"
	"github.com/rs/xid"
	"github.com/urfave/cli"
)

func DeployCmd() cli.Command {
	return cli.Command{
		Name:  "deploy",
		Usage: "it-chain ivm deploy [icode git url] [ssh path] [password]",
		Action: func(c *cli.Context) error {

			gitUrl := c.Args().Get(0)
			sshPath := c.Args().Get(1)
			password := c.Args().Get(2)
			deploy(gitUrl, sshPath, password)

			return nil
		},
	}
}

func deploy(gitUrl string, sshPath string, password string) {

	config := conf.GetConfiguration()
	client := rpc.NewClient(config.Engine.Amqp)

	defer client.Close()

	deployCommand := command.Deploy{
		ICodeId:  xid.New().String(),
		Url:      gitUrl,
		SshPath:  sshPath,
		Password: password,
	}

	log.Printf("[Cmd] deploying icode...")
	log.Printf("[Cmd] This may take a few minutes")

	err := client.Call("ivm.deploy", deployCommand, func(icode ivm.ICode, err rpc.Error) {

		if !err.IsNil() {
			log.Printf("fail to deploy icode err: [%s]", err.Message)
			return
		}

		log.Printf("[Cmd] icode has deployed - icodeID: [%s]", icode.ID)
	})

	if err != nil {
		log.Fatal(err.Error())
	}
}
