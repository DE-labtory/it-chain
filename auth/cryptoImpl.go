package auth

import (
	"reflect"
	"errors"
	"crypto/elliptic"
	"crypto/sha256"
)

type cryptoHelper struct {

	priKey			Key
	pubKey			Key

	keyManager		keyManager
	signers 		map[reflect.Type]signer
	verifiers 		map[reflect.Type]verifier
	keyGenerators 	map[reflect.Type]keyGenerator

}

func NewCrypto(path string, keyGenOpts KeyGenOpts) (Crypto, error) {

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

	ch := &cryptoHelper{
		keyManager:		*km,
		signers: 		signers,
		verifiers: 		verifiers,
		keyGenerators: 	keyGenerators,
	}

	err := ch.loadKey()
	if err != nil ||
		ch.priKey.Algorithm() != keyGenOpts.Algorithm() ||
			ch.pubKey.Algorithm() != keyGenOpts.Algorithm() {
			err := ch.generateKey(keyGenOpts)
			if err != nil {
				return nil, err
			}
	}

	return ch, nil

}

func (ch *cryptoHelper) generateKey(opts KeyGenOpts) (err error) {

	// remove all exist key file in specific path
	err = ch.keyManager.removeKey()
	if err != nil {
		return err
	}

	if opts == nil {
		return errors.New("Invalid KeyGen Options")
	}

	keyGenerator, found := ch.keyGenerators[reflect.TypeOf(opts)]
	if !found {
		return errors.New("Invalid KeyGen Options")
	}

	pri, pub, err := keyGenerator.GenerateKey(opts)
	if err != nil {
		return errors.New("Failed to generate a Key")
	}

	err = ch.keyManager.Store(pri, pub)
	if err != nil {
		return errors.New("Failed to store a Key")
	}

	ch.priKey, ch.pubKey = pri, pub
	return nil

}

func (ch *cryptoHelper) Sign(data []byte, opts SignerOpts) ([]byte, error) {

	var err error

	err = ch.loadKey()
	if err != nil {
		return nil, errors.New("Key is not exist.")
	}

	if len(data) == 0 {
		return nil, errors.New("invalid data.")
	}

	if ch.priKey == nil {
		return nil, errors.New("Private key is not exist.")
	}

	signer, found := ch.signers[reflect.TypeOf(ch.priKey)]
	if !found {
		return nil, errors.New("unsupported key type.")
	}

	hash := sha256.New()
	hash.Write(data)
	digest := hash.Sum(nil)

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

// load private and public key data to cryptoHelper if they are exist in path
func (ch *cryptoHelper) loadKey() (err error) {

	pri, pub, err := ch.keyManager.Load()
	if err != nil {
		return err
	}

	ch.priKey, ch.pubKey = pri, pub

	return nil

}

// return private key and public key if they are exist in path
func (ch *cryptoHelper) GetKey() (pri, pub Key, err error) {

	if ch.priKey == nil || ch.pubKey == nil {
		err := ch.loadKey()
		if err != nil {
			return nil, nil, err
		}
	}

	return ch.priKey, ch.pubKey, nil

}