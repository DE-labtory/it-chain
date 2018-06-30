package p2p

type GrpcCommandService interface {
	DeliverLeaderInfo(connectionId string, leader Leader) error
	DeliverPeerLeaderTable(connectionId string, peerLeaderTable PeerLeaderTable)error
	DeliverRequestVoteMessages(connectionIds ...string) error
	DeliverVoteLeaderMessage(connectionId string) error
	DeliverUpdateLeaderMessage(connecionId string, leader Leader) error
}
