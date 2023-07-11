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
func TestInit(t *testing.T) {
	resp, err1 := httpClient.Init()
	assert.Equal(t, resp.StatusCode, 200)
	assert.Equal(t, err1, nil)

	body, err2 := ioutil.ReadAll(resp.Body)
	assert.Equal(t, err2, nil)

	var data map[string]string
	err3 := json.Unmarshal(body, &data)
	assert.Equal(t, err3, nil)

	keyId = data["keyId"]
	fmt.Printf("KeyId: %s\n", keyId)
}

// call encrypt url on the on-premise server
func TestEncrypt(t *testing.T) {
	res, err := httpClient.Encrypt(keyId, plaintext)
	assert.Equal(t, err, nil)
	ciphertext = res

	// encryptResp, encryptRespErr := httpClient.Encrypt(keyId, plaintext)
	// assert.Equal(t, encryptResp.StatusCode, 200)
	// assert.Equal(t, encryptRespErr, nil)

	// encryptBody, encryptBodyErr := ioutil.ReadAll(encryptResp.Body)
	// assert.Equal(t, encryptBodyErr, nil)

	// var data2 map[string][]byte
	// unmarshalErr2 := json.Unmarshal(encryptBody, &data2)
	// assert.Equal(t, unmarshalErr2, nil)

	// ciphertext = data2["ciphertext"]

	// fmt.Printf("ciphertext: %v\n", string(ciphertext))
}

// call decrypt url on the on-premise server
func TestDecrypt(t *testing.T) {
	res, err := httpClient.Decrypt(keyId, ciphertext)
	assert.Equal(t, err, nil)
	newPlaintext := res
	assert.Equal(t, newPlaintext, plaintext, "did not obtain back the inital plaintext")

	// decryptResp, decryptRespErr := httpClient.Decrypt(keyId, ciphertext)
	// assert.Equal(t, decryptResp.StatusCode, 200)
	// assert.Equal(t, decryptRespErr, nil)

	// decryptBody, decryptBodyErr := ioutil.ReadAll(decryptResp.Body)
	// assert.Equal(t, decryptBodyErr, nil)

	// var data3 map[string][]byte
	// unmarshalErr3 := json.Unmarshal(decryptBody, &data3)
	// assert.Equal(t, unmarshalErr3, nil)

	// newPlaintext := data3["plaintext"]
	// assert.Equal(t, newPlaintext, plaintext)

	// fmt.Println(plaintext)
}
