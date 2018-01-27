package auth

import (
	"crypto/ecdsa"
	"math/big"
	"encoding/asn1"
	"errors"
	"crypto/rand"
	"crypto/elliptic"
	"crypto/sha256"
)

type ecdsaSignature struct {
	R, S *big.Int
}

type ecdsaSigner struct{}

func marshalECDSASignature(r, s *big.Int) (signature []byte, err error) {
	return asn1.Marshal(ecdsaSignature{r, s})
}

func unmarshalECDSASignature(signature []byte) (*big.Int, *big.Int, error) {
	ecdsaSig := new(ecdsaSignature)
	_, err := asn1.Unmarshal(signature, ecdsaSig)
	if err != nil {
		return nil, nil, errors.New("failed to unmarshal")
	}

	if ecdsaSig.R == nil {
		return nil, nil, errors.New("invalid signature")
	}
	if ecdsaSig.S == nil {
		return nil, nil, errors.New("invalid signature")
	}

	if ecdsaSig.R.Sign() != 1 {
		return nil, nil, errors.New("invalid signature")
	}
	if ecdsaSig.S.Sign() != 1 {
		return nil, nil, errors.New("invalid signature")
	}

	return ecdsaSig.R, ecdsaSig.S, nil
}

func (signer *ecdsaSigner) Sign(key Key, digest []byte, opts SignerOpts) ([]byte, error) {

	r, s, err := ecdsa.Sign(rand.Reader, key.(*ecdsaPrivateKey).priv, digest)
	if err != nil {
		return nil, err
	}

	signature, err := marshalECDSASignature(r, s)
	if err != nil {
		return nil, err
	}

	return signature, nil
}

type ecdsaVerifier struct{}

func (v *ecdsaVerifier) Verify(key Key, signature, digest []byte, opts SignerOpts) (bool, error) {

	r, s, err := unmarshalECDSASignature(signature)
	if err != nil {
		return false, err
	}

	valid := ecdsa.Verify(key.(*ecdsaPublicKey).pub, digest, r, s)
	if !valid {
		return valid, errors.New("failed to verify")
	}

	return valid, nil
}

type ecdsaPrivateKey struct {
	priv *ecdsa.PrivateKey
}

func (key *ecdsaPrivateKey) SKI() (ski []byte) {

	if key.priv == nil {
		return nil
	}

	data := elliptic.Marshal(key.priv.Curve, key.priv.PublicKey.X, key.priv.PublicKey.Y)

	hash := sha256.New()
	hash.Write(data)
	return hash.Sum(nil)

}

func (key *ecdsaPrivateKey) PublicKey() (pub Key, err error) {
	return &ecdsaPublicKey{&key.priv.PublicKey}, nil
}

type ecdsaPublicKey struct {
	pub *ecdsa.PublicKey
}

func (key *ecdsaPublicKey) SKI() (ski []byte) {

	if key.pub == nil {
		return nil
	}

	data := elliptic.Marshal(key.pub.Curve, key.pub.X, key.pub.Y)

	hash := sha256.New()
	hash.Write(data)
	return hash.Sum(nil)

}