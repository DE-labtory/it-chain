package auth

type Signer interface {

	Sign(key Key, digest []byte, opts SignerOpt) (signature []byte, err error)

}

type Verifier interface {

	Verify(key Key, signature, digest []byte, opts SignerOpt) (valid bool, err error)

}