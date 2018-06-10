package p2p

type MessageDispatcher interface {
	RequestLeaderInfo(nodeId NodeId) error
	DeliverLeaderInfo(nodeId NodeId, leader Node) error
	RequestNodeList(nodeId NodeId) error
	ResponseTable(toNode Node, nodes []Node) error
	SendLeaderUpdateMessage(nodeId NodeId, leader Node) error
	SendDeliverNodeListMessage(nodeId NodeId, nodeList []Node) error
}
