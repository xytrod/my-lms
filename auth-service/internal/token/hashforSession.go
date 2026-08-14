package token

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashForSession(secret string) string {
	res := sha256.Sum256([]byte(secret))
	return hex.EncodeToString(res[:])
}
