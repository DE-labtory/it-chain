package icode

type CommandService interface {
	SendBlockExecuteResultCommand(results []Result, blockId string) error
}
