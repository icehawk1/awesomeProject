package blockchain

import (
	"crypto/sha256"
	"fmt"
)

func computeSha256Hex(input string) string {
	return fmt.Sprintf("%X", computeSha256(input))
}

func computeSha256(input string) []byte {
	hash := sha256.New()
	hash.Write([]byte(input))
	return hash.Sum(nil)
}

type Hashable interface {
	ComputeHash() string
	ComputeHashByte() string
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}