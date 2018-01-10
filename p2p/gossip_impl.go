package p2p

type GossipServiceImpl struct{
	gossipTable GossipTable
}

//gossip message를 peersIP에 전파
func (gs GossipServiceImpl) Gossip(gossipTable GossipTable, peersIP []string){

}

func (gs GossipServiceImpl) Listen(gossipTable GossipTable){

}

func (gs GossipServiceImpl) Pull(){

}

func (gs GossipServiceImpl) Stop(){

}

func (gs GossipServiceImpl) PeersInfoOfChannel(){

}

