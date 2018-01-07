package auth

import "crypto/ecdsa"

type ecdsaSigner struct{}

func (s *ecdsaSigner) Sign(key Key, digest []byte, opts SignerOpt) (signature []byte, err error) {

}

type ecdsaVerifier struct{}

func (v *ecdsaVerifier) Verify(key Key, signature, digest []byte, opts SignerOpt) (valid bool, err error) {

}

type ecdsaPublicKey struct {
	pub *ecdsa.PublicKey
}

type ecdsaPrivateKey struct {
	priv *ecdsa.PrivateKey
}