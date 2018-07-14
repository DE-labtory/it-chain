package txpool

type LeaderId struct {
	Id string
}

func (lid LeaderId) ToString() string {
	return string(lid.Id)
}

//Aggregate root must implement aggregate interface
type Leader struct {
	LeaderId LeaderId
}
