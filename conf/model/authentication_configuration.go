package model

type AuthenticationConfiguration struct {
	KeyType string
	KeyPath string
}

func NewAuthenticationConfiguration() AuthenticationConfiguration {
	return AuthenticationConfiguration{
		KeyType: "RSA",
		KeyPath: ".it-chain/",
	}
}
