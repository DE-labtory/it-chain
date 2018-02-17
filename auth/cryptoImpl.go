package auth

import (
	"reflect"
	"errors"
	"crypto/elliptic"
	"crypto/sha256"
)

type cryptoImpl struct {

	priKey			Key
	pubKey			Key

	keyManager		keyManager
	signers 		map[reflect.Type]signer
	verifiers 		map[reflect.Type]verifier
	keyGenerators 	map[reflect.Type]keyGenerator

}

func NewCryptoImpl(path string, keyGenOpts KeyGenOpts) (Crypto, error) {

	km := &keyManager{}
	km.Init(path)

	signers := make(map[reflect.Type]signer)
	signers[reflect.TypeOf(&rsaPrivateKey{})] = &rsaSigner{}
	signers[reflect.TypeOf(&ecdsaPrivateKey{})] = &ecdsaSigner{}

	verifiers := make(map[reflect.Type]verifier)
	verifiers[reflect.TypeOf(&rsaPublicKey{})] = &rsaVerifier{}
	verifiers[reflect.TypeOf(&ecdsaPublicKey{})] = &ecdsaVerifier{}

	keyGenerators := make(map[reflect.Type]keyGenerator)
	keyGenerators[reflect.TypeOf(&RSAKeyGenOpts{})] = &rsaKeyGenerator{2048}
	keyGenerators[reflect.TypeOf(&ECDSAKeyGenOpts{})] = &ecdsaKeyGenerator{elliptic.P521()}

	ci := &cryptoImpl{
		keyManager:		*km,
		signers: 		signers,
		verifiers: 		verifiers,
		keyGenerators: 	keyGenerators,
	}

	err := ci.loadKey()
	if err != nil ||
		ci.priKey.Algorithm() != keyGenOpts.Algorithm() ||
			ci.pubKey.Algorithm() != keyGenOpts.Algorithm() {
			err := ci.generateKey(keyGenOpts)
			if err != nil {
				return nil, err
			}
	}

	return ci, nil

}

func (ci *cryptoImpl) generateKey(opts KeyGenOpts) (err error) {

	// remove all exist key file in specific path
	err = ci.keyManager.removeKey()
	if err != nil {
		return err
	}

	if opts == nil {
		return errors.New("Invalid KeyGen Options")
	}

	keyGenerator, found := ci.keyGenerators[reflect.TypeOf(opts)]
	if !found {
		return errors.New("Invalid KeyGen Options")
	}

	pri, pub, err := keyGenerator.GenerateKey(opts)
	if err != nil {
		return errors.New("Failed to generate a Key")
	}

	err = ci.keyManager.Store(pri, pub)
	if err != nil {
		return errors.New("Failed to store a Key")
	}

	ci.priKey, ci.pubKey = pri, pub
	return nil

}

func (ci *cryptoImpl) Sign(data []byte, opts SignerOpts) ([]byte, error) {

	var err error

	err = ci.loadKey()
	if err != nil {
		return nil, errors.New("Key is not exist.")
	}

	if len(data) == 0 {
		return nil, errors.New("invalid data.")
	}

	if ci.priKey == nil {
		return nil, errors.New("Private key is not exist.")
	}

	signer, found := ci.signers[reflect.TypeOf(ci.priKey)]
	if !found {
		return nil, errors.New("unsupported key type.")
	}

	hash := sha256.New()
	hash.Write(data)
	digest := hash.Sum(nil)

	signature, err := signer.Sign(ci.priKey, digest, opts)
	if err != nil {
		return nil, errors.New("signing error is occurred")
	}

	return signature, err

}

func (ci *cryptoImpl) Verify(key Key, signature, digest []byte, opts SignerOpts) (bool, error) {

	if key == nil {
		return false, errors.New("invalid key")
	}

	if len(signature) == 0 {
		return false, errors.New("invalid signature")
	}

	if len(digest) == 0 {
		return false, errors.New("invalid digest")
	}

	verifier, found := ci.verifiers[reflect.TypeOf(key)]
	if !found {
		return false, errors.New("unsupported key type")
	}

	valid, err := verifier.Verify(key, signature, digest, opts)
	if err != nil {
		return false, errors.New("verifying error is occurred")
	}

	return valid, nil

}

// load private and public key data to cryptoHelper if they are exist in path
func (ci *cryptoImpl) loadKey() (err error) {

	pri, pub, err := ci.keyManager.Load()
	if err != nil {
		return err
	}

	ci.priKey, ci.pubKey = pri, pub

	return nil

}

// return private key and public key if they are exist in path
func (ci *cryptoImpl) GetKey() (pri, pub Key, err error) {

	if ci.priKey == nil || ci.pubKey == nil {
		err := ci.loadKey()
		if err != nil {
			return nil, nil, err
		}
	}

	return ci.priKey, ci.pubKey, nil

}