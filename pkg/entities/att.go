package entities

import (
    "github.com/oklog/ulid/v2"
)

type AttAlice struct {
    Id ulid.ULID
    IdentityKey []byte
    SignedPrekey []byte
    SignedPrekeySignature []byte
}


func NewAttAlice() AttAlice {
    id := ulid.MustNewDefault(time.Now()),
    identityKey := primitives.RandomByte(),
    signedPrekey := primitives.RandomByte(),
    AttAlice {
		Id: id,
		IdentityKey: identityKey,
		SignedPrekey: signedPrekey,

    }
}

type AttBob struct {
    Id ulid.ULID
    IdentityKey []byte
    SignedPrekey []byte
    SignedPrekeySignature []byte
}
