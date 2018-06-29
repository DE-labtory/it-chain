package p2p

type GrpcCommandService interface {
	DeliverLeaderInfo(nodeId PeerId, leader Leader) error
	DeliverPeerLeaderTable(connectionId string, peerLeaderTable PeerLeaderTable)error
}
