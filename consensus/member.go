package consensus

type MemberId struct {
	Id string
}

type Member struct {
	MemberId MemberId
}

func (mid MemberId) ToString() string {
	return string(mid.Id)
}

func (m *Member) StringMemberId() string {
	return m.MemberId.ToString()
}

func (m Member) GetId() string {
	return m.StringMemberId()
}
