package builder

import (
	"math/rand"
	"github.com/boke0/att/pkg/state"
	"github.com/oklog/ulid/v2"
	"testing"
	"time"
)

func TestBuildAttTree(t *testing.T) {
	entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)

	alice := state.AttAliceState{
		Id:                    ulid.MustNew(ulid.Timestamp(time.Now()), entropy),
		EphemeralKey:          []byte{0x01},
		EphemeralKeySignature: []byte{0x01},
		IsActive:              true,
	}
	bob := state.AttBobState{
		Id:                    ulid.MustNew(ulid.Timestamp(time.Now()), entropy),
		EphemeralKey:          []byte{0x01},
		EphemeralKeySignature: []byte{0x01},
		IsActive:              false,
	}
	charly := state.AttBobState{
		Id:                    ulid.MustNew(ulid.Timestamp(time.Now()), entropy),
		EphemeralKey:          []byte{0x01},
		EphemeralKeySignature: []byte{0x01},
		IsActive:              false,
	}
	dave := state.AttBobState{
		Id:                    ulid.MustNew(ulid.Timestamp(time.Now()), entropy),
		EphemeralKey:          []byte{0x01},
		EphemeralKeySignature: []byte{0x01},
		IsActive:              false,
	}

	states := []state.AttState{
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

	tree := BuildAttTree(states)


	t.Run("unbalanced", func(t *testing.T) {
		if tree.Root.Right.Id != dave.Id.String() {
			t.Errorf("Root's right node id was %s (wants %s)", tree.Root.Right.Id, dave.Id)
		}
		if tree.Root.Left.Right.Id != charly.Id.String() {
			t.Errorf("Root's left's right node id was %s (wants %s)", tree.Root.Left.Right.Id, charly.Id)
		}
		if tree.Root.Left.Left.Right.Id != bob.Id.String() {
			t.Errorf("Root's left's left's right node id was %s (wants %s)", tree.Root.Left.Left.Right.Id, bob.Id)
		}
		if tree.Root.Left.Left.Left.Id != alice.Id.String() {
			t.Errorf("Root's left's left's left node id was %s (wants %s)", tree.Root.Left.Left.Left.Id, alice.Id)
		}
	})

	charly.IsActive = true
	states = []state.AttState{
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

	tree = BuildAttTree(states)

	t.Run("balanced", func(t *testing.T) {
		if tree.Root.Right.Right.Id != dave.Id.String() {
			t.Errorf("Root's right's right node id was %s (wants %s)", tree.Root.Right.Right.Id, dave.Id)
		}
		if tree.Root.Right.Left.Id != charly.Id.String() {
			t.Errorf("Root's right's left node id was %s (wants %s)", tree.Root.Right.Left.Id, charly.Id)
		}
		if tree.Root.Left.Right.Id != bob.Id.String() {
			t.Errorf("Root's left's right's right node id was %s (wants %s)", tree.Root.Left.Right.Id, bob.Id)
		}
		if tree.Root.Left.Left.Id != alice.Id.String() {
			t.Errorf("Root's left's left's left node id was %s (wants %s)", tree.Root.Left.Left.Id, alice.Id)
		}
	})
}
