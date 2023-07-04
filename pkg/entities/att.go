package entities

import (
    "crypto/rand"
	"github.com/boke0/att/pkg/primitives"
    "github.com/google/uuid"
)

type AttAlice struct {
    Id string
    IdentityKey primitives.PrivateKey
    SignedPrekey primitives.PrivateKey
    SignedPrekeySignature []byte
}


func NewAttAlice() AttAlice {
    id := uuid.NewUUID().String()
    identityKey := primitives.RandomByte()
    signedPrekey := primitives.RandomByte()
    sig, _ := primitives.Sign(rand.Reader, identityKey, signedPrekey)
    return AttAlice {
		Id: id,
		IdentityKey: identityKey,
		SignedPrekey: signedPrekey,
        SignedPrekeySignature: sig,
    }
}

func (a AttAlice) ToBob() AttBob {
    return AttBob {
        Id: a.Id,
        IdentityKey: primitives.AsPublic(a.IdentityKey),
        SignedPrekey: primitives.AsPublic(a.SignedPrekey),
        SignedPrekeySignature: a.SignedPrekeySignature,
    }
}

type AttBob struct {
    Id string
    IdentityKey primitives.PublicKey
    SignedPrekey primitives.PublicKey
    SignedPrekeySignature []byte
}
