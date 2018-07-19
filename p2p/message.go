package p2p

type LeaderInfoRequestMessage struct {
	TimeUnix int64
}

type UpdateLeaderMessage struct {
	Peer Peer
}

type PLTableMessage struct {
	PLTable PLTable
}

type RequestVoteMessage struct {
}

type VoteMessage struct {
}
