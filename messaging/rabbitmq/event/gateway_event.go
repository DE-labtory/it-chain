package event

// 메세지 전달하는 이벤트 구조이다.
// 복수의 수신자와, 내용, 프로토콜로써 정의된다.
type MessageDeliverEvent struct {
	Recipients []string
	Body       []byte
	Protocol   string
}

type MessageReceiveEvent struct {
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
