package infra

import (
	"github.com/it-chain/bifrost"
	"github.com/it-chain/bifrost/client"
	"github.com/it-chain/heimdall/key"
	"github.com/it-chain/it-chain-Engine/gateway"
)

type GrpcDialService struct {
	store  *bifrost.ConnectionStore
	priKey key.PriKey
	pubKey key.PubKey
}

func (g GrpcDialService) Dial(address string) (gateway.Connection, error) {

	connection, err := client.Dial(g.buildDialOption(address))

	if err != nil {
		return gateway.Connection{}, err
	}

	g.store.GetConnection(connection.GetID())
}

func (g GrpcDialService) buildDialOption(address string) (string, client.ClientOpts, client.GrpcOpts) {

	clientOpt := client.ClientOpts{
		Ip:     address,
		PriKey: g.priKey,
		PubKey: g.pubKey,
	}

	grpcOpt := client.GrpcOpts{
		TlsEnabled: false,
		Creds:      nil,
	}

	return address, clientOpt, grpcOpt
}
