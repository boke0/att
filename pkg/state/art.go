package state

import (
	"crypto/sha256"

	"github.com/boke0/att/pkg/primitives"
	. "github.com/boke0/att/pkg/primitives"
)

type ArtState struct {
	Alice *ArtAliceState
	Bob   *ArtBobState
}

func (a ArtState) Id() string {
	if a.Alice != nil {
		return a.Alice.Id
	} else {
		return a.Bob.Id
	}
}

func (a ArtState) IsAlice() bool {
	if a.Alice != nil {
		return true
	} else {
		return false
	}
}

func (a ArtState) PrivateKey() *primitives.PrivateKey {
	if a.Alice != nil {
		k := a.Alice.PrivateKey()
		return &k
	} else {
		k := a.Bob.PrivateKey()
		return &k
	}
}

func (a ArtState) PublicKey() primitives.PublicKey {
	if a.Alice != nil {
		return a.Alice.PublicKey()
	} else {
		return a.Bob.PublicKey()
	}
}

type ArtAliceState struct {
	Id                    string
	EphemeralKey          primitives.PrivateKey
	SetupKey              primitives.PublicKey
	IsInitiator           bool
}

func (a ArtAliceState) PrivateKey() primitives.PrivateKey {
	if a.IsInitiator {
		return a.EphemeralKey
	}else{
		result := DiffieHellman(a.EphemeralKey, a.SetupKey)
		hashed := sha256.Sum256(result)
		key := primitives.PrivateKey(hashed[:])
		return key
	}
}

func (a ArtAliceState) PublicKey() primitives.PublicKey {
	return AsPublic(a.PrivateKey())
}

type ArtBobState struct {
	Id                    string
	EphemeralKey          primitives.PublicKey
	SetupKey              primitives.PrivateKey
	IsInitiator           bool
}

func (a ArtBobState) PrivateKey() primitives.PrivateKey {
	result :=DiffieHellman(a.SetupKey, a.EphemeralKey)
	hashed := sha256.Sum256(result)
	key := primitives.PrivateKey(hashed[:])
	return key
}

func (a ArtBobState) PublicKey() primitives.PublicKey {
	if a.IsInitiator {
		return a.EphemeralKey
	}else{
		return AsPublic(a.PrivateKey())
	}
}
