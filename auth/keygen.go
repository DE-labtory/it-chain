package auth

import (
	"crypto/rsa"
	"crypto/rand"
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"
	"errors"
)

type rsaKeyGenerator struct {
	bits int
}

func (keygen *rsaKeyGenerator) GenerateKey(opts KeyGenOpts) (pri, pub Key, err error) {

	if keygen.bits <= 0 {
		return nil, nil, errors.New("Bits length should be bigger than 0")
	}

	generatedKey, err := rsa.GenerateKey(rand.Reader, keygen.bits)

	if err != nil {
		return nil, nil, fmt.Errorf("Failed to generate RSA key : %s", err)
	}

	pri = &rsaPrivateKey{generatedKey}
	pub, err = pri.(*rsaPrivateKey).PublicKey()
	if err != nil {
		return nil, nil, err
	}

	return

}

type ecdsaKeyGenerator struct {
	curve elliptic.Curve
}

func (keygen *ecdsaKeyGenerator) GenerateKey(opts KeyGenOpts) (pri, pub Key, err error) {

	if keygen.curve == nil {
		return nil, nil, errors.New("Curve value have not to be nil")
	}

	generatedKey, err := ecdsa.GenerateKey(keygen.curve, rand.Reader)

	if err != nil {
		return nil, nil, fmt.Errorf("Failed to generate ECDSA key : %s", err)
	}

	pri = &ecdsaPrivateKey{generatedKey}
	pub, err = pri.(*ecdsaPrivateKey).PublicKey()
	if err != nil {
		return nil, nil, err
	}

	return

}
