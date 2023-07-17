package qrng

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/kaijun123/kubernetes-kms/pkg/http"
	"github.com/kaijun123/kubernetes-kms/pkg/util"
)

var _ util.Service = (*qrngRemoteService)(nil)

// No need for transformer
type qrngRemoteService struct {
	keyId      string
	httpClient *http.HTTPClient // Used for calling the apis on the on-premise server
}

// Calls the `Encrypt()` method of the httpClient
func (s *qrngRemoteService) Encrypt(ctx context.Context, uid string, plaintext []byte) (*util.EncryptResponseBody, error) {
	fmt.Println("Calling Encrypt()......")

	ciphertext, err := s.httpClient.Encrypt(s.keyId, plaintext)
	if err != nil {
		return nil, err
	}

	fmt.Println("End of Encrypt()......")
	return &util.EncryptResponseBody{
		KeyId:      s.keyId,
		Ciphertext: ciphertext,
		Annotations: map[string][]byte{
			"mockAnnotationKey": []byte("1"),
		},
	}, nil
}

// Calls the `Decrypt()` method of the httpClient
func (s *qrngRemoteService) Decrypt(ctx context.Context, uid string, req *util.DecryptRequestBody) ([]byte, error) {
	fmt.Println("Calling Decrypt()......")

	if req.KeyId != s.keyId {
		return nil, errors.New("invalid keyID")
	}

	plaintext, err := s.httpClient.Decrypt(req.KeyId, req.Ciphertext)
	if err != nil {
		return nil, err
	}

	fmt.Println("End of Decrypt()......")
	return plaintext, nil
}

// Status returns the api_version, health_status and key_id of the KMS plugin.
// The API server considers the key_id returned from the Status procedure call to be authoritative.
// If an EncryptRequest procedure call returns a key_id that is different from Status, the response is thrown away and the plugin is considered unhealthy.
// In this methodm, we perform a simple encrypt/decrypt operation to verify the plugin's connectivity with On-Premise server.
func (s *qrngRemoteService) Status(ctx context.Context) (*util.StatusResponseBody, error) {
	fmt.Println("Calling Status()......")

	plaintext := util.GenerateRandomString(32)
	// fmt.Println("plaintext: ", plaintext)

	ciphertext, err := s.httpClient.Encrypt(s.keyId, []byte(plaintext))
	if err != nil {
		return nil, err
	}
	newPlaintext, err := s.httpClient.Decrypt(s.keyId, ciphertext)
	if err != nil {
		return nil, err
	}
	// fmt.Println("newPlaintext: ", string(newPlaintext))

	if string(newPlaintext) != plaintext {
		err := errors.New("error: Plaintext obtained after decryption is no the same as the plaintext before encryption")
		return nil, err
	}

	resp := &util.StatusResponseBody{
		Version: "v2beta1",
		Healthz: "ok",
		KeyId:   s.keyId,
	}
	fmt.Println("End of Status()......")
	return resp, nil
}

// NewQrngRemoteService creates an instance of qrngRemoteService.
// When creating a new instance of qrngRemoteService, you need to obtain a new keyId from the qrng.
func NewQrngRemoteService() (*qrngRemoteService, error) {
	fmt.Println("Calling NewQrngRemoteService()......")
	// fmt.Println("before creating http client pointer")

	httpClient := http.NewHTTPClient()
	// fmt.Println("after creating http client pointer")
	fmt.Println(httpClient)

	// fmt.Println("before calling Init()")
	res, err := httpClient.Init()
	if err != nil {
		log.Fatalf("error during init: %s", err)
		return nil, err
	}
	// fmt.Println("after calling Init()")

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("error while retrieving response body: %s", responseBody)
		return nil, err
	}
	var initResponse util.InitResponse
	json.Unmarshal(responseBody, &initResponse)
	// fmt.Println("KeyId: ", initResponse.KeyId)

	qRemoteService := &qrngRemoteService{
		keyId:      initResponse.KeyId,
		httpClient: httpClient,
	}
	fmt.Println("End of NewQrngRemoteService()......")
	return qRemoteService, nil
}
