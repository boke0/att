package tree

import (
    "testing"
	"math/rand"
	"time"
	"github.com/oklog/ulid/v2"
	"github.com/boke0/att/pkg/state"
	"github.com/boke0/att/pkg/primitives"
)

func TestExchange(t *testing.T) {
	entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)

	alice := state.ArtAliceState{
		Id:                    ulid.MustNew(ulid.Timestamp(time.Now()), entropy),
		EphemeralKey:          primitives.RandomByte(),
		EphemeralKeySignature: []byte{0x01},
	}

	bob := state.ArtBobState{
		Id:                    ulid.MustNew(ulid.Timestamp(time.Now()), entropy),
		EphemeralKey:          primitives.AsPublic(primitives.RandomByte()),
		EphemeralKeySignature: []byte{0x01},
	}
	charly := state.ArtBobState{
		Id:                    ulid.MustNew(ulid.Timestamp(time.Now()), entropy),
		EphemeralKey:          primitives.AsPublic(primitives.RandomByte()),
		EphemeralKeySignature: []byte{0x01},
	}

	tree := Tree[state.ArtState] {
		Root: TreeNode[state.ArtState]{
			Left: &TreeNode[state.ArtState]{
				Left: &TreeNode[state.ArtState]{
					Id:   alice.Id.String(),
					Peer: &state.ArtState{ Alice: &alice },
				},
				Right: &TreeNode[state.ArtState]{
					Id:   bob.Id.String(),
					Peer: &state.ArtState{ Bob: &bob },
				},
			},
			Right: &TreeNode[state.ArtState] {
				Id:   charly.Id.String(),
				Peer: &state.ArtState{ Bob: &charly },
			},
		},
	}
	t.Run("exchange", func(t *testing.T) {
		if tree.Root.Right.Id != charly.Id.String() {
			t.Errorf("Root's right's node id was %s (wants %s)", tree.Root.Right.Left.Id, charly.Id)
		}
		if tree.Root.Left.Right.Id != bob.Id.String() {
			t.Errorf("Root's left's right node id was %s (wants %s)", tree.Root.Left.Right.Id, bob.Id)
		}
		if tree.Root.Left.Left.Id != alice.Id.String() {
			t.Errorf("Root's left's left node id was %s (wants %s)", tree.Root.Left.Left.Id, alice.Id)
		}
	})
}
