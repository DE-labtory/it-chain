package common

// it-chain의 공통적인 설정을 담는 구조체이다.
type CommonConfiguration struct {
	BootNodeIp string
	NodeIp     string
	Messaging  MessagingConfiguration
}

func NewCommonConfiguration() CommonConfiguration {
	return CommonConfiguration{
		BootNodeIp: "127.0.0.1:4444",
		NodeIp:     "127.0.0.1:4444",
		Messaging:  NewMessagingConfiguration(),
	}
}
