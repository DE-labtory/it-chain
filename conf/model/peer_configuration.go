package model

type PeerConfiguration struct {
	Empty string
}

func NewPeerConfiguration() PeerConfiguration {
	return PeerConfiguration{
		Empty: "empty",
	}
}
