package peer

import (
	"github.com/urfave/cli"
)

var peerCmd = cli.Command{
	Name:        "peer",
	Aliases:     []string{"p"},
	Usage:       "options for peer",
	Subcommands: []cli.Command{},
}

func PeerCmd() cli.Command {
	peerCmd.Subcommands = append(peerCmd.Subcommands, StartCmd())
	peerCmd.Subcommands = append(peerCmd.Subcommands, DialCmd())
	return peerCmd
}
