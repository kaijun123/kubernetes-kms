package util

import (
	"crypto/rand"
	"encoding/base64"
)

// Generates random strings for cryptographic purposes.
// Specify the number of bytes of the string using the length parameter
func GenerateRandomString(length int) string {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(b)
}
