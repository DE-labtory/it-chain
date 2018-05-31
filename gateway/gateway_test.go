package gateway_test

import (
	"os"
	"sync"
	"testing"

	"github.com/it-chain/bifrost"
	"github.com/it-chain/it-chain-Engine/gateway"
	"github.com/it-chain/midgard/bus/rabbitmq"
)

func TestStart(t *testing.T) {

	wg := sync.WaitGroup{}
	wg.Add(1)

	amqpUrl := "amqp://guest:guest@localhost:5672/"
	keyPath := "~test"
	grpcUrl := "127.0.0.1:7777"

	go gateway.Start(amqpUrl, grpcUrl, keyPath)

	defer func() {
		os.RemoveAll(keyPath)
		os.RemoveAll("~test2")
		gateway.Stop()
	}()

	wg.Wait()
}

func TestStart2(t *testing.T) {

	wg := sync.WaitGroup{}
	wg.Add(1)

	amqpUrl := "amqp://guest:guest@localhost:5672/"
	keyPath := "~test3"
	rabbitmqClient := rabbitmq.Connect(amqpUrl)

	defer func() {
		os.RemoveAll(keyPath)
		rabbitmqClient.Close()
	}()

	//create connection store by bifrost which is it-chain's own lib for implementing p2p network
	ConnectionStore := bifrost.NewConnectionStore()

	pri, pub := gateway.LoadKeyPair(keyPath)
	//load key

	//createHandler
	connectionHandler := gateway.NewConnectionCommandHandler(ConnectionStore, pri, pub, rabbitmqClient)
	connectionHandler.HandleConnectionCreate(gateway.ConnectionCreateCommand{
		Address: "127.0.0.1:7777",
	})

	wg.Wait()
}

func TestStart3(t *testing.T) {

	wg := sync.WaitGroup{}
	wg.Add(1)

	amqpUrl := "amqp://guest:guest@localhost:5672/"
	keyPath := "~test4"
	rabbitmqClient := rabbitmq.Connect(amqpUrl)

	defer func() {
		os.RemoveAll(keyPath)
		rabbitmqClient.Close()
	}()

	//create connection store by bifrost which is it-chain's own lib for implementing p2p network
	ConnectionStore := bifrost.NewConnectionStore()

	pri, pub := gateway.LoadKeyPair(keyPath)
	//load key

	//createHandler
	connectionHandler := gateway.NewConnectionCommandHandler(ConnectionStore, pri, pub, rabbitmqClient)
	connectionHandler.HandleConnectionCreate(gateway.ConnectionCreateCommand{
		Address: "127.0.0.1:7777",
	})

	wg.Wait()
}
