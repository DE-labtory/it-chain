package blockchain

type CommandService interface {
	SendBlockValidateCommand(block Block) error
}
