package primitives

import (
	"crypto/sha256"
	"encoding/hex"
)

func Hash(str string) string {
	bytes := sha256.Sum256([]byte(str))
	return hex.EncodeToString(bytes[:])
}
