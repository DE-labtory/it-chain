package peer

import (
	"fmt"
	"syscall"

	"log"

	"github.com/urfave/cli"
)

func StopCmd() cli.Command {

	return cli.Command{
		Name:    "stop",
		Aliases: []string{"t"},
		Usage:   "stop peer",
		Action: func(c *cli.Context) error {
			fmt.Println("peer is terminating...")
			stop()
			return nil
		},
	}
}

//todo delete my.pid?
func stop() {

	pid, err := GetValue("my.pid")

	if err != nil {
		log.Fatalln(err.Error())
	}

	err = syscall.Kill(pid, syscall.SIGKILL)

	if err != nil {
		log.Fatalln(err.Error())
	}

	log.Printf("it-chain peer is terminating")
}
