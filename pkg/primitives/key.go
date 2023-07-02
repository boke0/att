package primitives

import (
	"crypto/rand"

	"golang.org/x/crypto/curve25519"
)

func AsPublic(key []byte) []byte {
	pub, _ := curve25519.X25519(key, curve25519.Basepoint)
	return pub
}


func DiffieHellman(priv []byte, pub  []byte) []byte {
	result, _ := curve25519.X25519(priv, pub)
	return result
}

func RandomByte() []byte {
	b := make([]byte, curve25519.ScalarSize)
	rand.Read(b)
	return b
}
