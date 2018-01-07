package auth

import "crypto/rsa"

type rsaSigner struct{}

func (s *rsaSigner) Sign(key Key, digest []byte, opts SignerOpt) (signature []byte, err error) {

}

type rsaVerifier struct{}

func (v *rsaVerifier) Verify(key Key, signature, digest []byte, opts SignerOpt) (valid bool, err error) {
}

type rsaPublicKey struct {
	pub *rsa.PublicKey
}

type rsaPrivateKey struct {
	priv *rsa.PrivateKey
}