/*
 * Copyright 2018 It-chain
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package blockchain

import "github.com/urfave/cli"

var blockchainCmd = cli.Command{
	Name:        "blockchain",
	Aliases:     []string{"b"},
	Usage:       "options for blockchain",
	Subcommands: []cli.Command{},
}

func BlockchainCmd() cli.Command {
	blockchainCmd.Subcommands = append(blockchainCmd.Subcommands, CreateGenesisBlock())
	return blockchainCmd
}
