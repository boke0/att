package entities

import (
    "github.com/google/uuid"
)

type ArtAlice struct {
    Id uuid.UUID
    IdentityKey []byte
    SignedPrekey []byte
    SignedPrekeySignature []byte
}

type ArtBob struct {
    Id uuid.UUID
    IdentityKey []byte
    SignedPrekey []byte
    SignedPrekeySignature []byte
}
