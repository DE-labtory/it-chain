package event

type MessageDeliverEvent struct {
	Recipients []string
	Body       []byte
	Protocol   string
}

type NewConnEvent struct {
	Id      string
	Address string
}

type ConnCmdCreate struct {
	Address string
}
