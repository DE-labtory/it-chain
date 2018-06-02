package gateway

type GrpcService interface {
	Dial(address string) (Connection, error)
	CloseConnection(connID string)
	SendMessages(message []byte, protocol string, connIDs ...string)
}
