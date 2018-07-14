package icode

type ContainerService interface {
	StartContainer(meta Meta) error
	StopContainer(id ID) error
	ExecuteTransaction(tx Transaction) (*Result, error)
}
