package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateRandomString(t *testing.T) {
	string1 := GenerateRandomString(10)
	// log.Printf("%s", string1)
	string2 := GenerateRandomString(10)
	// log.Printf("%s", string2)

	assert.NotEqual(t, string1, string2)
}
