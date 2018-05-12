package peer

import (
	"fmt"

	"github.com/urfave/cli"
)

func DialCmd() cli.Command {

	return cli.Command{
		Name:  "dial",
		Usage: "dial to peer",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "address",
				Usage: "peer address",
			},
		},
		Action: func(c *cli.Context) error {
			fmt.Println(c.String("address"))
			fmt.Println("peer is dialing...")
			dial()
			return nil
		},
	}
}

//start peer
func dial() {
	//gateway.Start()
}
