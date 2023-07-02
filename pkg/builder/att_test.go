package builder

import (
	"time"
	"testing"
	"github.com/boke0/att/pkg/entities"
	"github.com/boke0/att/pkg/states"
	"github.com/boke0/att/pkg/primitives"
	"github.com/oklog/ulid/v2"
)

func TestBuildAttTree(t *testing.T) {
	alice = entities.AttAlice {
		Id: ulid.MustNewDefault(time.Now()),
		IdentityKey: primitives.RandomByte(),
	}

	alice_state = alice.UpdateKey()
	bob_state = bob.ToBob().UpdateKey()
	charly_state = charly.ToBob().UpdateKey()
	dave_state = dave.ToBob().UpdateKey()

	states := []states.AttState {
		{
			Alice: 
		}
	}
	tree := BuildAttTree(states)
}
