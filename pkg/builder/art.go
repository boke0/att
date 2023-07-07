package builder

import (
	"fmt"
	"sort"

	"github.com/boke0/att/pkg/primitives"
	. "github.com/boke0/att/pkg/primitives"
	. "github.com/boke0/att/pkg/state"
	. "github.com/boke0/att/pkg/tree"
	"github.com/oklog/ulid/v2"
)

func BuildArtTree(states []ArtState, keys map[string]primitives.PublicKey) Tree[ArtState] {
    sort.Slice(states, func(i, j int) bool { return states[i].Id() < states[j].Id() })
	
	id := Hash(states[0].Id() + states[1].Id())
	var treeNode TreeNode[ArtState]
	treeNode = TreeNode[ArtState]{
		Id: id,
		Count: 2,
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
	if k, ok := keys[treeNode.Id]; ok {
		if !treeNode.IsAliceSide() {
			treeNode.PublicKey = &k
		}
	}
	if len(states) > 2 {
		for i := 2; i<len(states); i+= 1 {
			treeNode = insertToArt(treeNode, states[i], keys)
		}
	}
	tree := Tree[ArtState] {
		Root: treeNode,
	}
	return tree
}

func insertToArt(t TreeNode[ArtState], state ArtState, keys map[string]primitives.PublicKey) TreeNode[ArtState] {
	if t.Left == nil || t.Right == nil {
		id := Hash(t.Id + state.Id())
		
		treeNode := TreeNode[ArtState]{
			Id: id,
			Count: t.Count + 1,
			Left: &t,
			Right: &TreeNode[ArtState]{
				Count: 1,
				Id: state.Id(),
				Peer: state,
			},
		}

		if k, ok := keys[treeNode.Id]; ok {
			if !treeNode.IsAliceSide() {
				treeNode.PublicKey = &k
			}
		}
		return treeNode
	}else if t.Left.Count > t.Right.Count {
		right := insertToArt(*t.Right, state, keys)
		t.Right = &right
		t.Count = t.Left.Count + t.Right.Count
		if k, ok := keys[t.Id]; ok {
			if !t.IsAliceSide() {
				t.PublicKey = &k
			}
		}
		return t
	}else{
		id := Hash(t.Id + state.Id())
		treeNode := TreeNode[ArtState]{
			Id: id,
			Count: t.Count + 1,
			Left: &t,
			Right: &TreeNode[ArtState]{
				Count: 1,
				Id: state.Id(),
				Peer: state,
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

func AttachArtKeys(tree *Tree[ArtState], publicKeys map[string]PublicKey) {
	attachArtKeys(&tree.Root, publicKeys)
}

func attachArtKeys(treeNode *TreeNode[ArtState], keys map[string]PublicKey) {
	k := keys[treeNode.Id]
	if k != nil {
		treeNode.PublicKey = &k
	}
	if treeNode.Left != nil {
		attachArtKeys(treeNode.Left, keys)
	}
	if treeNode.Right != nil {
		attachArtKeys(treeNode.Right, keys)
	}
}

func PrintArtTree(node *TreeNode[ArtState], space int) {
	if node == nil {
		return
	}
	space += 2
	PrintArtTree(node.Right, space)
	for i := 0; i < space; i++ {
		print(" ")
	}
	if node.PublicKey != nil {
		fmt.Printf("%s %x\n", node.Id[:26], (*node.PublicKey)[:8])
	}else if node.Peer != nil {
		if node.Peer.IsAlice() {
			fmt.Printf("a: %d\n", ulid.MustParse(node.Id).Time())
		}else{
			fmt.Printf("b: %d\n", ulid.MustParse(node.Id).Time())
		}
	}else{
		fmt.Printf("%s\n", node.Id[:26])
	}
	println("")
	PrintArtTree(node.Left, space)
}
