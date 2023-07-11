package http

import (
	// "log"

	"bytes"
	"encoding/json"

	// "fmt"
	"io/ioutil"
	"log"
	"net/http"

	"k8s.io/kms/util"
)

const (
	encryptUrl = "http://localhost:8080/encrypt"
	decryptUrl = "http://localhost:8080/decrypt"
	initUrl    = "http://localhost:8080/init"
)

type HTTPClient struct {
	encryptUrl string
	decryptUrl string
	initUrl    string
}

// call encrypt url on the on-premise server
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
		log.Fatalf("error: %s", err)
		return nil, err
	}

	// Retrieves data from reponse
	encryptBody, encryptBodyErr := ioutil.ReadAll(resp.Body)
	if encryptBodyErr != nil {
		log.Fatalf("error: %s", encryptBodyErr)
		return nil, encryptBodyErr
	}

	var data map[string][]byte
	unmarshalErr := json.Unmarshal(encryptBody, &data)
	if unmarshalErr != nil {
		log.Fatalf("error: %s", unmarshalErr)
		return nil, unmarshalErr
	}

	ciphertext := data["ciphertext"]

	// fmt.Println("ciphertext: ", ciphertext)
	return ciphertext, nil
}

// call decrypt url on the on-premise serve
func (c *HTTPClient) Decrypt(keyId string, ciphertext []byte) ([]byte, error) {
	// Create the request body
	requestBody := util.DecryptRequestBody{
		KeyId:      keyId,
		Ciphertext: ciphertext,
	}

	// Marshal the request body into JSON
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	// Perform the HTTP POST request with the JSON request body
	resp, err := http.Post(c.decryptUrl, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Fatalf("error: %s", err)
		return nil, err
	}

	// Retrieves data from reponse
	decryptBody, decryptBodyErr := ioutil.ReadAll(resp.Body)
	if decryptBodyErr != nil {
		log.Fatalf("error: %s", decryptBodyErr)
		return nil, decryptBodyErr
	}

	var data map[string][]byte
	unmarshalErr := json.Unmarshal(decryptBody, &data)
	if unmarshalErr != nil {
		log.Fatalf("error: %s", unmarshalErr)
		return nil, unmarshalErr
	}

	plaintext := data["plaintext"]

	// fmt.Println("plaintext: ", plaintext)
	return plaintext, nil
}

// call Connect url on the on-premise server
func (c *HTTPClient) Init() (*http.Response, error) {
	resp, err := http.Get(c.initUrl)
	if err != nil {
		log.Fatalf("error during initialising the http client: %s", err)
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
