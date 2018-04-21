package model

type Leader struct {
	peerId PeerID
}

func (l *Leader) StringPeerId() string {
	return l.peerId.ToString()
}
