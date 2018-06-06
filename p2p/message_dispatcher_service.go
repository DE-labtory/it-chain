package p2p

type MessageDispatcher interface {
	RequestLeaderInfo(peer Node) error
	DeliverLeaderInfo(toPeer Node, leader Node) error
	RequestTable(toNode Node) error
	ResponseTable(toNode Node, nodes []Node) error
	SendLeaderUpdateMessage(leader Node) error
	SendDeliverNodeListMessage(toNode Node) error // update node repository in specific node
}
