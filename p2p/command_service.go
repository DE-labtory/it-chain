package p2p

type CommandService interface {
	Dial(ipAddress string) error
}