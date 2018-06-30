package consensus

type LeaderId struct {
	Id string
}

type Leader struct {
	LeaderId LeaderId
}

func (lid LeaderId) ToString() string {
	return string(lid.Id)
}

func (l *Leader) StringLeaderId() string {
	return l.LeaderId.ToString()
}

func (l Leader) GetId() string {
	return l.StringLeaderId()
}
