package builder

import (
	"sort"

	. "github.com/boke0/att/pkg/primitives"
	. "github.com/boke0/att/pkg/state"
	. "github.com/boke0/att/pkg/tree"
)

func BuildArtTree(states []ArtState) Tree[ArtState] {
    sort.Slice(states, func(i, j int) bool { return states[i].Id() < states[j].Id() })
	
	id := Hash(states[0].Id() + states[1].Id())
	var treeNode TreeNode[ArtState]
	treeNode = TreeNode[ArtState]{
		Id: id,
		Left: &TreeNode[ArtState]{
			Id: states[0].Id(),
			Peer: &states[0],
			PublicKey: nil,
		},
		Right: &TreeNode[ArtState]{
			Id: states[1].Id(),
			Peer: &states[1],
			PublicKey: nil,
		},
	}
	if len(states) > 2 {
		for i := 2; i<len(states); i+= 1 {
			treeNode = insertToArt(treeNode, states[i])
		}
	}
	tree := Tree[ArtState] {
		Root: treeNode,
	}
	return tree
}

func insertToArt(t TreeNode[ArtState], state ArtState) TreeNode[ArtState] {
	if t.Left == nil || t.Right == nil {
		id := Hash(t.Id + state.Id())
		
		return TreeNode[ArtState]{
			Id: id,
			Left: &t,
			Right: &TreeNode[ArtState]{
				Id: state.Id(),
				Peer: state,
			},
		}
	}else if t.Left.Count() > t.Right.Count() {
		right := insertToArt(*t.Right, state)
		t.Right = &right
		return t
	}else{
		id := Hash(t.Id + state.Id())
		return TreeNode[ArtState]{
			Id: id,
			Left: &t,
			Right: &TreeNode[ArtState]{
				Id: state.Id(),
				Peer: state,
			},
		}

	}
}

func AttachKeys(tree *Tree[ArtState], publicKeys map[string]PublicKey) {
	attachKeys(&tree.Root, publicKeys)
}

func attachKeys(treeNode *TreeNode[ArtState], keys map[string]PublicKey) {
	k := keys[treeNode.Id]
	if k != nil {
		treeNode.PublicKey = &k
	}
	if treeNode.Left != nil {
		attachKeys(treeNode.Left, keys)
	}
	if treeNode.Right != nil {
		attachKeys(treeNode.Right, keys)
	}
}
