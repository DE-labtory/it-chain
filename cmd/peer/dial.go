package peer

import (
	"fmt"

	"sync"

	"encoding/json"
	"log"

	"github.com/it-chain/it-chain-Engine/conf"
	"github.com/it-chain/it-chain-Engine/messaging/rabbitmq"
	"github.com/it-chain/it-chain-Engine/messaging/rabbitmq/event"
	"github.com/it-chain/it-chain-Engine/messaging/rabbitmq/topic"
	"github.com/streadway/amqp"
	"github.com/urfave/cli"
)

func DialCmd() cli.Command {

	var peerAddress string

	return cli.Command{
		Name:  "dial",
		Usage: "dial to peer",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "address",
				Usage:       "peer address",
				Destination: &peerAddress,
			},
		},
		Action: func(c *cli.Context) error {
			fmt.Println(c.String("address"))
			fmt.Println("peer is dialing...")
			dial(peerAddress)
			return nil
		},
	}
}

//start peer
func dial(peerAddress string) {

	config := conf.GetConfiguration()
	mq := rabbitmq.Connect(config.Common.Messaging.Url)

	wg := sync.WaitGroup{}
	wg.Add(1)

	var tmpReceiver = func(delivery amqp.Delivery) {
		fmt.Println(delivery.Body)
		wg.Done()
	}

	mq.Consume(topic.ConnCreateEvent.String(), tmpReceiver)

	ConnCreatedCmd := event.ConnCreateCmd{}
	ConnCreatedCmd.Address = peerAddress

	b, err := json.Marshal(ConnCreatedCmd)

	if err != nil {
		log.Fatal(err)
	}

	err = mq.Publish(topic.ConnCreateCmd.String(), b)

	if err != nil {
		log.Fatal(err)
	}

	wg.Wait()
}
