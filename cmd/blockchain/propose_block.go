package blockchain

import (
	"time"

	"fmt"

	"github.com/it-chain/it-chain-Engine/blockchain"
	"github.com/it-chain/it-chain-Engine/conf"
	"github.com/it-chain/it-chain-Engine/txpool"
	"github.com/it-chain/midgard"
	"github.com/it-chain/midgard/bus/rabbitmq"
	"github.com/rs/xid"
	"github.com/urfave/cli"
)

func ProposeBlockCmd() cli.Command {
	return cli.Command{
		Name:  "propose-block",
		Usage: "it-chain blockchain propose-block",
		Action: func(c *cli.Context) error {
			proposeBlock()
			return nil
		},
	}
}

func proposeBlock() {
	config := conf.GetConfiguration()
	client := rabbitmq.Connect(config.Common.Messaging.Url)
	defer client.Close()

	command := blockchain.ProposeBlockCommand{
		CommandModel: midgard.CommandModel{
			ID: xid.New().String(),
		},
		Transactions: getTxList(getTime()),
	}

	fmt.Println(fmt.Sprintf("commmand ID : %s", command.GetID()))
	client.Publish("Command", "block.propose", command)
}

func getTxList(testingTime time.Time) []txpool.Transaction {
	return []txpool.Transaction{
		{
			PublishPeerId: "p01",
			TxId:          "tx01",
			TimeStamp:     testingTime,
			TxData: txpool.TxData{
				Jsonrpc: "jsonRPC01",
				Method:  "invoke",
				Params: txpool.Param{
					Function: "function01",
					Args:     []string{"arg1", "arg2"},
				},
				ID: "txdata01",
			},
		},
		{
			PublishPeerId: "p02",
			TxId:          "tx02",
			TimeStamp:     testingTime,
			TxData: txpool.TxData{
				Jsonrpc: "jsonRPC02",
				Method:  "invoke",
				Params: txpool.Param{
					Function: "function02",
					Args:     []string{"arg1", "arg2"},
				},
				ID: "txdata02",
			},
		},
		{
			PublishPeerId: "p03",
			TxId:          "tx03",
			TimeStamp:     testingTime,
			TxData: txpool.TxData{
				Jsonrpc: "jsonRPC03",
				Method:  "invoke",
				Params: txpool.Param{
					Function: "function03",
					Args:     []string{"arg1", "arg2"},
				},
				ID: "txdata03",
			},
		}, {
			PublishPeerId: "p04",
			TxId:          "tx04",
			TimeStamp:     testingTime,
			TxData: txpool.TxData{
				Jsonrpc: "jsonRPC04",
				Method:  "invoke",
				Params: txpool.Param{
					Function: "function04",
					Args:     []string{"arg1", "arg2"},
				},
				ID: "txdata04",
			},
		},
	}
}

func getTime() time.Time {
	testingTime, _ := time.Parse("Jan 2, 2006 at 3:04pm (MST)", "Feb 3, 2013 at 7:54pm (UTC)")
	return testingTime
}
