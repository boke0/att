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
        k := a.Alice.PrivateKey()
        return &k
    }else{
        k := a.Bob.PrivateKey()
        return &k
    }
}

func (a ArtState) PublicKey() []byte {
    if a.Alice != nil {
        return a.Alice.PublicKey()
    }else{
        return a.Bob.PublicKey()
    }
}

type ArtAliceState struct {
    Id ulid.ULID
    EphemeralKey []byte
    EphemeralKeySignature []byte
    SetupKey []byte
    IsInitiator bool
}

func (a ArtAliceState) PrivateKey() []byte {
    return DiffieHellman(a.EphemeralKey, a.SetupKey)
}

func (a ArtAliceState) PublicKey() []byte {
    return AsPublic(a.PrivateKey())
}

type ArtBobState struct {
    Id ulid.ULID
    EphemeralKey []byte
    EphemeralKeySignature []byte
    SetupKey []byte
    IsInitiator bool
}

func (a ArtBobState) PrivateKey() []byte {
    return DiffieHellman(a.EphemeralKey, a.SetupKey)
}

func (a ArtBobState) PublicKey() []byte {
    return AsPublic(a.PrivateKey())
}
