package p2p


//gossip 최상위 service
type GossipService interface{
	Gossip(gossipTable GossipTable, peersIP []string)
	Listen(GossipTable)
	Pull()
	Stop()
	PeersInfoOfChannel(channel string) []*PeerInfo
}