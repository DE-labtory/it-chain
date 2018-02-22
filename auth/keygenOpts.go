package auth

type RSAKeyGenOpts struct {}

type ECDSAKeyGenOpts struct {}

func (opts *RSAKeyGenOpts) Algorithm() string {
	return RSA
}

func (opts *ECDSAKeyGenOpts) Algorithm() string {
	return ECDSA
}
