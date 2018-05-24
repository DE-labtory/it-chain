package txpool

type Leader struct {
	leaderId LeaderId
}

func (l *Leader) StringLeaderId() string {
	return l.leaderId.ToString()
}

type LeaderId struct {
	id string
}

func (lid LeaderId) ToString() string {
	return string(lid.id)
}
