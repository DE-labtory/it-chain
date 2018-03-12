package parliament

type Leader struct{
	ID PeerID
}

func (l Leader) GetStringID() string{
	return l.ID.ID
}