package auth

import (
	"crypto/rsa"
	"crypto/rand"
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"
)

type rsaKeyGenOpts struct {}

type rsaKeyGenerator struct {
	bits int
}

func (keygen *rsaKeyGenerator) KeyGenerate(opts KeyGenOpts) (key Key, err error) {

	generatedKey, err := rsa.GenerateKey(rand.Reader, keygen.bits)

	if err != nil {
		return nil, fmt.Errorf("Failed to generate RSA key : %s", err)
	}

	return &rsaPrivateKey{generatedKey}, nil

}

type ecdsaKeyGenOpts struct {}

type ecdsaKeyGenerator struct {
	curve elliptic.Curve
}

func (keygen *ecdsaKeyGenerator) KeyGenerate(opts KeyGenOpts) (key Key, err error) {

	generatedKey, err := ecdsa.GenerateKey(keygen.curve, rand.Reader)

	if err != nil {
		return nil, fmt.Errorf("Failed to generate ECDSA key : %s", err)
	}

	return &ecdsaPrivateKey{generatedKey}, nil

}
