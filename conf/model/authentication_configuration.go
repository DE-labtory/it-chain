package model

type AuthenticationConfiguration struct {
	KeyType int
	KeyPath string
}

func NewAuthenticationConfiguration() AuthenticationConfiguration {
	return AuthenticationConfiguration{
		KeyType: 0,
		KeyPath: ".it-chain/",
	}
}
