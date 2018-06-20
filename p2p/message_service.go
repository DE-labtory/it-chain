package p2p

type MessageService interface {
	RequestLeaderInfo(nodeId NodeId) error
	DeliverLeaderInfo(nodeId NodeId, leader Leader) error
	RequestNodeList(nodeId NodeId) error
	DeliverNodeList(nodeId NodeId, nodes []Node) error
	DeliverNode(nodeId NodeId, node Node) error
}
