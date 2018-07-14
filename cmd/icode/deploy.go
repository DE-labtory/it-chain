package icode

import (
	"github.com/it-chain/it-chain-Engine/conf"
	"github.com/it-chain/it-chain-Engine/icode"
	"github.com/it-chain/midgard/bus/rabbitmq"
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
		Url:     gitUrl,
		SshPath: sshPath,
	}
	client.Publish("Command", "icode.deploy", command)
}
