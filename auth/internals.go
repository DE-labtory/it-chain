package auth

import "crypto"

type SignerOpts interface {
	crypto.SignerOpts
}

type signer interface {

	Sign(key Key, digest []byte, opts SignerOpts) ([]byte, error)

}

type verifier interface {

	Verify(key Key, signature, digest []byte, opts SignerOpts) (bool, error)

}

type KeyGenOpts interface {

	Algorithm() string

}

type keyGenerator interface {

	GenerateKey(opts KeyGenOpts) (pri, pub Key, err error)

}