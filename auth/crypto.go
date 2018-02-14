package auth

import (
	"crypto/rsa"
	"crypto"
)

type Key interface {

	SKI() (ski []byte)

}

type Crypto interface {

	Sign(digest []byte, opts SignerOpts) ([]byte, error)

	Verify(key Key, signature, digest []byte, opts SignerOpts) (bool, error)

	GenerateKey(opts KeyGenOpts) (pri, pub Key, err error)

	LoadKey() (pri, pub Key, err error)
}

var DefaultRSAOption = &rsa.PSSOptions{SaltLength:rsa.PSSSaltLengthEqualsHash, Hash:crypto.SHA256}