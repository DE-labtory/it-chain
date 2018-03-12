package parliament

type Member struct{
	ID PeerID
}

func (m Member) GetStringID() string{
	return m.ID.ID
}
