package event

//todo define event
type LeaderChangeEvent struct {
	PeerId string
}

//todo define event
type PeerDisconnectEvent struct {
	PeerId string
}

//todo define event
type PeerConnectEvent struct {
	PeerId string
}

//todo define event
type LeaderInfoDeliverCmd struct {
	timeUnix int64
}

//todo define event
type LeaderInfoReceiveEvent struct {
	PeerId  string
	Address string
}
