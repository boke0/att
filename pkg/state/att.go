package state

import (
	"github.com/boke0/att/pkg/primitives"
	. "github.com/boke0/att/pkg/primitives"
	"github.com/oklog/ulid/v2"
)

type AttState struct {
    Alice *AttAliceState
    Bob *AttBobState
}

func (a AttState) Id() string {
    if a.Alice != nil {
        return a.Alice.Id.String()
    }else{
        return a.Bob.Id.String()
    }
}

func (a AttState) IsActive() bool {
    if a.Alice != nil {
        return a.Alice.IsActive
    }else{
        return a.Bob.IsActive
    }
}

func (a AttState) IsAlice() bool {
    if a.Alice != nil {
        return true
    }else{
        return false
    }
}

func (a AttState) PrivateKey() *primitives.PrivateKey {
    if a.Alice != nil {
        return &a.Alice.EphemeralKey
    }else{
        return nil
    }
}

func (a AttState) PublicKey() primitives.PublicKey {
    if a.Alice != nil {
        return AsPublic(a.Alice.EphemeralKey)
    }else{
        return a.Bob.EphemeralKey
    }
}

type AttAliceState struct {
    Id ulid.ULID
    EphemeralKey primitives.PrivateKey
    EphemeralKeySignature []byte
    IsActive bool
}

type AttBobState struct {
    Id ulid.ULID
    EphemeralKey primitives.PublicKey
    EphemeralKeySignature []byte
    IsActive bool
}
