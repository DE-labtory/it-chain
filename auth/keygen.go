package auth

import (
	"crypto/rsa"
	"crypto/rand"
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"
	"errors"
)

type RSAKeyGenerator struct {
	bits int
}

func (keygen *RSAKeyGenerator) KeyGenerate(opts KeyGenOpts) (key Key, err error) {

	if keygen.bits <= 0 {
		return nil, errors.New("Bits length should be bigger than 0")
	}

	generatedKey, err := rsa.GenerateKey(rand.Reader, keygen.bits)

	if err != nil {
		return nil, fmt.Errorf("Failed to generate RSA key : %s", err)
	}

	return &rsaPrivateKey{generatedKey}, nil

}

type ECDSAKeyGenerator struct {
	curve elliptic.Curve
}

func (keygen *ECDSAKeyGenerator) KeyGenerate(opts KeyGenOpts) (key Key, err error) {

	if keygen.curve == nil {
		return nil, errors.New("Curve value have not to be nil")
	}

	generatedKey, err := ecdsa.GenerateKey(keygen.curve, rand.Reader)

	if err != nil {
		return nil, fmt.Errorf("Failed to generate ECDSA key : %s", err)
	}

	return &ecdsaPrivateKey{generatedKey}, nil

}
