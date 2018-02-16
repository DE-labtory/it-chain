package auth

import (
	"crypto/rsa"
	"crypto"
)

type Key interface {

	SKI() (ski []byte)

	Algorithm() string

}

type Crypto interface {

	Sign(digest []byte, opts SignerOpts) ([]byte, error)

	Verify(key Key, signature, digest []byte, opts SignerOpts) (bool, error)

	GetKey() (pri, pub Key, err error)

}

var DefaultRSAOption = &rsa.PSSOptions{SaltLength:rsa.PSSSaltLengthEqualsHash, Hash:crypto.SHA256}