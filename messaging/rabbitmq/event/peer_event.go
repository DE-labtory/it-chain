package event

var LeaderInfoRequestProtocol string = "LeaderInfoRequestProtocol"
var LeaderInfoDeliverProtocol string = "LeaderInfoDeliverProtocol"

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
type LeaderInfoRequestCmd struct {
	TimeUnix int64
}

//todo define event
type LeaderInfoPublishEvent struct {
	LeaderId string
	Address  string
}
