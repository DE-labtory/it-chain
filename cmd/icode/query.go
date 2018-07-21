package icode

import (
	"fmt"

	"github.com/it-chain/engine/common/amqp/rpc"
	"github.com/it-chain/engine/common/command"
	"github.com/it-chain/engine/conf"
	"github.com/it-chain/engine/icode"
	"github.com/it-chain/midgard"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"github.com/urfave/cli"
)

func QueryCmd() cli.Command {
	return cli.Command{
		Name:  "deploy",
		Usage: "it-chain icode query [icode id] [function name] [args...]",
		Action: func(c *cli.Context) error {

			if c.NArg() < 3 {
				return errors.New("not enough args")
			}
			icodeId := c.Args().Get(0)
			functionName := c.Args().Get(1)
			args := make([]string, 0)
			for i := 2; i < c.NArg(); i++ {
				args = append(args, c.Args().Get(i))
			}
			query(icodeId, functionName, args)

			return nil
		},
	}
}
func query(id string, functionName string, args []string) {
	config := conf.GetConfiguration()
	conn, err := amqp.Dial(config.Engine.Amqp)
	defer conn.Close()
	if err != nil {
		fmt.Println("error while create amqp conn")
	}
	client, err := rpc.NewRpcClient(*conn, "icode")
	if err != nil {
		fmt.Println("error while create rpc client")
	}
	result := new(icode.Result)
	comm := command.Query{
		CommandModel: midgard.CommandModel{
			ID: id,
		},
		Function: functionName,
		Args:     args,
	}
	err = client.Call("ICodeApi.Query", comm, result)
	if err != nil {
		fmt.Println("error while call query")
	}
	fmt.Println(result)
}
