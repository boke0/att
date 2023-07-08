package builder

import (
	"math/rand"
	"testing"
	"time"

	"github.com/boke0/att/pkg/primitives"
	"github.com/boke0/att/pkg/state"
	"github.com/oklog/ulid/v2"
)

func TestBuildArtTree(t *testing.T) {
	entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)

	suk := primitives.RandomByte()
	alice := state.ArtAliceState{
		Id:                    ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String(),
		EphemeralKey:          primitives.RandomByte(),
		SetupKey: suk,
	}
	bob := state.ArtBobState{
		Id:                    ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String(),
		EphemeralKey:          primitives.AsPublic(primitives.RandomByte()),
		SetupKey: suk,
	}
	charly := state.ArtBobState{
		Id:                    ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String(),
		EphemeralKey:          primitives.AsPublic(primitives.RandomByte()),
		SetupKey: suk,
	}
	dave := state.ArtBobState{
		Id:                    ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String(),
		EphemeralKey:          primitives.AsPublic(primitives.RandomByte()),
		SetupKey: suk,
	}

	states := []state.ArtState{
		{
			Alice: &alice,
		},
		{
			Bob: &bob,
		},
		{
			Bob: &charly,
		},
		{
			Bob: &dave,
		},
	}

	tree := BuildArtTree(states, make(map[string]primitives.PublicKey))

	t.Run("balanced", func(t *testing.T) {
		if tree.Root.Right.Right.Id != dave.Id {
			t.Errorf("Root's right's right node id was %s (wants %s)", tree.Root.Right.Right.Id, dave.Id)
		}
		if tree.Root.Right.Left.Id != charly.Id {
			t.Errorf("Root's right's left node id was %s (wants %s)", tree.Root.Right.Left.Id, charly.Id)
		}
		if tree.Root.Left.Right.Id != bob.Id {
			t.Errorf("Root's left's right's right node id was %s (wants %s)", tree.Root.Left.Right.Id, bob.Id)
		}
		if tree.Root.Left.Left.Id != alice.Id {
			t.Errorf("Root's left's left's left node id was %s (wants %s)", tree.Root.Left.Left.Id, alice.Id)
		}
	})
}
