package auth

type Key interface {}

type Crypto interface {

	Sign(key Key, digest []byte, opts SignerOpts) (signature []byte, err error)

	Verify(key Key, signature, digest []byte, opts SignerOpts) (valid bool, err error)

	KeyGenerate(opts KeyGenOpts) (key Key, err error)

}