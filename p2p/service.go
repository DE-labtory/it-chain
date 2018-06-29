package p2p

type PeerService interface{
	Dial(ipAddress string) error
}
