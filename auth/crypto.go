package auth

type Key interface {}

type Crypto interface {

	Sign(key Key, digest []byte, opts SignerOpt) (signature []byte, err error)

	Verify(key Key, signature, digest []byte, opts SignerOpt) (valid bool, err error)

}