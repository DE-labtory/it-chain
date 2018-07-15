package blockchain

import "github.com/urfave/cli"

var blockchainCmd = cli.Command{
	Name:        "blockchain",
	Aliases:     []string{"p"},
	Usage:       "options for blockchain",
	Subcommands: []cli.Command{},
}

func BlockchainCmd() cli.Command {
	blockchainCmd.Subcommands = append(blockchainCmd.Subcommands, ProposeBlockCmd())
	return blockchainCmd
}
