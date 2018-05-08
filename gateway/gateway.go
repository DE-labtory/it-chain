package main

import (
	"log"

	"github.com/it-chain/bifrost"
	"github.com/it-chain/bifrost/mux"
	"github.com/it-chain/bifrost/server"
	"github.com/it-chain/heimdall/key"
	"github.com/it-chain/it-chain-Engine/conf"
	"github.com/it-chain/it-chain-Engine/messaging"
	"github.com/it-chain/it-chain-Engine/messaging/topic"
)

var DefaultMux *mux.DefaultMux
var ConnectionStore *bifrost.ConnectionStore

//set mux for connection
func init() {

	DefaultMux = mux.New()

	DefaultMux.Handle("chat", func(message bifrost.Message) {
		log.Printf("%s", message.Data)
	})

	DefaultMux.Handle("join", func(message bifrost.Message) {
		log.Printf("%s", message.Data)
	})

	ConnectionStore = bifrost.NewConnectionStore()
}

func Start() error {

	config := conf.GetConfiguration()

	mq := messaging.NewRabbitmq(config.Common.Messaging.Url)
	mq.Start()

	eventConsumer := NewEventConsumer(ConnectionStore)

	mq.Consume(topic.MessageDeliverEvent.String(), eventConsumer.HandleMessageDeliverEvent)

	pri, pub := loadKeyPair(config.Authentication.KeyPath)

	s := server.New(bifrost.KeyOpts{PriKey: pri, PubKey: pub})

	s.OnConnection(OnConnection)
	s.OnError(OnError)

	s.Listen(config.GrpcGateway.Ip)

	return nil
}

func OnConnection(connection bifrost.Connection) {

	connection.Handle(DefaultMux)
	ConnectionStore.AddConnection(connection)

	defer connection.Close()

	if err := connection.Start(); err != nil {
		connection.Close()
	}
}

func OnError(err error) {
	log.Fatalln(err.Error())
}

func loadKeyPair(keyPath string) (key.PriKey, key.PubKey) {

	km, err := key.NewKeyManager(keyPath)

	if err != nil {
		log.Fatal(err.Error())
	}

	pri, pub, err := km.GetKey()

	if err == nil {
		return pri, pub
	}

	pri, pub, err = km.GenerateKey(key.KeyGenOpts(conf.GetConfiguration().Authentication.KeyType))

	if err != nil {
		log.Fatal(err.Error())
	}

	return pri, pub
}
