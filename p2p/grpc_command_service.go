package p2p

type GrpcCommandService interface {
	RequestLeaderInfo(nodeId PeerId) error
	DeliverLeaderInfo(nodeId PeerId, leader Leader) error
	RequestPeerList(nodeId PeerId) error
	DeliverPeerList(nodeId PeerId, nodes []Peer) error
	DeliverPeer(nodeId PeerId, node Peer) error
}
