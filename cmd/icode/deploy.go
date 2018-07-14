package icode

import (
	"fmt"

	"github.com/it-chain/it-chain-Engine/conf"
	"github.com/it-chain/it-chain-Engine/icode"
	"github.com/it-chain/midgard"
	"github.com/it-chain/midgard/bus/rabbitmq"
	"github.com/rs/xid"
	"github.com/urfave/cli"
)

func DeployCmd() cli.Command {
	return cli.Command{
		Name:  "deploy",
		Usage: "it-chain icode deploy [icode git url] [ssh path]",
		Action: func(c *cli.Context) error {

			gitUrl := c.Args().Get(0)
			sshPath := c.Args().Get(1)
			deploy(gitUrl, sshPath)

			return nil
		},
	}
}
func deploy(gitUrl string, sshPath string) {
	config := conf.GetConfiguration()
	client := rabbitmq.Connect(config.Common.Messaging.Url)
	defer client.Close()
	command := icode.DeployCommand{
		CommandModel: midgard.CommandModel{
			ID: xid.New().String(),
		},
		Url:     gitUrl,
		SshPath: sshPath,
	}
	fmt.Println(fmt.Sprintf("deploying ID : %s", command.GetID()))
	client.Publish("Command", "icode.deploy", command)
}
