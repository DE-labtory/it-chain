package p2p

type PeerService interface {
	Dial(ipAddress string) error
	DeliverLeaderInfo(connectionId string, leader Leader) error
	DeliverPeerLeaderTable(connectionId string, peerLeaderTable PeerLeaderTable) error
}

type RaftService interface {
	DeliverRequestVoteMessages(connectionIds ...string) error
	DeliverVoteLeaderMessage(connectionId string) error
	DeliverUpdateLeaderMessage(connecionId string, leader Leader) error
}
