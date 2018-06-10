package p2p

type MessageDispatcher interface {
	RequestLeaderInfo(peer Node) error
	DeliverLeaderInfo(toPeer Node, leader Leader) error
	RequestTable(toNode Node) error
	ResponseTable(toNode Node, nodes []Node) error
	DeliverNodeList(toNode Node, nodes []*Node) error
}
