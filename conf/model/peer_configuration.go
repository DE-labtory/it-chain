package model

type PeerConfiguration struct {
	LeaderElection string
}

func NewPeerConfiguration() PeerConfiguration {
	return PeerConfiguration{
		LeaderElection: "RAFT",
	}
}
