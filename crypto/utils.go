package crypto

import (
	"bytes"
	"crypto/sha1"
)

func Hash(in []byte) []byte {
	hash := sha1.Sum(in)
	return hash[:]
}

func HashVerify(in []byte, hash []byte) bool {
	h := Hash(in)
	return bytes.Equal(h, hash)
}
