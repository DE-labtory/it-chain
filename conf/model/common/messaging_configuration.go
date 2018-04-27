package common

type MessagingConfiguration struct {
	Url string
}

func NewMessagingConfiguration() MessagingConfiguration {
	return MessagingConfiguration{
		Url: "amqp://guest:guest@localhost:5672/",
	}
}
