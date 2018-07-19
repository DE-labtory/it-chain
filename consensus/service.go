package consensus

type ParliamentService interface {
	Elect(parliament Parliament) ([]*Representative, error)
}
