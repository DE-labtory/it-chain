package blockchain

type NodeId struct {
	Id string
}

type Node struct {
	IpAddress string
	NodeId NodeId
}

type NodeRepository interface {
	AddNode(node Node) error
	DeleteNode(node Node) error
}

