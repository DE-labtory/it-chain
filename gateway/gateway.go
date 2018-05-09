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
var pri key.PriKey
var pub key.PubKey
var config *conf.Configuration
var mq *messaging.Rabbitmq

//set mux for connection
func init() {

	config := conf.GetConfiguration()

	DefaultMux = mux.New()

	DefaultMux.Handle("chat", func(message bifrost.Message) {
		log.Printf("%s", message.Data)
	})

	DefaultMux.Handle("join", func(message bifrost.Message) {
		log.Printf("%s", message.Data)
	})

	ConnectionStore = bifrost.NewConnectionStore()
	pri, pub = loadKeyPair(config.Authentication.KeyPath)
	mq = messaging.NewRabbitmq(config.Common.Messaging.Url)
	mq.Start()
}

func Start() error {

	mq.Consume(topic.MessageDeliverEvent.String(), HandleMessageDeliverEvent)
	mq.Consume(topic.ConnCmdCreate.String(), HandleConnCmdCreate)

	s := server.New(bifrost.KeyOpts{PriKey: pri, PubKey: pub})

	s.OnConnection(OnConnection)
	s.OnError(OnError)

	s.Listen(config.GrpcGateway.Ip)

	return nil
}

func OnConnection(connection bifrost.Connection) {

	connection.Handle(DefaultMux)
	ConnectionStore.AddConnection(connection)

	PublishNewConnEvent(connection)

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

	pri, pub, err = km.GenerateKey(convertToKeyGenOpts(conf.GetConfiguration().Authentication.KeyType))

	if err != nil {
		log.Fatal(err.Error())
	}

	return pri, pub
}

func convertToKeyGenOpts(keyType string) key.KeyGenOpts {

	switch keyType {
	case "RSA1024":
		return key.RSA1024
	case "RSA2048":
		return key.RSA2048
	case "RSA4096":
		return key.RSA4096
	default:
		return key.RSA1024
	}
}
