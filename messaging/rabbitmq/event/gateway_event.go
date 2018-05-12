package event

type MessageDeliverEvent struct {
	Recipients []string
	Body       []byte
	Protocol   string
}

type MessageReceiveEvent struct {
	SenderID string
	Body     []byte
	Protocol string
}

type ConnCreateEvent struct {
	Id      string
	Address string
}

//todo @junbeomlee grpc-gateway가 이벤트 발행해야함
type ConnTerminateEvent struct {
	Id      string
	Address string
}

type ConnCreateCmd struct {
	Address string
}
