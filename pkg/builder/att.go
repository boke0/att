package builder

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"time"

	. "github.com/boke0/att/pkg/primitives"
	. "github.com/boke0/att/pkg/state"
	. "github.com/boke0/att/pkg/tree"
	"github.com/oklog/ulid/v2"
)

func BuildAttTree(states []AttState) Tree[AttState] {
	test := true
    sort.Slice(states, func(i, j int) bool { return states[i].Id() < states[j].Id() })
	for test {
		test = false
		for i := 0; i+3 < len(states); i += 2 {
			if !(states[i].IsActive() || states[i+1].IsActive()) {
				if states[i+2].IsActive() || states[i+3].IsActive() {
					t1 := states[i]
					t2 := states[i+1]
					states[i] = states[i+2]
					states[i+1] = states[i+3]
					states[i+2] = t1
					states[i+3] = t2
					test = true
				}
			}
			if (states[i].IsActive() || states[i+1].IsActive()) && (states[i+2].IsActive() || states[i+3].IsActive()){
				var (
					prevActivatedAt *time.Time
					nextActivatedAt *time.Time
				)
				if states[i].IsActive() {
					prevActivatedAt = states[i].ActivatedAt()
				}
				if states[i+1].IsActive() && (prevActivatedAt == nil || (prevActivatedAt != nil && prevActivatedAt.After(*states[i+1].ActivatedAt()))) {
					prevActivatedAt = states[i+1].ActivatedAt()
				}
				if states[i+2].IsActive() {
					nextActivatedAt = states[i+2].ActivatedAt()
				}
				if states[i+3].IsActive() && (prevActivatedAt == nil || (prevActivatedAt != nil && prevActivatedAt.After(*states[i+3].ActivatedAt()))) {
					nextActivatedAt = states[i+3].ActivatedAt()
				}
				if (prevActivatedAt == nil && nextActivatedAt != nil) || (prevActivatedAt != nil && nextActivatedAt != nil && prevActivatedAt.After(*nextActivatedAt)) {
					t1 := states[i]
					t2 := states[i+1]
					states[i] = states[i+2]
					states[i+1] = states[i+3]
					states[i+2] = t1
					states[i+3] = t2
					test = true
				}
			}
		}
	}
	id := Hash(states[0].Id() + states[1].Id())
	var treeNode TreeNode[AttState]
	if states[0].Alice != nil {
		treeNode = TreeNode[AttState]{
			Id: id,
			Left: &TreeNode[AttState]{
				Id: states[0].Id(),
				Peer: &states[0],
				PublicKey: nil,
			},
			Right: &TreeNode[AttState]{
				Id: states[1].Id(),
				Peer: &states[1],
				PublicKey: nil,
			},
		}
	}else{
		treeNode = TreeNode[AttState]{
			Id: id,
			Left: &TreeNode[AttState]{
				Id: states[0].Id(),
				Peer: &states[0],
				PublicKey: nil,
			},
			Right: &TreeNode[AttState]{
				Id: states[1].Id(),
				Peer: &states[1],
				PublicKey: nil,
			},
		}
	}
	if len(states) > 2 {
		for i := 2; i<len(states); i+= 2 {
			if len(states) > i + 1 {
				if states[i].IsActive() || states[i + 1].IsActive() {
					treeNode = addTwo(treeNode, states[i], states[i + 1])
				}else{
					treeNode = add(treeNode, states[i])
					treeNode = add(treeNode, states[i + 1])
				}
			}else{
				treeNode = add(treeNode, states[i])
			}
		}
	}
	tree := Tree[AttState] {
		Root: treeNode,
	}
	return tree
}

func insertToAtt(t TreeNode[AttState], state AttState) TreeNode[AttState] {
	if t.Left == nil || t.Right == nil {
		id := sha256.Sum256([]byte(t.Id + state.Id()))
		
		return TreeNode[AttState]{
			Id: hex.EncodeToString(id[:]),
			Left: &t,
			Right: &TreeNode[AttState]{
				Id: state.Id(),
				Peer: &state,
			},
		}
	}else if t.Left.Count() > t.Right.Count() && (isActive(*t.Right) || state.IsActive()) {
		right := insertToAtt(*t.Right, state)
		t.Right = &right
		return t
	}else{
		id := sha256.Sum256([]byte(t.Id + state.Id()))
		return TreeNode[AttState]{
			Id: hex.EncodeToString(id[:]),
			Left: &t,
			Right: &TreeNode[AttState]{
				Id: state.Id(),
				Peer: &state,
			},
		}

	}
}

func addTwo(t TreeNode[AttState], state1 AttState, state2 AttState) TreeNode[AttState] {
	t = insertToAtt(t, state1)
	t = insertToAtt(t, state2)
    return t
}

func add(t TreeNode[AttState], state AttState) TreeNode[AttState] {
    idBytes := sha256.Sum256([]byte(t.Id + state.Id()))
	return TreeNode[AttState] {
		Id: hex.EncodeToString(idBytes[:]),
		Left: &t,
		Right: &TreeNode[AttState] {
			Id: state.Id(),
			Peer: &state,
		},
	}
}

func isActive(t TreeNode[AttState]) bool {
	if t.Peer != nil {
		state, _ := t.Peer.(*AttState)
		return state.IsActive()
	}else{
		return isActive(*t.Left) || isActive(*t.Right)
	}
}

func AttachAttKeys(tree *Tree[AttState], publicKeys map[string]PublicKey) {
	attachAttKeys(&tree.Root, publicKeys)
}

func attachAttKeys(treeNode *TreeNode[AttState], keys map[string]PublicKey) {
	k := keys[treeNode.Id]
	if k != nil {
		treeNode.PublicKey = &k
	}
	if treeNode.Left != nil {
		attachAttKeys(treeNode.Left, keys)
	}
	if treeNode.Right != nil {
		attachAttKeys(treeNode.Right, keys)
	}
}

func PrintTree(node *TreeNode[AttState], space int) {
	if node == nil {
		return
	}
	space += 2
	PrintTree(node.Right, space)
	for i := 0; i < space; i++ {
		print(" ")
	}
	if node.PublicKey != nil {
		fmt.Printf("%s %x\n", node.Id[:26], (*node.PublicKey)[:8])
	}else if node.Peer != nil {
		peer := node.Peer.(*AttState)
		if node.Peer.IsAlice() {
			fmt.Printf("a: %d %x %t\n", ulid.MustParse(node.Id).Time(), node.Peer.PublicKey()[:8], peer.IsActive())
		}else{
			fmt.Printf("b: %d %x %t\n", ulid.MustParse(node.Id).Time(), node.Peer.PublicKey()[:8], peer.IsActive())
		}
	}else{
		fmt.Printf("%s\n", node.Id[:26])
	}
	println("")
	PrintTree(node.Left, space)
}
