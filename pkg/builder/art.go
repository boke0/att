package builder

import (
	"crypto/sha256"
	"fmt"
	"sort"

	"github.com/boke0/att/pkg/primitives"
	. "github.com/boke0/att/pkg/primitives"
	. "github.com/boke0/att/pkg/state"
	. "github.com/boke0/att/pkg/tree"
	"github.com/oklog/ulid/v2"
)

func BuildArtTree(states []ArtState, keys map[string]primitives.PublicKey, initialize bool) Tree[ArtState] {
    sort.Slice(states, func(i, j int) bool { return states[i].Id() < states[j].Id() })
	
	id := Hash(states[0].Id() + states[1].Id())
	var treeNode TreeNode[ArtState]
	var (
		lpk primitives.PublicKey
		rpk primitives.PublicKey
	)
	if k, ok := keys[states[0].Id()]; ok && !states[0].IsAlice() {
		lpk = k
	}else{
		lpk = states[0].PublicKey()
	}
	if k, ok := keys[states[1].Id()]; ok && !states[1].IsAlice() {
		rpk = k
	}else{
		rpk = states[1].PublicKey()
	}
	treeNode = TreeNode[ArtState]{
		Id: id,
		Count: 2,
		Left: &TreeNode[ArtState]{
			Id: states[0].Id(),
			Peer: &states[0],
			PrivateKey: states[0].PrivateKey(),
			PublicKey: &lpk,
		},
		Right: &TreeNode[ArtState]{
			Id: states[1].Id(),
			Peer: &states[1],
			PrivateKey: states[1].PrivateKey(),
			PublicKey: &rpk,
		},
	}
	if k, ok := keys[treeNode.Id]; ok {
		if !treeNode.IsAliceSide() {
			treeNode.PublicKey = &k
		}
	}else if !treeNode.IsAliceSide() && initialize {
		result := primitives.DiffieHellman(
			*treeNode.Left.Peer.PrivateKey(),
			treeNode.Right.Peer.PublicKey(),
		)
        hashed := sha256.Sum256(result)
		key := primitives.PrivateKey(hashed[:])
        pub := primitives.AsPublic(key)
		treeNode.PrivateKey = &key
		treeNode.PublicKey = &pub
	}
	if len(states) > 2 {
		for i := 2; i<len(states); i+= 1 {
			treeNode = insertToArt(treeNode, states[i], keys, initialize)
		}
	}
	tree := Tree[ArtState] {
		Root: treeNode,
	}
	return tree
}

