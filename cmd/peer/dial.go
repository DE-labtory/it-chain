package peer

import (
	"fmt"

	"sync"

	"log"

	"reflect"

	"github.com/it-chain/it-chain-Engine/conf"
	"github.com/it-chain/it-chain-Engine/grpc_gateway"
	"github.com/it-chain/midgard/bus/rabbitmq"
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

var wg = sync.WaitGroup{}

type ErrorEventHandler struct {
}

func (e ErrorEventHandler) ErrorCreated(error grpc_gateway.ErrorCreatedEvent) {
	fmt.Println(error)
	//wg.Done()
}

type ConnectionEventHandler struct {
}

func (e ErrorEventHandler) ConnectionCreated(event grpc_gateway.ConnectionCreatedEvent) {
	fmt.Println(event)
	//wg.Done()
}

//start peer
func dial(peerAddress string) {

	config := conf.GetConfiguration()
	client := rabbitmq.Connect(config.Common.Messaging.Url)

	defer client.Close()

	client.Subscribe("Event", "Error", &ErrorEventHandler{})
	client.Subscribe("Event", "Connection", &ConnectionEventHandler{})

	command := grpc_gateway.ConnectionCreateCommand{
		Address: peerAddress,
	}

	log.Println(reflect.TypeOf(command))

	wg.Add(1)
	err := client.Publish("Command", "Connection", command)

	if err != nil {
		log.Fatal(err)
	}

	wg.Wait()
}
