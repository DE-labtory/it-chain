package icode

import (
	"github.com/it-chain/engine/conf"
	"github.com/it-chain/engine/icode"
	"github.com/it-chain/midgard"
	"github.com/it-chain/midgard/bus/rabbitmq"
	"github.com/urfave/cli"
)

func UnDeployCmd() cli.Command {
	return cli.Command{
		Name:  "undeploy",
		Usage: "it-chain icode undeploy [icode id] ",
		Action: func(c *cli.Context) error {

			icodeId := c.Args().Get(0)
			unDeploy(icodeId)
			return nil
		},
	}
}
func unDeploy(icodeId string) {
	config := conf.GetConfiguration()
	client := rabbitmq.Connect(config.Engine.Amqp)
	defer client.Close()
	command := icode.UnDeployCommand{
		CommandModel: midgard.CommandModel{
			ID: icodeId,
		},
	}
	client.Publish("Command", "icode.undeploy", command)
}
