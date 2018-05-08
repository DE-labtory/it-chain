package main

import (
	"github.com/it-chain/heimdall/auth"
	"github.com/it-chain/it-chain-Engine/conf"
	"github.com/it-chain/it-chain-Engine/messaging"
	"github.com/it-chain/heimdall/key"
	"net/http"
	"log"
)

func Start() error {

	config := conf.GetConfiguration()

	loadKeyPair(config.Authentication.KeyPath)



	key.NewKeyManager("~/asd")

	http.Handle()

	signer := NewSigner()

	mc :=

	mq := messaging.NewRabbitmq(config.Common.Messaging.Url)
	return nil
}

func loadKeyPair(keyPath string) (key.PriKey, key.PubKey){

	km, err := key.NewKeyManager(keyPath)

	if err != nil {
		log.Fatal(err.Error())
	}

	pri, pub, err := km.GetKey()

	if err != nil{
		km.GenerateKey(conf.GetConfiguration().Authentication.KeyType)
	}

	return pri,pub
}
