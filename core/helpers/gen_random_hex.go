package helpers

import (
	"crypto/rand"
	"encoding/hex"
)

func GenRandomHexString(length int) string {
	buf := make([]byte, (length+1)/2)

	rand.Read(buf)

	return hex.EncodeToString(buf)[:length]
}
