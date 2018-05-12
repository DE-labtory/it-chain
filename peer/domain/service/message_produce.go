package service

type Publish func(topic string, data []byte) error

type MessageProducer interface {
	RequestLeaderInfo(ipAddress string) error
}
