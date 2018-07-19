package model

// it-chain의 공통적이고 중요한 설정을 담는 구조체이다.
type EngineConfiguration struct {
	KeyPath              string
	Mode                 string
	Amqp                 string
	BootstrapNodeAddress string
}

func NewEngineConfiguration() EngineConfiguration {
	return EngineConfiguration{
		KeyPath:              ".it-chain/",
		Mode:                 "solo",
		BootstrapNodeAddress: "127.0.0.1:5555",
		Amqp:                 "amqp://guest:guest@localhost:5672/",
	}
}
