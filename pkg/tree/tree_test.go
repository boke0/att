package tree

import (
	"fmt"
	"bytes"
    "testing"
	"math/rand"
	"time"
	"github.com/oklog/ulid/v2"
	"github.com/boke0/att/pkg/state"
	"github.com/boke0/att/pkg/primitives"
)

func TestExchange(t *testing.T) {
	entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)

	alice := state.AttAliceState{
		Id:                    ulid.MustNew(ulid.Timestamp(time.Now()), entropy),
		EphemeralKey:          primitives.RandomByte(),
		EphemeralKeySignature: []byte{0x01},
	}
	bob := state.AttAliceState{
		Id:                    ulid.MustNew(ulid.Timestamp(time.Now()), entropy),
		EphemeralKey:          primitives.RandomByte(),
		EphemeralKeySignature: []byte{0x01},
	}
	charly := state.AttAliceState{
		Id:                    ulid.MustNew(ulid.Timestamp(time.Now()), entropy),
		EphemeralKey:          primitives.RandomByte(),
		EphemeralKeySignature: []byte{0x01},
	}

	fmt.Printf("a: %x\n", alice.EphemeralKey)
	fmt.Printf("b: %x\n", bob.EphemeralKey)
	fmt.Printf("c: %x\n", charly.EphemeralKey)

	aliceTree := Tree[state.AttState] {
		Root: TreeNode[state.AttState]{
			Id: primitives.Hash(primitives.Hash(alice.Id.String() + bob.Id.String()) + charly.Id.String()),
			Left: &TreeNode[state.AttState]{
				Id: primitives.Hash(alice.Id.String() + bob.Id.String()),
				Left: &TreeNode[state.AttState]{
					Id:   alice.Id.String(),
					Peer: &state.AttState{ Alice: &alice },
				},
				Right: &TreeNode[state.AttState]{
					Id:   bob.Id.String(),
					Peer: &state.AttState{ Bob: toBob(bob) },
				},
			},
			Right: &TreeNode[state.AttState] {
				Id:   charly.Id.String(),
				Peer: &state.AttState{ Bob: toBob(charly) },
			},
		},
	}
	fmt.Printf("aliceTree: %t\n", aliceTree.Root.Left.Left.Peer.IsAlice())
	bobTree := Tree[state.AttState] {
		Root: TreeNode[state.AttState]{
			Id: primitives.Hash(primitives.Hash(alice.Id.String() + bob.Id.String()) + charly.Id.String()),
			Left: &TreeNode[state.AttState]{
				Id: primitives.Hash(alice.Id.String() + bob.Id.String()),
				Left: &TreeNode[state.AttState]{
					Id:   alice.Id.String(),
					Peer: &state.AttState{ Bob: toBob(alice) },
				},
				Right: &TreeNode[state.AttState]{
					Id:   bob.Id.String(),
					Peer: &state.AttState{ Alice: &bob },
				},
			},
			Right: &TreeNode[state.AttState] {
				Id:   charly.Id.String(),
				Peer: &state.AttState{ Bob: toBob(charly) },
			},
		},
	}
	charlyTree := Tree[state.AttState] {
		Root: TreeNode[state.AttState]{
			Id: primitives.Hash(primitives.Hash(alice.Id.String() + bob.Id.String()) + charly.Id.String()),
			Left: &TreeNode[state.AttState]{
				Id: primitives.Hash(alice.Id.String() + bob.Id.String()),
				Left: &TreeNode[state.AttState]{
					Id:   alice.Id.String(),
					Peer: &state.AttState{ Bob: toBob(alice) },
				},
				Right: &TreeNode[state.AttState]{
					Id:   bob.Id.String(),
					Peer: &state.AttState{ Bob: toBob(bob) },
				},
			},
			Right: &TreeNode[state.AttState] {
				Id:   charly.Id.String(),
				Peer: &state.AttState{ Alice: &charly },
			},
		},
	}
	t.Run("exchange", func(t *testing.T) {
		aliceKey, keys := aliceTree.DiffieHellman()
		bobKey, _ := bobTree.DiffieHellman()

		charlyTree.AttachKeys(keys)
		charlyKey, _ := charlyTree.DiffieHellman()

		if !bytes.Equal(aliceKey, bobKey) {
			t.Errorf("alice and bob key exchange failed")
		}
		if !bytes.Equal(bobKey, charlyKey) {
			t.Errorf("bob and charly key exchange failed")
		}
		if !bytes.Equal(aliceKey, charlyKey) {
			t.Errorf("alice and charly key exchange failed")
		}
	})
}

func toBob(alice state.AttAliceState) *state.AttBobState {
	return &state.AttBobState {
		Id:                    alice.Id,
		EphemeralKey:          primitives.AsPublic(alice.EphemeralKey),
		EphemeralKeySignature: alice.EphemeralKeySignature,
	}
}
