package consensus

type RepresentativeId string

type Representative struct{
	Id RepresentativeId
}

func (r Representative) GetIdString() string{
	return string(r.Id)
}

func NewRepresentative(Id string) *Representative{

	return &Representative{Id:RepresentativeId(Id)}
}