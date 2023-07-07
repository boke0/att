package entities

import (
	crand "crypto/rand"
	"encoding/hex"

	mrand "math/rand"
	"time"

	"github.com/boke0/att/pkg/builder"
	"github.com/boke0/att/pkg/messages"
	"github.com/boke0/att/pkg/primitives"
	"github.com/oklog/ulid/v2"

	//"github.com/boke0/att/pkg/builder"
	"github.com/boke0/att/pkg/state"
	"github.com/google/uuid"
)

type AttAlice struct {
	Id                    string
	IdentityKey           primitives.PrivateKey
	SignedPrekey          primitives.PrivateKey
	SignedPrekeySignature []byte
	states                map[string]state.AttState
	keys                  map[string]primitives.PublicKey
}

func NewAttAlice() AttAlice {
	id := uuid.NewString()
	identityKey := primitives.RandomByte()
	signedPrekey := primitives.RandomByte()
	sig, _ := primitives.Sign(crand.Reader, identityKey, primitives.AsPublic(signedPrekey))
	return AttAlice{
		Id:                    id,
		IdentityKey:           identityKey,
		SignedPrekey:          signedPrekey,
		SignedPrekeySignature: sig,
		states:                map[string]state.AttState{},
		keys:                  map[string]primitives.PublicKey{},
	}
}

func (a *AttAlice) Bob() AttBob {
	return AttBob{
		Id:                    a.Id,
		IdentityKey:           primitives.AsPublic(a.IdentityKey),
		SignedPrekey:          primitives.AsPublic(a.SignedPrekey),
		SignedPrekeySignature: a.SignedPrekeySignature,
		Alice:                 a,
	}
}

func (a *AttAlice) Initialize(bobs map[string]AttBob) messages.AttMessage {
	entropy := ulid.Monotonic(mrand.New(mrand.NewSource(time.Now().UnixNano())), 0)
	a.states = map[string]state.AttState{}
	users := map[string]messages.AttUserInitializeMessage{}

	cnt := 0
	initiatorRoomUserId := ulid.MustNew(uint64(cnt), entropy)
	initializeKey := primitives.RandomByte()
	publicInitializeKey := primitives.AsPublic(initializeKey)
	initializeKeySignature, _ := primitives.Sign(crand.Reader, a.IdentityKey, publicInitializeKey)

	astate := state.AttState{
		Alice: &state.AttAliceState{
			Id:                    initiatorRoomUserId,
			EphemeralKey:          initializeKey,
			EphemeralKeySignature: initializeKeySignature,
		},
	}
	a.states[a.Id] = astate
	users[a.Id] = messages.AttUserInitializeMessage{
		RoomUserId:     astate.Alice.Id.String(),
		SignedPrekeyId: primitives.HashBytes(publicInitializeKey),
	}

	for _, bob := range bobs {
		if ok := primitives.Verify(bob.IdentityKey, bob.SignedPrekey, bob.SignedPrekeySignature); !ok {
			panic("invalid signed prekey signature")
		}
		cnt += 1
		id := ulid.MustNew(uint64(cnt), entropy)
		a.states[bob.Id] = state.AttState{
			Bob: &state.AttBobState{
				Id:                    id,
				EphemeralKey:          bob.SignedPrekey,
				EphemeralKeySignature: bob.SignedPrekeySignature,
			},
		}
		users[bob.Id] = messages.AttUserInitializeMessage{
			RoomUserId:     id.String(),
			SignedPrekeyId: primitives.HashBytes(bob.SignedPrekey),
		}
	}

	return messages.AttMessage{
		InitializeMessage: &messages.AttInitializeMessage{
			InitiatorId:            a.Id,
			InitiatorRoomUserId:    initiatorRoomUserId.String(),
			InitializeKey:          hex.EncodeToString(publicInitializeKey),
			InitializeKeySignature: hex.EncodeToString(initializeKeySignature),
			Users:                  users,
		},
	}
}

