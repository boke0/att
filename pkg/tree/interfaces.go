package tree

type IPeer interface {
    IsAlice() bool
    PublicKey() []byte
    PrivateKey() *[]byte
    Id() string
}

