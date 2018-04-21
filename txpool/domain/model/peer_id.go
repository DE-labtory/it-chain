package model

type PeerID struct {
	id string
}

func (p PeerID) ToString() string {
	return string(p.id)
}
