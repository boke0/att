package tree

import (
	"crypto/sha256"
    //"fmt"

	"github.com/boke0/att/pkg/primitives"
)

type Tree[Peer IPeer] struct {
    Root TreeNode[Peer]
}

type TreeNode[Peer IPeer] struct {
    Id string
    Peer IPeer
    PublicKey *primitives.PublicKey
    Left *TreeNode[Peer]
    Right *TreeNode[Peer]
    IsActive bool
    Count int
}

func (t Tree[IPeer]) DiffieHellman() ([]byte, map[string]primitives.PublicKey) {
    key, publicKeys := t.Root.DiffieHellman()
    delete(publicKeys, t.Root.Id)
    return key, publicKeys
}

func (t TreeNode[IPeer]) IsAlice() bool {
    return (t.Peer != nil && t.Peer.IsAlice())
}

func (t TreeNode[IPeer]) IsAliceSide() bool {
    return (t.Peer != nil && t.Peer.IsAlice()) || (t.Left != nil && t.Left.IsAliceSide()) || (t.Right != nil && t.Right.IsAliceSide())
}

func (t TreeNode[IPeer]) DiffieHellman() ([]byte, map[string]primitives.PublicKey) {
    if t.IsAlice() {
        return *t.Peer.PrivateKey(), make(map[string]primitives.PublicKey)
    }else if t.Peer != nil {
        return t.Peer.PublicKey(), make(map[string]primitives.PublicKey)
    }else if t.PublicKey != nil && !t.IsAliceSide() {
        return *t.PublicKey, make(map[string]primitives.PublicKey)
    }else{
        var (
            privateKey, publicKey []byte
            nodeLeftPublicKeys, nodeRightPublicKeys map[string]primitives.PublicKey
        )

        if t.Left.IsAliceSide() && !t.Right.IsAliceSide() {
            privateKey, nodeLeftPublicKeys = t.Left.DiffieHellman()
            publicKey, nodeRightPublicKeys = t.Right.DiffieHellman()
            if len(privateKey) == 0 {
                panic("private key is empty")
            }
            if len(publicKey) == 0 {
                panic("public key is empty")
            }
        }else if !t.Left.IsAliceSide() && t.Right.IsAliceSide() {
            privateKey, nodeLeftPublicKeys = t.Right.DiffieHellman()
            publicKey, nodeRightPublicKeys = t.Left.DiffieHellman()
            if len(privateKey) == 0 {
                panic("private key is empty")
            }
            if len(publicKey) == 0 {
                panic("public key is empty")
            }
        }else{
            panic("invalid tree structure")
        }

        //fmt.Printf("%x %x\n", primitives.AsPublic(privateKey), publicKey)

        result := primitives.DiffieHellman(privateKey, publicKey)
        key := sha256.Sum256(result)
        pub := primitives.AsPublic(key[:])
        nodePublicKeys := map[string]primitives.PublicKey {
            t.Id: pub,
        }
        nodePublicKeys = merge(nodePublicKeys, nodeLeftPublicKeys)
        nodePublicKeys = merge(nodePublicKeys, nodeRightPublicKeys)

        return key[:], nodePublicKeys
    }
}

/*func (t TreeNode[IPeer]) Count() int {
    if t.Peer != nil {
        return 1
    }else{
        return t.Left.Count() + t.Right.Count()
    }
}*/

func merge(m ...map[string]primitives.PublicKey) map[string]primitives.PublicKey {
    ans := make(map[string]primitives.PublicKey, 0)

    for _, c := range m {
        for k, v := range c {
            ans[k] = v
        }
    }
    return ans
}

func (tree *Tree[any]) AttachKeys(publicKeys map[string]primitives.PublicKey) {
	tree.Root.attachKeys(publicKeys)
}

func (treeNode *TreeNode[any]) attachKeys(keys map[string]primitives.PublicKey) {
	k := keys[treeNode.Id]
	if k != nil {
		treeNode.PublicKey = &k
	}
	if treeNode.Left != nil {
		treeNode.Left.attachKeys(keys)
	}
	if treeNode.Right != nil {
		treeNode.Right.attachKeys(keys)
	}
}

func (treeNode *TreeNode[IPeer]) IsIn(senderId string) bool {
	if treeNode.Id == senderId {
        return true
	}
	if treeNode.Left != nil && treeNode.Right != nil {
		return treeNode.Left.IsIn(senderId) || treeNode.Right.IsIn(senderId)
	}
    return false
}
