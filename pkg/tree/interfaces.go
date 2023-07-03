package tree

import "github.com/boke0/att/pkg/primitives"

type IPeer interface {
    IsAlice() bool
    PublicKey() primitives.PublicKey
    PrivateKey() *primitives.PrivateKey
    Id() string
}