func (a *AttAlice) Send(mes []byte) messages.AttMessage {
	ephemeral_key := primitives.RandomByte()
	public_ephemeral_key := primitives.AsPublic(ephemeral_key)
	ephemeral_key_signature, _ := primitives.Sign(crand.Reader, a.IdentityKey, public_ephemeral_key)

	a.states[a.Id].Alice.EphemeralKey = ephemeral_key
	a.states[a.Id].Alice.EphemeralKeySignature = ephemeral_key_signature
	if a.states[a.Id].Alice.ActivatedAt == nil {
		t := time.Now()
		a.states[a.Id].Alice.ActivatedAt = &t
	}

	states := []state.AttState{}
	for _, state := range a.states {
		states = append(states, state)
	}

	tree := builder.BuildAttTree(states, a.keys)

	key, key_bytes := tree.DiffieHellman()
	
	keys := map[string]messages.AttPublicKey{}
	for nid, key_byte := range key_bytes {
		sig, _ := primitives.Sign(crand.Reader, a.IdentityKey, key_byte)
		keys[nid] = messages.AttPublicKey{
			SenderId:           a.Id,
			PublicKey:          hex.EncodeToString(key_byte),
			PublicKeySignature: hex.EncodeToString(sig),
		}
	}

	return messages.AttMessage{
		TextMessage: &messages.AttTextMessage{
			SenderId:              a.Id,
			EphemeralKey:          hex.EncodeToString(public_ephemeral_key),
			EphemeralKeySignature: hex.EncodeToString(ephemeral_key_signature),
			Keys:                  keys,
			Payload:               hex.EncodeToString(primitives.Encrypt(mes, key)),
		},
	}
}

func (a *AttAlice) Receive(mes messages.AttMessage, bobs map[string]AttBob) []byte {
	if mes.InitializeMessage != nil {
		{
			ephemeralKey, _ := hex.DecodeString(mes.InitializeMessage.InitializeKey)
			ephemeralKeySignature, _ := hex.DecodeString(mes.InitializeMessage.InitializeKeySignature)
			if ok := primitives.Verify(bobs[mes.InitializeMessage.InitiatorId].IdentityKey, ephemeralKey, ephemeralKeySignature); !ok {
				panic("invalid initialize key signature")
			}
			a.states[mes.InitializeMessage.InitiatorId] = state.AttState{
				Bob: &state.AttBobState{
					Id:                    ulid.MustParse(mes.InitializeMessage.InitiatorRoomUserId),
					EphemeralKey:          ephemeralKey,
					EphemeralKeySignature: ephemeralKeySignature,
				},
			}
		}
		for user_uuid, user := range mes.InitializeMessage.Users {
			if user_uuid == a.Id {
				ephemeralKey := a.SignedPrekey
				ephemeralKeySignature := a.SignedPrekeySignature
				a.states[a.Id] = state.AttState{
					Alice: &state.AttAliceState{
						Id:                    ulid.MustParse(user.RoomUserId),
						EphemeralKey:          ephemeralKey,
						EphemeralKeySignature: ephemeralKeySignature,
					},
				}
			} else if user_uuid != mes.InitializeMessage.InitiatorId {
				bob := bobs[user_uuid]
				ephemeralKey := bob.SignedPrekey
				ephemeralKeySignature := bob.SignedPrekeySignature
				if ok := primitives.Verify(bob.IdentityKey, ephemeralKey, ephemeralKeySignature); !ok {
					panic("invalid initialize key signature")
				}
				a.states[bob.Id] = state.AttState{
					Bob: &state.AttBobState{
						Id:                    ulid.MustParse(user.RoomUserId),
						EphemeralKey:          ephemeralKey,
						EphemeralKeySignature: ephemeralKeySignature,
					},
				}
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
		a.states[mes.TextMessage.SenderId].Bob.EphemeralKeySignature = ephemeralKeySignature
		if a.states[mes.TextMessage.SenderId].Bob.ActivatedAt == nil {
			t := time.Now()
			a.states[mes.TextMessage.SenderId].Bob.ActivatedAt = &t
		}

		states := []state.AttState{}
		for _, state := range a.states {
			states = append(states, state)
		}

		tree := builder.BuildAttTree(states, a.keys)
		key, _ := tree.DiffieHellman()
		
		cipher, _ := hex.DecodeString(mes.TextMessage.Payload)
		payload := primitives.Decrypt(cipher, key)

		return payload
	}
	return []byte{}
}

type AttBob struct {
	Id                    string
	IdentityKey           primitives.PublicKey
	SignedPrekey          primitives.PublicKey
	SignedPrekeySignature []byte
	Alice                 *AttAlice
}
