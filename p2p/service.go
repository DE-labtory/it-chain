package p2p

type PeerService interface {
	Dial(ipAddress string) error
	DeliverLeaderInfo(connectionId string, leader Leader) error
	DeliverPeerLeaderTable(connectionId string, peerLeaderTable PeerLeaderTable) error
}

type VotingService interface {

	CountDownLeftTimeBy(tick int64)
	SetState(state string)
	GetState() string
	ResetLeftTime() int64
	GetLeftTime() int64
	GetVoteCount() int
	ResetVoteCount()
	CountUp()
	DeliverRequestVoteMessages(connectionIds ...string) error
	DeliverVoteLeaderMessage(connectionId string) error
	DeliverUpdateLeaderMessage(connecionId string, leader Leader) error
}