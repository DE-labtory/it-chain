package model

type Member struct {
	peerId PeerID
}

func (m Member) GetStringID() string {
	return m.peerId.ID
}
