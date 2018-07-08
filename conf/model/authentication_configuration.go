package model

type AuthenticationConfiguration struct {
	KeyType string
	KeyPath string
}

func NewAuthenticationConfiguration() AuthenticationConfiguration {
	return AuthenticationConfiguration{
		KeyType: "ECDSA256",
		KeyPath: ".it-chain/",
	}
}
