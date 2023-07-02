package state

import (
	. "github.com/boke0/att/pkg/primitives"
    "github.com/oklog/ulid/v2"
)

type ArtState struct {
    Alice *ArtAliceState
    Bob *ArtBobState
}

func (a ArtState) Id() string {
    if a.Alice != nil {
        return a.Alice.Id.String()
    }else{
        return a.Bob.Id.String()
    }
}

func (a ArtState) IsAlice() bool {
    if a.Alice != nil {
        return true
    }else{
        return false
    }
}

func (a ArtState) PrivateKey() *[]byte {
    if a.Alice != nil {
        return &a.Alice.EphemeralKey
    }else{
        return nil
    }
}

func (a ArtState) PublicKey() []byte {
    if a.Alice != nil {
        return AsPublic(a.Alice.EphemeralKey)
    }else{
        return a.Bob.EphemeralKey
    }
}

type ArtAliceState struct {
    Id ulid.ULID
    EphemeralKey []byte
    EphemeralKeySignature []byte
}

type ArtBobState struct {
    Id ulid.ULID
    EphemeralKey []byte
    EphemeralKeySignature []byte
}
