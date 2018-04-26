package event

type MessageDeliveryEvent struct {
	Recipients []string
	Body       []byte
	Protocol   string
}

type NewConn struct {
	Id      string
	Address string
}
