package peer

import (
	"fmt"

	"github.com/it-chain/it-chain-Engine/gateway"
	"github.com/urfave/cli"
)

func StartCmd() cli.Command {

	return cli.Command{
		Name:    "start",
		Aliases: []string{"s"},
		Usage:   "start peer as background",
		Action: func(c *cli.Context) error {
			fmt.Println("peer is starting...")
			start()
			return nil
		},
	}
}

//start peer
func start() {
	gateway.Start()
}
