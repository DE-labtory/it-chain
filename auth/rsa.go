package auth

import (
	"crypto/rsa"
	"errors"
	"crypto/rand"
)

type rsaSigner struct{}

func (s *rsaSigner) Sign(key Key, digest []byte, opts SignerOpts) (signature []byte, err error) {

	if opts == nil {
		return nil, errors.New("invalid options")
	}

	return key.(*rsaPrivateKey).priv.Sign(rand.Reader, digest, opts)

}

type rsaVerifier struct{}

func (v *rsaVerifier) Verify(key Key, signature, digest []byte, opts SignerOpts) (valid bool, err error) {

	if opts == nil {
		return false, errors.New("invalid options")
	}

	switch opts.(type) {
	case *rsa.PSSOptions:
		err := rsa.VerifyPSS(key.(*rsaPublicKey).pub,
			(opts.(*rsa.PSSOptions)).Hash,
				digest, signature, opts.(*rsa.PSSOptions))

		if err != nil {
			return false, errors.New("verifying error occurred")
		}

		return true, nil
	default:
		return false, errors.New("invalid options")
	}

}

type rsaPrivateKey struct {
	priv *rsa.PrivateKey
}

func (key *rsaPrivateKey) PublicKey() (pub Key, err error) {
	return &rsaPublicKey{&key.priv.PublicKey}, nil
}

type rsaPublicKey struct {
	pub *rsa.PublicKey
}