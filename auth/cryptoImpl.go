package auth

import (
	"reflect"
	"errors"
	"crypto/elliptic"
)

type cryptoHelper struct {

	priKey			Key
	pubKey			Key

	keyManager		keyManager
	signers 		map[reflect.Type]signer
	verifiers 		map[reflect.Type]verifier
	keyGenerators 	map[reflect.Type]keyGenerator

}

func NewCrypto(path string) (Crypto, error) {

	km := &keyManager{}
	km.Init(path)

	signers := make(map[reflect.Type]signer)
	signers[reflect.TypeOf(&rsaPrivateKey{})] = &rsaSigner{}
	signers[reflect.TypeOf(&ecdsaPrivateKey{})] = &ecdsaSigner{}

	verifiers := make(map[reflect.Type]verifier)
	verifiers[reflect.TypeOf(&rsaPublicKey{})] = &rsaVerifier{}
	verifiers[reflect.TypeOf(&ecdsaPublicKey{})] = &ecdsaVerifier{}

	keyGenerators := make(map[reflect.Type]keyGenerator)
	keyGenerators[reflect.TypeOf(&RSAKeyGenOpts{})] = &rsaKeyGenerator{1024}
	keyGenerators[reflect.TypeOf(&ECDSAKeyGenOpts{})] = &ecdsaKeyGenerator{elliptic.P256()}

	ch := &cryptoHelper{
		keyManager:		*km,
		signers: 		signers,
		verifiers: 		verifiers,
		keyGenerators: 	keyGenerators,
	}

	return ch, nil

}

func (ch *cryptoHelper) Sign(digest []byte, opts SignerOpts) ([]byte, error) {

	var err error

	ch.priKey, ch.pubKey, err = ch.LoadKey()
	if err != nil {
		return nil, errors.New("Key is not exist.")
	}

	if len(digest) == 0 {
		return nil, errors.New("invalid digest.")
	}

	if ch.priKey == nil {
		return nil, errors.New("Private key is not exist.")
	}

	signer, found := ch.signers[reflect.TypeOf(ch.priKey)]
	if !found {
		return nil, errors.New("unsupported key type.")
	}

	signature, err := signer.Sign(ch.priKey, digest, opts)
	if err != nil {
		return nil, errors.New("signing error is occurred")
	}

	return signature, err

}

func (ch *cryptoHelper) Verify(key Key, signature, digest []byte, opts SignerOpts) (bool, error) {

	if key == nil {
		return false, errors.New("invalid key")
	}

	if len(signature) == 0 {
		return false, errors.New("invalid signature")
	}

	if len(digest) == 0 {
		return false, errors.New("invalid digest")
	}

	verifier, found := ch.verifiers[reflect.TypeOf(key)]
	if !found {
		return false, errors.New("unsupported key type")
	}

	valid, err := verifier.Verify(key, signature, digest, opts)
	if err != nil {
		return false, errors.New("verifying error is occurred")
	}

	return valid, nil

}

func (ch *cryptoHelper) GenerateKey(opts KeyGenOpts) (pri, pub Key, err error) {

	if opts == nil {
		return nil, nil, errors.New("Invalid KeyGen Options")
	}

	keyGenerator, found := ch.keyGenerators[reflect.TypeOf(opts)]
	if !found {
		return nil,nil, errors.New("Invalid KeyGen Options")
	}

	pri, pub, err = keyGenerator.GenerateKey(opts)
	if err != nil {
		return nil, nil, errors.New("Failed to generate a Key")
	}

	err = ch.keyManager.Store(pri, pub)
	if err != nil {
		return nil, nil, errors.New("Failed to store a Key")
	}

	return pri, pub, nil

}

func (ch *cryptoHelper) LoadKey() (pri, pub Key, err error) {

	pri, pub, err = ch.keyManager.Load()
	if err != nil {
		return nil, nil, err
	}

	return pri, pub, err

}