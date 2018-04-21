package model

type Leader struct {
	peerId PeerID
}

func (l Leader) GetStringID() string {
	return l.peerId.ID
}
