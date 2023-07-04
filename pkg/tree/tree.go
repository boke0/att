package tree

import (
	"crypto/sha256"

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
}

func (t Tree[IPeer]) DiffieHellman() ([]byte, map[string]primitives.PublicKey) {
    key, publicKeys := t.DiffieHellman()
    delete(publicKeys, t.Root.Id)
    return key, publicKeys
}

func (t TreeNode[IPeer]) IsAlice() bool {
    return (t.Peer != nil && t.Peer.IsAlice())
}

func (t TreeNode[IPeer]) IsAliceSide() bool {
    return (t.Peer != nil && t.Peer.IsAlice()) || t.Left.IsAlice() || t.Right.IsAlice()
}

func (t TreeNode[IPeer]) DiffieHellman() ([]byte, map[string]primitives.PublicKey) {
    if t.IsAlice() {
        return *t.Peer.PrivateKey(), make(map[string]primitives.PublicKey)
    }else if t.Peer != nil {
        return t.Peer.PublicKey(), make(map[string]primitives.PublicKey)
    }else if t.PublicKey != nil {
        return *t.PublicKey, make(map[string]primitives.PublicKey)
    }else{
        var (
            privateKey, publicKey []byte
            nodeLeftPublicKeys, nodeRightPublicKeys map[string]primitives.PublicKey
        )

        if t.Left.IsAlice() {
            privateKey, nodeLeftPublicKeys = t.Left.DiffieHellman()
            publicKey, nodeRightPublicKeys = t.Right.DiffieHellman()
        }else{
            privateKey, nodeLeftPublicKeys = t.Right.DiffieHellman()
            publicKey, nodeRightPublicKeys = t.Left.DiffieHellman()
        }

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

func (t TreeNode[IPeer]) Count() int {
    if t.Peer != nil {
        return 1
    }else{
        return t.Left.Count() + t.Right.Count()
    }
}

func merge(m ...map[string]primitives.PublicKey) map[string]primitives.PublicKey {
    ans := make(map[string]primitives.PublicKey, 0)

    for _, c := range m {
        for k, v := range c {
            ans[k] = v
        }
    }
    return ans
}
