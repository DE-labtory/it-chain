package p2p

// GPPC 메세지를 위한 message.go
// topic 이름은 구조체이름을 이용한다.
type PeerListRequestMessage struct {
	TimeUnix int64
}

type LeaderInfoRequestMessage struct {
	TimeUnix int64
}

type PeerListResponseMessage struct {
	Peers []Peer
}

type LeaderInfoResponseMessage struct {
	LeaderId string
	Address  string
}

type UpdateLeaderMessage struct {
	Peer Peer
}

type PLTableMessage struct {
	PLTable PLTable
}

type RequestVoteMessage struct{}

type VoteMessage struct{}
