package auth

import (
	"reflect"
	"errors"
	"crypto/elliptic"
)

type cryptoHelper struct {

	privKey			Key
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

func (ch *cryptoHelper) Sign(digest []byte, opts SignerOpts) (signature []byte, err error) {

	if len(digest) == 0 {
		return nil, errors.New("invalid digest")
	}

	signer, found := ch.signers[reflect.TypeOf(key)]
	if !found {
		return nil, errors.New("unsupported key type")
	}

	signature, err = signer.Sign(key, digest, opts)
	if err != nil {
		return nil, errors.New("signing error is occurred")
	}

	return

}

func (c *cryptoHelper) Verify(key Key, signature, digest []byte, opts SignerOpts) (valid bool, err error) {

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

func (ch *cryptoHelper) GenerateKey(opts KeyGenOpts) (key Key, err error) {

	if opts == nil {
		return nil, errors.New("Invalid KeyGen Options")
	}

	keyGenerator, found := ch.keyGenerators[reflect.TypeOf(opts)]
	if !found {
		return nil, errors.New("Invalid KeyGen Options")
	}

	key, err = keyGenerator.GenerateKey(opts)
	if err != nil {
		return nil, errors.New("Failed to generate a Key")
	}

	err = ch.keyManager.Store(key)
	if err != nil {
		return nil, errors.New("Failed to store a Key")
	}

	return

}

func (ch *cryptoHelper) LoadKey() (pri Key, pub Key, err error) {

	pri, pub, err = ch.keyManager.LoadKey()
	if err != nil {
		return nil, nil, errors.New("Failed to Load Key")
	}

	return

}

//func (ch *cryptoHelper) KeyImport(data interface{}, opts KeyImporterOpts) (key Key, err error) {
//
//	if data == nil {
//		return nil, errors.New("Data have not to be NIL")
//	}
//
//	if opts == nil {
//		return nil, errors.New("Invalid KeyImporter Opts")
//	}
//
//	keyImporter, found := ch.keyImporters[reflect.TypeOf(opts)]
//	if !found {
//		return nil, errors.New("Invalid KeyImporter Opts")
//	}
//
//	key, err = keyImporter.Import(data)
//	if err != nil {
//		return nil, errors.New("Failed to import key from input data")
//	}
//
//	return
//}