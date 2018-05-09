//Subscribe event and do corresponding logic

package main

import (
	"encoding/json"
	"log"

	"github.com/it-chain/bifrost"
	"github.com/it-chain/bifrost/client"
	"github.com/it-chain/it-chain-Engine/messaging/event"
	"github.com/streadway/amqp"
)

func HandleMessageDeliverEvent(amqpMessage amqp.Delivery) {

	MessageDelivery := &event.MessageDeliverEvent{}
	if err := json.Unmarshal(amqpMessage.Body, MessageDelivery); err != nil {
		// fail to unmarshal event
		return
	}

	deliver(MessageDelivery.Recipients, MessageDelivery.Protocol, MessageDelivery.Body)
}

func HandleConnCmdCreate(amqpMessage amqp.Delivery) {

	ConnCmdCreate := &event.ConnCmdCreate{}
	if err := json.Unmarshal(amqpMessage.Body, ConnCmdCreate); err != nil {
		// fail to unmarshal event
		return
	}

	clientOpt := client.ClientOpts{
		Ip:     ConnCmdCreate.Address,
		PriKey: pri,
		PubKey: pub,
	}

	grpcOpt := client.GrpcOpts{
		TlsEnabled: false,
		Creds:      nil,
	}

	conn, err := client.Dial(ConnCmdCreate.Address, clientOpt, grpcOpt)

	if err != nil {
		log.Println(err.Error())
		return
	}

	OnConnection(conn)
}

func deliver(recipients []string, protocol string, data []byte) error {

	for _, recipient := range recipients {
		connection := ConnectionStore.GetConnection(bifrost.ConnID(recipient))

		if connection != nil {
			connection.Send(data, protocol, nil, nil)
		}
	}

	return nil
}
