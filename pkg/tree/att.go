package tree

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/boke0/att/pkg/primitives"
	. "github.com/boke0/att/pkg/state"
)


func (t TreeNode[State]) InsertToAtt(state AttState, keys map[string]primitives.PublicKey) TreeNode[State] {
	if t.Left == nil || t.Right == nil {
		id := sha256.Sum256([]byte(t.Id + state.Id()))
		
		treeNode := TreeNode[State]{
			Id: hex.EncodeToString(id[:]),
			IsActive: state.IsActive() || t.IsActive,
			Count: 2,
			Left: &t,
			Right: &TreeNode[State]{
				Id: state.Id(),
				IsActive: state.IsActive(),
				Count: 1,
				Peer: state,
			},
		}
		if k, ok := keys[treeNode.Id]; ok {
			if !treeNode.IsAliceSide() {
				treeNode.PublicKey = &k
			}
		}

		return treeNode
	}else if t.Left.Count > t.Right.Count && ((*t.Right).IsActive || state.IsActive()) {
		right := t.Right.InsertToAtt(state, keys)
		t.Right = &right
		t.Count = t.Left.Count + t.Right.Count
		t.IsActive = t.Left.IsActive || t.Right.IsActive

		if k, ok := keys[t.Id]; ok {
			if !t.IsAliceSide() {
				t.PublicKey = &k
			}
		}
		return t
	}else{
		id := sha256.Sum256([]byte(t.Id + state.Id()))
		treeNode := TreeNode[State]{
			Id: hex.EncodeToString(id[:]),
			IsActive: t.IsActive || state.IsActive(),
			Count: t.Count + 1,
			Left: &t,
			Right: &TreeNode[State]{
				Id: state.Id(),
				IsActive: state.IsActive(),
				Count: 1,
				Peer: &state,
			},
		}
		if k, ok := keys[treeNode.Id]; ok {
			if !treeNode.IsAliceSide() {
				treeNode.PublicKey = &k
			}
		}
		return treeNode

	}
}

func (t TreeNode[State]) AddTwo(state1 AttState, state2 AttState, keys map[string]primitives.PublicKey) TreeNode[State] {
	t = t.InsertToAtt(state1, keys).InsertToAtt(state2, keys)
    return t
}

func (t TreeNode[State]) Add(state AttState, keys map[string]primitives.PublicKey) TreeNode[State] {
    idBytes := sha256.Sum256([]byte(t.Id + state.Id()))
	treeNode := TreeNode[State] {
		Id: hex.EncodeToString(idBytes[:]),
		IsActive: t.IsActive || state.IsActive(),
		Count: t.Count + 1,
		Left: &t,
		Right: &TreeNode[State] {
			Id: state.Id(),
			IsActive: state.IsActive(),
			Count: t.Count + 1,
			Peer: &state,
		},
	}
	if k, ok := keys[treeNode.Id]; ok {
		if !treeNode.IsAliceSide() {
			treeNode.PublicKey = &k
		}
	}
	return treeNode
}

func IsActive(t TreeNode[AttState]) bool {
	if t.Peer != nil {
		state, _ := t.Peer.(*AttState)
		return state.IsActive()
	}else{
		return IsActive(*t.Left) || IsActive(*t.Right)
	}
}
