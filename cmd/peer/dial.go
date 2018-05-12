package peer

import (
	"fmt"

	"sync"

	"encoding/json"
	"log"

	"github.com/it-chain/it-chain-Engine/conf"
	"github.com/it-chain/it-chain-Engine/gateway"
	"github.com/it-chain/it-chain-Engine/messaging/rabbitmq"
	"github.com/it-chain/it-chain-Engine/messaging/rabbitmq/event"
	"github.com/it-chain/it-chain-Engine/messaging/rabbitmq/topic"
	"github.com/streadway/amqp"
	"github.com/urfave/cli"
)

func DialCmd() cli.Command {

	return cli.Command{
		Name:  "dial",
		Usage: "it-chain peer dial [ip]",
		Action: func(c *cli.Context) error {

			address := c.Args().Get(0)
			fmt.Printf("peer is dialing on [%s]", address)
			dial(address)

			return nil
		},
	}
}

//start peer
func dial(peerAddress string) {

	config := conf.GetConfiguration()
	mq := rabbitmq.Connect(config.Common.Messaging.Url)

	defer mq.Close()

	wg := sync.WaitGroup{}
	wg.Add(1)

	var tmpConnCmdCreateReceiver = func(delivery amqp.Delivery) {
		fmt.Println(delivery.Body)
		wg.Done()
	}

	var tmpGatewayErrorEventReceiver = func(delivery amqp.Delivery) {

		gatewayErrorEvent := &gateway.ErrorEvent{}

		err := json.Unmarshal(delivery.Body, gatewayErrorEvent)

		if err != nil {

		}

		if gatewayErrorEvent.Event == "ConnCreateCmd" {
			fmt.Printf("fail to dial peer [%s]", gatewayErrorEvent.Err)
			wg.Done()
		}
	}

	mq.Consume(topic.ConnCreateEvent.String(), tmpConnCmdCreateReceiver)
	mq.Consume("GatewayErrorEvent", tmpGatewayErrorEventReceiver)

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
