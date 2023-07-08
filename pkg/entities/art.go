package entities

import (
	crand "crypto/rand"
	"crypto/sha256"
	"encoding/hex"

	"github.com/boke0/att/pkg/builder"
	"github.com/boke0/att/pkg/messages"
	"github.com/boke0/att/pkg/primitives"
	"github.com/boke0/att/pkg/state"
)

type ArtAlice struct {
	Id                    string
	IdentityKey           []byte
	SignedPrekey          []byte
	SignedPrekeySignature []byte
	states                map[string]state.ArtState
	keys                  map[string]primitives.PublicKey
}

func NewArtAlice(id string) ArtAlice {
	identityKey := primitives.RandomByte()
	signedPrekey := primitives.RandomByte()
	sig, _ := primitives.Sign(crand.Reader, identityKey, primitives.AsPublic(signedPrekey))
	return ArtAlice{
		Id:                    id,
		IdentityKey:           identityKey,
		SignedPrekey:          signedPrekey,
		SignedPrekeySignature: sig,
		states:                make(map[string]state.ArtState),
		keys:                  make(map[string]primitives.PublicKey),
	}
}

func (a *ArtAlice) Bob() ArtBob {
	return ArtBob{
		Id:                    a.Id,
		IdentityKey:           primitives.AsPublic(a.IdentityKey),
		SignedPrekey:          primitives.AsPublic(a.SignedPrekey),
		SignedPrekeySignature: a.SignedPrekeySignature,
		Alice:                 a,
	}
}

func (a *ArtAlice) Initialize(bobs map[string]ArtBob) messages.ArtMessage {
	a.states = map[string]state.ArtState{}
	sendKeys := map[string]messages.ArtPublicKey{}

	cnt := 0
	setupKey := primitives.RandomByte()
	publicSetupKey := primitives.AsPublic(setupKey)
	setupKeySignature, _ := primitives.Sign(crand.Reader, a.IdentityKey, publicSetupKey)

	astate := state.ArtState{
		Alice: &state.ArtAliceState{
			Id:           a.Id,
			SetupKey:     setupKey,
			EphemeralKey: a.SignedPrekey,
			IsInitiator:  true,
		},
	}
	a.states[a.Id] = astate

	for _, bob := range bobs {
		if ok := primitives.Verify(bob.IdentityKey, bob.SignedPrekey, bob.SignedPrekeySignature); !ok {
			panic("invalid signed prekey signature")
		}
		cnt += 1
		a.states[bob.Id] = state.ArtState{
			Bob: &state.ArtBobState{
				Id:           bob.Id,
				EphemeralKey: bob.SignedPrekey,
				IsInitiator:  false,
				SetupKey:     setupKey,
			},
		}
	}

	states := []state.ArtState{}
	for _, state := range a.states {
		states = append(states, state)
	}

	tree := builder.BuildArtTree(states, map[string]primitives.PublicKey{})
	_, keys := tree.DiffieHellman()
	a.keys = keys

	foundKeys := builder.GetAllArtPublicKeys(&tree)
	for nid, key_byte := range foundKeys {
		sig, _ := primitives.Sign(crand.Reader, a.IdentityKey, key_byte)
		sendKeys[nid] = messages.ArtPublicKey{
			SenderId:           a.Id,
			PublicKey:          hex.EncodeToString(key_byte),
			PublicKeySignature: hex.EncodeToString(sig),
		}
	}
	for nid, key := range keys {
		sig, _ := primitives.Sign(crand.Reader, a.IdentityKey, key)
		sendKeys[nid] = messages.ArtPublicKey{
			SenderId:           a.Id,
			PublicKey:          hex.EncodeToString(key),
			PublicKeySignature: hex.EncodeToString(sig),
		}
	}

	return messages.ArtMessage{
		InitializeMessage: &messages.ArtInitializeMessage{
			InitiatorId:       a.Id,
			SetupKey:          hex.EncodeToString(publicSetupKey),
			SetupKeySignature: hex.EncodeToString(setupKeySignature),
			Keys:              sendKeys,
		},
	}
}

