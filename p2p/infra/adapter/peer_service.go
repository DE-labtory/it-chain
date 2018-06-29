package adapter

import "github.com/it-chain/it-chain-Engine/gateway"

type PeerService struct {
	publish Publish
}

func (ps *PeerService) Dial(ipAddress string) error{
	command := gateway.ConnectionCreateCommand{
		Address: ipAddress,
	}
	ps.publish("Command", "connection.create", command)
	return nil
}

