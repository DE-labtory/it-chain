package blockchain

type NodeId struct {
	Id string
}

type Node struct {
	IpAddress string
	NodeId NodeId
}
