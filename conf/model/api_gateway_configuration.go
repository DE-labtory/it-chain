package model

type ApiGatewayConfiguration struct {
	Address string
	Port    string
}

func NewApiGatewayConfiguration() ApiGatewayConfiguration {
	return ApiGatewayConfiguration{
		Address: "127.0.0.1",
		Port:    "4444",
	}
}
