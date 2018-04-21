package event

type MessageDeliveryEvent struct {
	Recipients []string
	Body       []byte
	Protocol   string
}
