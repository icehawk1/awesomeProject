package main

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

