package model

type GrpcGatewayConfiguration struct {
	Empty string
	Ip    string
}

func NewGrpcGatewayConfiguration() GrpcGatewayConfiguration {
	return GrpcGatewayConfiguration{
		Empty: "empty",
		Ip:    "127.0.0.1:13579",
	}
}
