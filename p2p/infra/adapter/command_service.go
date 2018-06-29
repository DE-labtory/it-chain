package adapter

import "github.com/it-chain/it-chain-Engine/gateway"

type CommandService struct{
	publish Publish
}

func (cs *CommandService) Dial(ipAddress string) error{
	command := gateway.ConnectionCreateCommand{
		Address: ipAddress,
	}
	cs.publish("Command", "connection.create", command)
}
