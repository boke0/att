package primitives

import (
	"crypto/rand"
	"crypto/ed25519"

	"golang.org/x/crypto/curve25519"
	"filippo.io/edwards25519"
	"filippo.io/edwards25519/field"
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


func ToEd25519(p []byte) (ed25519.PublicKey, error) {
	A, err := convertMont(p)
	if err != nil {
		return nil, err
	}
	return A.Bytes(), nil
}


// calculateKeyPair converts a Montgomery private key k to a twisted Edwards
// public key and private key (A, a) as defined in
// https://signal.org/docs/specifications/xeddsa/#elliptic-curve-conversions
//
//   calculate_key_pair(k):
//       E = kB
//       A.y = E.y
//       A.s = 0
//       if E.s == 1:
//           a = -k (mod q)
//       else:
//           a = k (mod q)
//       return A, a
func calculateKeyPair(p []byte) ([]byte, *edwards25519.Scalar, error) {
	var pA edwards25519.Point
	var sa edwards25519.Scalar

	k, err := (&edwards25519.Scalar{}).SetBytesWithClamping(p)
	if err != nil {
		return nil, nil, err
	}

	pub := pA.ScalarBaseMult(k).Bytes()
	signBit := (pub[31] & 0x80) >> 7

	if signBit == 1 {
		sa.Negate(k)
		// Set sig bit to 0
		pub[31] &= 0x7F
	} else {
		sa.Set(k)
	}

	return pub, &sa, nil
}

var one = (&field.Element{}).One()

// convertMont converts from a Montgomery u-coordinate to a twisted Edwards
// point P, according to
// https://signal.org/docs/specifications/xeddsa/#elliptic-curve-conversions
//
//   convert_mont(u):
//     umasked = u (mod 2|p|)
//     P.y = u_to_y(umasked)
//     P.s = 0
//     return P
func convertMont(u []byte) (*edwards25519.Point, error) {
	um, err := (&field.Element{}).SetBytes(u)
	if err != nil {
		return nil, err
	}

	// y = (u - 1)/(u + 1)
	a := new(field.Element).Subtract(um, one)
	b := new(field.Element).Add(um, one)
	y := new(field.Element).Multiply(a, b.Invert(b)).Bytes()

	// Set sign to 0
	y[31] &= 0x7F

	return (&edwards25519.Point{}).SetBytes(y)
}
