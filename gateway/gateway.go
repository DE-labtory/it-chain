package main

import (
	"github.com/it-chain/heimdall/auth"
	"github.com/it-chain/it-chain-Engine/conf"
	"github.com/it-chain/it-chain-Engine/messaging"
	"github.com/it-chain/heimdall/key"
	"net/http"
)

func Start() error {

	config := conf.GetConfiguration()

	key.NewKeyManager("~/asd")

	http.Handle()

	signer := NewSigner()

	mc :=

	mq := messaging.NewRabbitmq(config.Common.Messaging.Url)
	return nil
}
