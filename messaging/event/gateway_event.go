package event

type MessageDeliveryEvent struct {
	Recipients []string
	Body       []byte
	Protocol   string
}

type ConnectionCreated struct {
	Id      string
	Address string
}
