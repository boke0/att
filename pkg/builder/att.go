package builder

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"

	. "github.com/boke0/att/pkg/primitives"
	. "github.com/boke0/att/pkg/state"
	. "github.com/boke0/att/pkg/tree"
)

func BuildAttTree(states []AttState) Tree[AttState] {
	test := true
    sort.Slice(states, func(i, j int) bool { return states[i].Id() < states[j].Id() })
	for test {
		test = false
		for i := 0; i+1 < len(states); i += 2 {
			if !(states[i].IsActive() || states[i+1].IsActive()) && i+3 < len(states) {
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
				Peer: state,
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
				Peer: state,
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
	var (
		pub *[]byte
	)
    idBytes := sha256.Sum256([]byte(t.Id + state.Id()))
    return TreeNode[AttState] {
        Id: hex.EncodeToString(idBytes[:]),
        Peer: state,
        PublicKey: pub,
    }
}

func isActive(t TreeNode[AttState]) bool {
	if t.Peer != nil {
		state, _ := t.Peer.(AttState)
		return state.IsActive()
	}else{
		return isActive(*t.Left) || isActive(*t.Right)
	}
}
