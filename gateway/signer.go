package main

import (
	"crypto/sha512"

	"github.com/it-chain/bifrost/pb"
	"github.com/it-chain/heimdall/auth"
	"github.com/it-chain/heimdall/key"
)

type Signer struct {
	pri key.PriKey
}

func NewSigner(pri key.PriKey) *Signer {
	return &Signer{
		pri: pri,
	}
}

func (s Signer) SignEnvelope(envelope *pb.Envelope) *pb.Envelope {

	hash := sha512.New()
	hash.Write(envelope.Payload)
	digest := hash.Sum(nil)

	sig, err := auth.Sign(s.pri, digest, auth.EQUAL_SHA512.SignerOptsToPSSOptions())

	if err != nil {
		//signing error
		return nil
	}
	envelope.Signature = sig

	return envelope
}
