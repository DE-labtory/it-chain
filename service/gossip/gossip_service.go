package gossip

import "it-chain/common"

var logger = common.GetLogger("p2p package")

//gossip 최상위 service
type GossipService interface{

	//gossip 전파
	Push(gossipTable GossipTable, peersIP []string)

	//gossip의
	Listen(GossipTable)

	Pull() GossipTable

	Stop()

	Update()

	GetMyGossipTable() GossipTable
}