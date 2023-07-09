package tree

import (
	"crypto/sha256"

	"github.com/boke0/att/pkg/primitives"
)

func (t TreeNode[ArtState]) InsertToArt(state ArtState, keys map[string]primitives.PublicKey, initialize bool) TreeNode[ArtState] {
	if t.Left == nil || t.Right == nil {
		id := primitives.Hash(t.Id + state.Id())
		
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
				Peer: state,
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
		right := t.Right.InsertToArt(state, keys, initialize)
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
		id := primitives.Hash(t.Id + state.Id())
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
				Peer: state,
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
