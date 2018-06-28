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
	peers []Peer
}

type LeaderInfoResponseMessage struct {
	LeaderId string
	Address  string
}
