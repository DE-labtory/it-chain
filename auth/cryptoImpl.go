package auth

import (
	"reflect"
	"errors"
	"crypto/elliptic"
)

type collector struct {

	signers map[reflect.Type]Signer
	verifiers map[reflect.Type]Verifier
	keyGenerators map[reflect.Type]KeyGenerator

}

func NewCollector() (Crypto, error) {

	signers := make(map[reflect.Type]Signer)
	signers[reflect.TypeOf(&rsaPrivateKey{})] = &rsaSigner{}
	signers[reflect.TypeOf(&ecdsaPrivateKey{})] = &ecdsaSigner{}

	verifiers := make(map[reflect.Type]Verifier)
	verifiers[reflect.TypeOf(&rsaPublicKey{})] = &rsaVerifier{}
	verifiers[reflect.TypeOf(&ecdsaPublicKey{})] = &ecdsaVerifier{}

	keyGenerators := make(map[reflect.Type]KeyGenerator)
	keyGenerators[reflect.TypeOf(&RSAKeyGenOpts{false})] = &RSAKeyGenerator{1024}
	keyGenerators[reflect.TypeOf(&ECDSAKeyGenOpts{false})] = &ECDSAKeyGenerator{elliptic.P256()}

	coll := &collector{
		signers: 		signers,
		verifiers: 		verifiers,
		keyGenerators: 	keyGenerators,
	}

	return coll, nil

}

func (c *collector) Sign(key Key, digest []byte, opts SignerOpts) (signature []byte, err error) {

	if key == nil {
		return nil, errors.New("invalid key")
	}

	if len(digest) == 0 {
		return nil, errors.New("invalid digest")
	}

	signer, found := c.signers[reflect.TypeOf(key)]
	if !found {
		return nil, errors.New("unsupported key type")
	}

	signature, err = signer.Sign(key, digest, opts)
	if err != nil {
		return nil, errors.New("signing error is occurred")
	}

	return

}

func (c *collector) Verify(key Key, signature, digest []byte, opts SignerOpts) (valid bool, err error) {

	if key == nil {
		return false, errors.New("invalid key")
	}

	if len(signature) == 0 {
		return false, errors.New("invalid signature")
	}

	if len(digest) == 0 {
		return false, errors.New("invalid digest")
	}

	verifier, found := c.verifiers[reflect.TypeOf(key)]
	if !found {
		return false, errors.New("unsupported key type")
	}

	valid, err = verifier.Verify(key, signature, digest, opts)
	if err != nil {
		return false, errors.New("verifying error is occurred")
	}

	return

}

func (c *collector) KeyGenerate(opts KeyGenOpts) (key Key, err error) {

	if opts == nil {
		return nil, errors.New("Invalid KeyGen Options")
	}

	keyGenerator, found := c.keyGenerators[reflect.TypeOf(opts)]
	if !found {
		return nil, errors.New("Invalid KeyGen Options")
	}

	key, err = keyGenerator.KeyGenerate(opts)
	if err != nil {
		return nil, errors.New("Failed to generate a Key")
	}

	return

}