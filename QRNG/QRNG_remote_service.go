package qrng

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"k8s.io/kms/http"
	"k8s.io/kms/util"
)

var _ util.Service = (*qrngRemoteService)(nil)

// No need for transformer
type qrngRemoteService struct {
	keyId      string
	httpClient *http.HTTPClient // Used for calling the apis on the on-premise server
}

func (s *qrngRemoteService) Encrypt(ctx context.Context, req *util.EncryptRequestBody) (*util.EncryptResponseBody, error) {
	ciphertext, err := s.httpClient.Encrypt(req.KeyId, req.Plaintext)
	if err != nil {
		log.Fatal("error: ", err)
		return nil, err
	}

	return &util.EncryptResponseBody{
		KeyId:      s.keyId,
		Ciphertext: ciphertext,
		Annotations: map[string][]byte{
			"mockAnnotationKey": []byte("1"),
		},
	}, nil
}

func (s *qrngRemoteService) Decrypt(ctx context.Context, req *util.DecryptRequestBody) ([]byte, error) {
	plaintext, err := s.httpClient.Decrypt(req.KeyId, req.Ciphertext)
	if err != nil {
		log.Fatal("error: ", err)
		return nil, err
	}
	return plaintext, nil
}

// Status returns the health status of the KMS plugin.
// We perform a simple encrypt/decrypt operation to verify the plugin's connectivity with On-Premise server.
func (s *qrngRemoteService) Status(ctx context.Context) (*util.StatusResponseBody, error) {
	plaintext := util.GenerateRandomString(32)
	fmt.Println("plaintext: ", plaintext)

	ciphertext, err1 := s.httpClient.Encrypt(s.keyId, []byte(plaintext))
	if err1 != nil {
		log.Fatal("error: ", err1)
		return nil, err1
	}
	newPlaintext, err2 := s.httpClient.Decrypt(s.keyId, ciphertext)
	if err2 != nil {
		log.Fatal("error: ", err2)
		return nil, err2
	}
	if string(newPlaintext) != plaintext {
		log.Fatal("error: ", err2)
		return nil, err2
	}
	fmt.Println("newPlaintext: ", string(newPlaintext))

	resp := &util.StatusResponseBody{
		Version: "v2beta1",
		Healthz: "ok",
		KeyId:   s.keyId,
	}
	return resp, nil
}

// NewQrngRemoteService creates an instance of qrngRemoteService.
// When creating a new instance of qrngRemoteService, you need to obtain a new keyId from the qrng
func NewQrngRemoteService() (*qrngRemoteService, error) {
	httpClient := http.NewHTTPClient()
	res, _ := httpClient.Init()
	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("error while retrieving response body: %s", responseBody)
		return nil, err
	}
	var initResponse util.InitResponse
	json.Unmarshal(responseBody, &initResponse)
	fmt.Println(initResponse.KeyId)

	qRemoteService := &qrngRemoteService{
		keyId:      initResponse.KeyId,
		httpClient: httpClient,
	}
	return qRemoteService, nil
}
