package entities

import (
	"crypto/rand"

	"github.com/boke0/att/pkg/primitives"
	"github.com/google/uuid"
)

type ArtAlice struct {
	Id                    string
	IdentityKey           []byte
	SignedPrekey          []byte
	SignedPrekeySignature []byte
}

func NewArtAlice() ArtAlice {
	id := uuid.NewString()
	identityKey := primitives.RandomByte()
	signedPrekey := primitives.RandomByte()
	sig, _ := primitives.Sign(rand.Reader, identityKey, signedPrekey)
	return ArtAlice{
		Id:                    id,
		IdentityKey:           identityKey,
		SignedPrekey:          signedPrekey,
		SignedPrekeySignature: sig,
	}
}

type ArtBob struct {
	Id                    string
	IdentityKey           []byte
	SignedPrekey          []byte
	SignedPrekeySignature []byte
	Alice                 *AttAlice
}
