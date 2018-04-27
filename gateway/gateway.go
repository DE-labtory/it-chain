package main

import (
	"github.com/it-chain/bifrost"
	"github.com/it-chain/it-chain-Engine/messaging"
	"github.com/it-chain/it-chain-Engine/messaging/topic"
	"github.com/it-chain/it-chain-Engine/conf"
	"github.com/it-chain/it-chain-Engine/conf/model/common"
)


func Start() error {

	config := conf.GetConfiguration()

	messaging := messaging.NewMessaging(config.Common.Messaging.Url)
	ec        EventConsumer
	host      bifrost.BifrostHost

	if err := gw.messaging.Consume(topic.MessageDeliverEvent.String(), gw.ec.HandleMessageDeliverEvent); err != nil {
		return err
	}

	return nil
}
