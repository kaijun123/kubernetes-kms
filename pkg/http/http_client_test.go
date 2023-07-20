package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	httpClient = NewHTTPClient()
	plaintext  = []byte("plaintext")
	keyId      = ""
	ciphertext = []byte{}
)

// call get-key-id url on the on-premise server
func TestMain(t *testing.T) {
	resp, err1 := httpClient.Init()
	assert.Equal(t, resp.StatusCode, 200)
	assert.Equal(t, err1, nil)

	body, err2 := ioutil.ReadAll(resp.Body)
	assert.Equal(t, err2, nil)

	var data map[string]string
	err3 := json.Unmarshal(body, &data)
	assert.Equal(t, err3, nil)

	keyId = data["key_id"]
	fmt.Printf("KeyId: %s\n", keyId)
}

// call encrypt url on the on-premise server
func TestEncrypt(t *testing.T) {
	res, err := httpClient.Encrypt(keyId, plaintext)
	assert.Equal(t, err, nil)
	ciphertext = res
}

// call decrypt url on the on-premise server
func TestDecrypt(t *testing.T) {
	res, err := httpClient.Decrypt(keyId, ciphertext)
	assert.Equal(t, err, nil)
	newPlaintext := res
	assert.Equal(t, newPlaintext, plaintext, "did not obtain back the inital plaintext")
}
