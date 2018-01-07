package auth

import "crypto/ecdsa"

type ecdsaSigner struct{}

func (s *ecdsaSigner) Sign(key Key, digest []byte, opts SignerOpt) (signature []byte, err error) {
	return nil, nil
}

type ecdsaVerifier struct{}

func (v *ecdsaVerifier) Verify(key Key, signature, digest []byte, opts SignerOpt) (valid bool, err error) {
	return true, nil
}

type ecdsaPublicKey struct {
	pub *ecdsa.PublicKey
}

type ecdsaPrivateKey struct {
	priv *ecdsa.PrivateKey
}