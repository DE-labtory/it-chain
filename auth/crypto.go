package auth

type Key interface {

	SKI() (ski []byte)

}

type Crypto interface {

	Sign(digest []byte, opts SignerOpts) (signature []byte, err error)

	Verify(key Key, signature, digest []byte, opts SignerOpts) (valid bool, err error)

	GenerateKey(opts KeyGenOpts) (pri, pub Key, err error)

	LoadKey() (pri, pub Key, err error)
}