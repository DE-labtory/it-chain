package icode

type ContainerService interface {
	Start(meta Meta) error
	Stop(id ID) error
	Run(tx Transaction) (*Result, error)
}
