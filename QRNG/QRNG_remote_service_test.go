package qrng

import (
	"context"
	// "fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/kms/util"
)

var (
	qRemoteService *qrngRemoteService
	plaintext      = []byte("plaintext")
	ciphertext     = []byte("")
)

func testContext(t *testing.T) context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	return ctx
}

func TestNewQrngRemoteService(t *testing.T) {
	qrngRemoteService, err := NewQrngRemoteService()
	assert.Equal(t, err, nil, "qrngRemoteService initialisation error")
	qRemoteService = qrngRemoteService
}

func TestEncrypt(t *testing.T) {
	ctx := testContext(t)
	encryptRequestBody := &util.EncryptRequestBody{
		KeyId:     qRemoteService.keyId,
		Plaintext: plaintext,
	}
	encryptResponseBody, err := qRemoteService.Encrypt(ctx, encryptRequestBody)
	assert.Equal(t, err, nil)

	ciphertext = encryptResponseBody.Ciphertext
	// fmt.Println("keyId: ", encryptResponseBody.KeyId)
	// fmt.Println("ciphertext: ", ciphertext)
	// fmt.Println("Annotations: ", encryptResponseBody.Annotations)
}

func TestDecrypt(t *testing.T) {
	ctx := testContext(t)
	decryptRequestBody := &util.DecryptRequestBody{
		KeyId:      qRemoteService.keyId,
		Ciphertext: ciphertext,
	}
	newPlaintext, err := qRemoteService.Decrypt(ctx, decryptRequestBody)
	assert.Equal(t, err, nil)
	assert.Equal(t, newPlaintext, plaintext, "did not obtain back the inital plaintext")
	// fmt.Println("newPlaintext: ", newPlaintext)
}

func TestStatus(t *testing.T) {
	ctx := testContext(t)
	status, err := qRemoteService.Status(ctx)
	assert.Equal(t, err, nil)
	assert.Equal(t, status.Healthz, "ok")
}
