package model

type GrpcGatewayConfiguration struct {
	Address string
	Port    string
}

func NewGrpcGatewayConfiguration() GrpcGatewayConfiguration {
	return GrpcGatewayConfiguration{
		Address: "127.0.0.1",
		Port:    "13579",
	}
}