func (a *ArtAlice) Send(mes []byte) messages.ArtMessage {
	ephemeral_key := primitives.RandomByte()
	public_ephemeral_key := primitives.AsPublic(ephemeral_key)
	ephemeral_key_signature, _ := primitives.Sign(crand.Reader, a.IdentityKey, public_ephemeral_key)

	a.states[a.Id].Alice.EphemeralKey = ephemeral_key
	if !a.states[a.Id].Alice.IsInitiator {
		result := primitives.DiffieHellman(ephemeral_key, a.states[a.Id].Alice.SetupKey)
		hashed := sha256.Sum256(result)
		key := primitives.PrivateKey(hashed[:])
		pub := primitives.AsPublic(key)
		a.keys[a.Id] = pub
	}

	states := []state.ArtState{}
	for _, state := range a.states {
		states = append(states, state)
	}

	tree := builder.BuildArtTree(states, a.keys)

	key, key_bytes := tree.DiffieHellman()

	keys := map[string]messages.ArtPublicKey{}
	for nid, key_byte := range key_bytes {
		sig, _ := primitives.Sign(crand.Reader, a.IdentityKey, key_byte)
		keys[nid] = messages.ArtPublicKey{
			SenderId:           a.Id,
			PublicKey:          hex.EncodeToString(key_byte),
			PublicKeySignature: hex.EncodeToString(sig),
		}
	}
	pk := a.states[a.Id].PublicKey()
	sig, _ := primitives.Sign(crand.Reader, a.IdentityKey, pk)
	keys[a.Id] = messages.ArtPublicKey{
		SenderId:           a.Id,
		PublicKey:          hex.EncodeToString(pk),
		PublicKeySignature: hex.EncodeToString(sig),
	}

	return messages.ArtMessage{
		TextMessage: &messages.ArtTextMessage{
			SenderId:              a.Id,
			EphemeralKey:          hex.EncodeToString(public_ephemeral_key),
			EphemeralKeySignature: hex.EncodeToString(ephemeral_key_signature),
			Keys:                  keys,
			Payload:               hex.EncodeToString(primitives.Encrypt(mes, key)),
		},
	}
}

func (a *ArtAlice) Receive(mes messages.ArtMessage, bobs map[string]ArtBob) []byte {
	if mes.InitializeMessage != nil {
		{
			setupKey, _ := hex.DecodeString(mes.InitializeMessage.SetupKey)
			setupKeySignature, _ := hex.DecodeString(mes.InitializeMessage.SetupKeySignature)
			if ok := primitives.Verify(bobs[mes.InitializeMessage.InitiatorId].IdentityKey, setupKey, setupKeySignature); !ok {
				panic("invalid initialize key signature")
			}
			a.states[a.Id] = state.ArtState{
				Alice: &state.ArtAliceState{
					Id: a.Id,
					SetupKey:     setupKey,
					EphemeralKey: a.SignedPrekey,
					IsInitiator: false,
				},
			}
			for bobId, bob := range bobs {
				if ok := primitives.Verify(bob.IdentityKey, bob.SignedPrekey, bob.SignedPrekeySignature); !ok {
					panic("invalid initialize key signature")
				}
				a.states[bobId] = state.ArtState{
					Bob: &state.ArtBobState{
						Id: bobId,
						SetupKey:     setupKey,
						EphemeralKey: bob.SignedPrekey,
						IsInitiator: mes.InitializeMessage.InitiatorId == bobId,
					},
				}
			}
			a.keys = make(map[string]primitives.PublicKey)
			for nid, key := range mes.InitializeMessage.Keys {
				pk, _ := hex.DecodeString(key.PublicKey)
				pk_sig, _ := hex.DecodeString(key.PublicKeySignature)
				if ok := primitives.Verify(bobs[key.SenderId].IdentityKey, pk, pk_sig); !ok {
					panic("invalid public key signature")
				}
				a.keys[nid] = pk
			}
			states := []state.ArtState{}
			for _, state := range a.states {
				states = append(states, state)
			}
		}
		return []byte{}
	} else if mes.TextMessage != nil {
		for nid, key := range mes.TextMessage.Keys {
			pk, _ := hex.DecodeString(key.PublicKey)
			pk_sig, _ := hex.DecodeString(key.PublicKeySignature)
			if ok := primitives.Verify(bobs[key.SenderId].IdentityKey, pk, pk_sig); !ok {
				panic("invalid public key signature")
			}
			a.keys[nid] = pk
		}
		ephemeralKey, _ := hex.DecodeString(mes.TextMessage.EphemeralKey)
		ephemeralKeySignature, _ := hex.DecodeString(mes.TextMessage.EphemeralKeySignature)
		if ok := primitives.Verify(bobs[mes.TextMessage.SenderId].IdentityKey, ephemeralKey, ephemeralKeySignature); !ok {
			panic("invalid public key signature")
		}

		a.states[mes.TextMessage.SenderId].Bob.EphemeralKey = ephemeralKey
		if a.states[mes.TextMessage.SenderId].Bob.IsInitiator {
			a.keys[mes.TextMessage.SenderId] = ephemeralKey
		}else if a.states[a.Id].Alice.IsInitiator {
			result := primitives.DiffieHellman(a.states[mes.TextMessage.SenderId].Bob.SetupKey, ephemeralKey)
			hashed := sha256.Sum256(result)
			key := primitives.PrivateKey(hashed[:])
			pub := primitives.AsPublic(key)
			a.keys[mes.TextMessage.SenderId] = pub
		}

		states := []state.ArtState{}
		for _, state := range a.states {
			states = append(states, state)
		}

		tree := builder.BuildArtTree(states, a.keys)
		key, _ := tree.DiffieHellman()

		cipher, _ := hex.DecodeString(mes.TextMessage.Payload)
		payload := primitives.Decrypt(cipher, key)

		return payload
	}
	return []byte{}
}

type ArtBob struct {
	Id                    string
	IdentityKey           []byte
	SignedPrekey          []byte
	SignedPrekeySignature []byte
	Alice                 *ArtAlice
}
