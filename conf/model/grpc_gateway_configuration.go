package model

type GrpcGatewayConfiguration struct {
	Empty string
}

func NewGrpcGatewayConfiguration() GrpcGatewayConfiguration {
	return GrpcGatewayConfiguration{
		Empty: "empty",
	}
}
