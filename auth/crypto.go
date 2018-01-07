package auth

import (
	"crypto"
)

type Key interface {}

type SignerOpt interface {
	crypto.SignerOpts
}

type Crypto interface {

	Sign(key Key, data []byte, opt SignerOpt) ([]byte, error)

	Verify(key Key, sig []byte, data []byte, opts SignerOpt) (bool, error)

}