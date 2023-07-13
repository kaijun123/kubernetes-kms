package http

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/kaijun123/kubernetes-kms/pkg/util"
)

const (
	encryptUrl = "http://host.docker.internal:8080/encrypt"
	decryptUrl = "http://host.docker.internal:8080/decrypt"
	initUrl    = "http://host.docker.internal:8080/init"
)

type HTTPClient struct {
	encryptUrl string
	decryptUrl string
	initUrl    string
}

// call Encrypt api on the on-premise server
func (c *HTTPClient) Encrypt(keyId string, plaintext []byte) ([]byte, error) {
	// Create the request body
	requestBody := util.EncryptRequestBody{
		KeyId:     keyId,
		Plaintext: plaintext,
	}

	// Marshal the request body into JSON
	jsonBody, marshalErr := json.Marshal(requestBody)
	if marshalErr != nil {
		return nil, marshalErr
	}

	// Perform the HTTP POST request with the JSON request body
	resp, err := http.Post(c.encryptUrl, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	// Retrieves data from reponse
	encryptBody, encryptBodyErr := ioutil.ReadAll(resp.Body)
	if encryptBodyErr != nil {
		return nil, encryptBodyErr
	}

	var data map[string][]byte
	unmarshalErr := json.Unmarshal(encryptBody, &data)
	if unmarshalErr != nil {
		return nil, unmarshalErr
	}

	ciphertext := data["ciphertext"]

	// fmt.Println("ciphertext: ", ciphertext)
	return ciphertext, nil
}

// call Decrypt api on the on-premise serve
func (c *HTTPClient) Decrypt(keyId string, ciphertext []byte) ([]byte, error) {
	// Create the request body
	requestBody := util.DecryptRequestBody{
		KeyId:      keyId,
		Ciphertext: ciphertext,
	}

	// Marshal the request body into JSON
	jsonBody, marshalErr := json.Marshal(requestBody)
	if marshalErr != nil {
		return nil, marshalErr
	}

	// Perform the HTTP POST request with the JSON request body
	resp, err := http.Post(c.decryptUrl, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	// Retrieves data from reponse
	decryptBody, decryptBodyErr := ioutil.ReadAll(resp.Body)
	if decryptBodyErr != nil {
		return nil, decryptBodyErr
	}

	var data map[string][]byte
	unmarshalErr := json.Unmarshal(decryptBody, &data)
	if unmarshalErr != nil {
		return nil, unmarshalErr
	}

	plaintext := data["plaintext"]

	// fmt.Println("plaintext: ", plaintext)
	return plaintext, nil
}

// call Init api on the on-premise server
func (c *HTTPClient) Init() (*http.Response, error) {
	resp, err := http.Get(c.initUrl)
	if err != nil {
		return nil, err
	}
	// fmt.Println("response: ", resp)
	return resp, nil
}

// To be called when creating a new qrngRemoteService. ie calling NewQrngRemoteService
func NewHTTPClient() *HTTPClient {
	return &HTTPClient{
		encryptUrl: encryptUrl,
		decryptUrl: decryptUrl,
		initUrl:    initUrl,
	}
}
