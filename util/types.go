package util

import "context"

// Service allows encrypting and decrypting data using an external Key Management Service.
type Service interface {
	// Encrypt bytes to a ciphertext.
	Encrypt(ctx context.Context, req *EncryptRequestBody) (*EncryptResponseBody, error)
	// Decrypt a given bytearray to obtain the original data as bytes.
	Decrypt(ctx context.Context, req *DecryptRequestBody) ([]byte, error)
	// Status returns the status of the KMS.
	Status(ctx context.Context) (*StatusResponseBody, error)
}

type EncryptRequestBody struct {
	KeyId     string `json:"keyId"`
	Plaintext []byte `json:"plaintext"`
}

type DecryptRequestBody struct {
	KeyId      string `json:"keyId"`
	Ciphertext []byte `json:"plaintext"`
}

type InitResponse struct {
	KeyId string `json:"keyId"`
}

// EncryptResponse is the response from the Envelope service when encrypting data.
type EncryptResponseBody struct {
	Ciphertext  []byte
	KeyId       string
	Annotations map[string][]byte
}

// StatusResponse is the response from the Envelope service when getting the status of the service.
type StatusResponseBody struct {
	Version string
	Healthz string
	KeyId   string
}
