package auth

import (
	"reflect"
	"errors"
)

type impl struct {

	signers map[reflect.Type]Signer
	verifiers map[reflect.Type]Verifier

}

func New() (Crypto, error) {

	signers := make(map[reflect.Type]Signer)
	signers[reflect.TypeOf(&rsaPrivateKey{})] = &rsaSigner{}
	signers[reflect.TypeOf(&ecdsaPrivateKey{})] = &ecdsaSigner{}

	verifiers := make(map[reflect.Type]Verifier)
	verifiers[reflect.TypeOf(&rsaPublicKey{})] = &rsaVerifier{}
	verifiers[reflect.TypeOf(&ecdsaPublicKey{})] = &ecdsaVerifier{}

	impl := &impl{
		signers: 	signers,
		verifiers: 	verifiers,
	}

	return impl, nil

}

func (i *impl) Sign(key Key, digest []byte, opts SignerOpt) (signature []byte, err error) {

	if key == nil {
		return nil, errors.New("invalid key")
	}

	if len(digest) == 0 {
		return nil, errors.New("invalid digest")
	}

	signer, found := i.signers[reflect.TypeOf(key)]
	if !found {
		return nil, errors.New("unsupported key type")
	}

	signature, err = signer.Sign(key, digest, opts)
	if err != nil {
		return nil, errors.New("signing error is occurred")
	}

	return

}

func (i *impl) Verify(key Key, signature, digest []byte, opts SignerOpt) (valid bool, err error) {

	if key == nil {
		return false, errors.New("invalid key")
	}

	if len(signature) == 0 {
		return false, errors.New("invalid signature")
	}

	if len(digest) == 0 {
		return false, errors.New("invalid digest")
	}

	verifier, found := i.verifiers[reflect.TypeOf(key)]
	if !found {
		return false, errors.New("unsupported key type")
	}

	valid, err = verifier.Verify(key, signature, digest, opts)
	if err != nil {
		return false, errors.New("verifying error is occurred")
	}

	return

}