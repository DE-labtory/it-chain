package icode

import "github.com/urfave/cli"

var icodeCmd = cli.Command{
	Name:        "icode",
	Aliases:     []string{"p"},
	Usage:       "options for icode",
	Subcommands: []cli.Command{},
}

func IcodeCmd() cli.Command {
	icodeCmd.Subcommands = append(icodeCmd.Subcommands, DeployCmd())
	icodeCmd.Subcommands = append(icodeCmd.Subcommands, UnDeployCmd())
	return icodeCmd
}