func insertToArt(t TreeNode[ArtState], state ArtState, keys map[string]primitives.PublicKey, initialize bool) TreeNode[ArtState] {
	if t.Left == nil || t.Right == nil {
		id := Hash(t.Id + state.Id())
		
		var pk primitives.PublicKey
		if k, ok := keys[state.Id()]; ok && !state.IsAlice() {
			pk = k
		}else{
			pk = state.PublicKey()
		}
		treeNode := TreeNode[ArtState]{
			Id: id,
			Count: t.Count + 1,
			Left: &t,
			Right: &TreeNode[ArtState]{
				Count: 1,
				Id: state.Id(),
				Peer: &state,
				PrivateKey: state.PrivateKey(),
				PublicKey: &pk,
			},
		}

		if k, ok := keys[treeNode.Id]; ok {
			if !treeNode.IsAliceSide() {
				treeNode.PublicKey = &k
			}
		}else if !treeNode.IsAliceSide() && initialize {
			result := primitives.DiffieHellman(
				*treeNode.Left.PrivateKey,
				*treeNode.Right.PublicKey,
			)
			hashed := sha256.Sum256(result)
			key := primitives.PrivateKey(hashed[:])
			pub := primitives.AsPublic(key)
			treeNode.PrivateKey = &key
			treeNode.PublicKey = &pub
		}
		return treeNode
	}else if t.Left.Count > t.Right.Count {
		right := insertToArt(*t.Right, state, keys, initialize)
		t.Right = &right
		t.Count = t.Left.Count + t.Right.Count
		if k, ok := keys[t.Id]; ok {
			if !t.IsAliceSide() {
				t.PublicKey = &k
			}
		}else if !t.IsAliceSide() && initialize {
			result := primitives.DiffieHellman(
				*t.Left.PrivateKey,
				*t.Right.PublicKey,
			)
			hashed := sha256.Sum256(result)
			key := primitives.PrivateKey(hashed[:])
			pub := primitives.AsPublic(key)
			t.PrivateKey = &key
			t.PublicKey = &pub
		}
		return t
	}else{
		id := Hash(t.Id + state.Id())
		var pk primitives.PublicKey
		if k, ok := keys[state.Id()]; ok && !state.IsAlice() {
			pk = k
		}else{
			pk = state.PublicKey()
		}
		treeNode := TreeNode[ArtState]{
			Id: id,
			Count: t.Count + 1,
			Left: &t,
			Right: &TreeNode[ArtState]{
				Count: 1,
				Id: state.Id(),
				Peer: &state,
				PrivateKey: state.PrivateKey(),
				PublicKey: &pk,
			},
		}
		if k, ok := keys[treeNode.Id]; ok {
			if !treeNode.IsAliceSide() {
				treeNode.PublicKey = &k
			}
		}else if !treeNode.IsAliceSide() && initialize {
			result := primitives.DiffieHellman(
				*treeNode.Left.PrivateKey,
				*treeNode.Right.PublicKey,
			)
			hashed := sha256.Sum256(result)
			key := primitives.PrivateKey(hashed[:])
			pub := primitives.AsPublic(key)
			treeNode.PrivateKey = &key
			treeNode.PublicKey = &pub
		}
		return treeNode
	}
}

func AttachArtKeys(tree *Tree[ArtState], publicKeys map[string]PublicKey) {
	attachArtKeys(&tree.Root, publicKeys)
}

func GetAllArtPublicKeys(tree *Tree[ArtState]) map[string]PublicKey {
	return getAllArtPublicKeys(&tree.Root)
}

func getAllArtPublicKeys(treeNode *TreeNode[ArtState]) map[string]PublicKey {
	keys := map[string]PublicKey{}
	if treeNode.Left != nil {
		keys_ := getAllArtPublicKeys(treeNode.Left)
		for nid, k := range keys_ {
			keys[nid] = k
		}
	}
	if treeNode.Right != nil {
		keys_ := getAllArtPublicKeys(treeNode.Right)
		for nid, k := range keys_ {
			keys[nid] = k
		}
	}
	if treeNode.PublicKey != nil {
		keys[treeNode.Id] = *treeNode.PublicKey
	}
	return keys
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
		if node.Peer != nil {
			if node.Peer.IsAlice() {
				fmt.Printf("a: %d %x %x %t\n", ulid.MustParse(node.Id).Time(), primitives.AsPublic((*(*node).PrivateKey)), ((*node).Peer.PublicKey()), node.Peer.(*ArtState).Alice.IsInitiator)
			}else{
				fmt.Printf("b: %d %x %x %t\n", ulid.MustParse(node.Id).Time(), (*(*node).PublicKey), ((*node).Peer.PublicKey()), node.Peer.(*ArtState).Bob.IsInitiator)
			}
		}else{
			fmt.Printf("%s %x %d\n", node.Id[:8], (*node.PublicKey), node.Count)
		}
	}else if node.IsAliceSide() {
		fmt.Printf("asite:%s %t %t %d\n", node.Id[:8], node.PrivateKey != nil, node.PublicKey != nil, node.Count)
	}else{
		fmt.Printf("%s %t %t %d\n", node.Id[:8], node.PrivateKey != nil, node.PublicKey != nil, node.Count)
	}
	println("")
	PrintArtTree(node.Left, space)
}
