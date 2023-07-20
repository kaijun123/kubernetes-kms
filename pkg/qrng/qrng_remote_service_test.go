package qrng

import (
	"context"

	"testing"

	"github.com/kaijun123/kubernetes-kms/pkg/util"
	"github.com/stretchr/testify/assert"
)

const (
	version       = "v2beta1"
	testPlaintext = "lorem ipsum dolor sit amet"
)

func testContext(t *testing.T) context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	return ctx
}

func TestNewQrngRemoteService(t *testing.T) {
	ctx := testContext(t)

	plaintext := []byte(testPlaintext)

	qrngRemoteService, err := NewQrngRemoteService()
	assert.Equal(t, err, nil, "qrngRemoteService initialisation error")

	t.Run("should be able to encrypt and decrypt", func(t *testing.T) {
		encryptionResponseBody, err := qrngRemoteService.Encrypt(ctx, "", plaintext)
		assert.Equal(t, err, nil)
		assert.Equal(t, qrngRemoteService.KeyId, encryptionResponseBody.KeyId, "keyId should always be the same")
		assert.NotEqual(t, plaintext, encryptionResponseBody.Ciphertext, "plaintext and ciphertext cannot be the same")

		newPlaintext, err := qrngRemoteService.Decrypt(ctx, "", &util.DecryptRequestBody{
			KeyId:      qrngRemoteService.KeyId,
			Ciphertext: encryptionResponseBody.Ciphertext,
		})
		assert.Equal(t, err, nil)
		assert.Equal(t, newPlaintext, plaintext, "did not obtain back the inital plaintext")
	})

	t.Run("should return error when decrypt with an invalid keyId", func(t *testing.T) {
		encryptionResponseBody, err := qrngRemoteService.Encrypt(ctx, "", plaintext)
		assert.Equal(t, err, nil)
		assert.Equal(t, qrngRemoteService.KeyId, encryptionResponseBody.KeyId, "keyId should always be the same")
		assert.NotEqual(t, plaintext, encryptionResponseBody.Ciphertext, "plaintext and ciphertext cannot be the same")

		_, err = qrngRemoteService.Decrypt(ctx, "", &util.DecryptRequestBody{
			KeyId:      encryptionResponseBody.KeyId + "1",
			Ciphertext: encryptionResponseBody.Ciphertext,
		})

		if err.Error() != "invalid keyID" {
			t.Errorf("should have returned an invalid keyID error. Got %v, requested keyID: %q, remote service keyID: %q", err, encryptionResponseBody.KeyId+"1", qrngRemoteService.KeyId)
		}
	})

	t.Run("should return status data", func(t *testing.T) {
		status, err := qrngRemoteService.Status(ctx)
		assert.Equal(t, err, nil)

		if status.Healthz != "ok" {
			t.Errorf("want: %q, have: %q", "ok", status.Healthz)
		}
		if len(status.KeyId) == 0 {
			t.Errorf("want: len(keyID) > 0, have: %d", len(status.KeyId))
		}
		if status.Version != version {
			t.Errorf("want %q, have: %q", version, status.Version)
		}
	})
}
