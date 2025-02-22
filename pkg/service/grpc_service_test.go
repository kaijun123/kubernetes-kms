/*
Copyright 2023 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package service

import (
	"bytes"
	"context"
	"encoding/base64"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"k8s.io/apimachinery/pkg/util/wait"

	kmsapi "github.com/kaijun123/kubernetes-kms/apis/v2"
	"github.com/kaijun123/kubernetes-kms/pkg/qrng"
	"github.com/kaijun123/kubernetes-kms/pkg/util"
)

const version = "v2beta1"

func TestGRPCService(t *testing.T) {
	t.Parallel()

	defaultTimeout := 30 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	t.Cleanup(cancel)

	address := filepath.Join(os.TempDir(), "kmsv2.sock")
	plaintext := []byte("lorem ipsum dolor sit amet")
	r := rand.New(rand.NewSource(time.Now().Unix()))
	id, err := makeID(r.Read)
	if err != nil {
		t.Fatal(err)
	}

	// // Start the gRPC server
	// kmsService := newBase64Service(id)
	// server := NewGRPCService(address, defaultTimeout, kmsService)

	remoteKMSService, err := qrng.NewQrngRemoteService()
	if err != nil {
		t.Fatal(err)
	}
	server := NewGRPCService(address, defaultTimeout, remoteKMSService)
	go func() {
		if err := server.ListenAndServe(); err != nil {
			panic(err)
		}
	}()
	t.Cleanup(server.Shutdown)

	// Start the gRPC client
	client := newClient(t, address)

	// make sure the gRPC server is up before running tests
	if err := wait.PollImmediateUntilWithContext(ctx, time.Second, func(ctx context.Context) (bool, error) {
		ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()

		_, err := client.Status(ctx, &kmsapi.StatusRequest{})
		if err != nil {
			t.Logf("failed to get kms status: %v", err)
		}

		return err == nil, nil
	}); err != nil {
		t.Fatal(err)
	}

	t.Run("should be able to encrypt and decrypt through unix domain sockets", func(t *testing.T) {
		t.Parallel()

		encRes, err := client.Encrypt(ctx, &kmsapi.EncryptRequest{
			Plaintext: plaintext,
			Uid:       id,
		})
		if err != nil {
			t.Fatal(err)
		}

		if bytes.Equal(plaintext, encRes.Ciphertext) {
			t.Fatal("plaintext and ciphertext shouldn't be equal!")
		}

		decRes, err := client.Decrypt(ctx, &kmsapi.DecryptRequest{
			Ciphertext:  encRes.Ciphertext,
			KeyId:       encRes.KeyId,
			Annotations: encRes.Annotations,
			Uid:         id,
		})
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(decRes.Plaintext, plaintext) {
			t.Errorf("want: %q, have: %q", plaintext, decRes.Plaintext)
		}
	})

	t.Run("should return status data", func(t *testing.T) {
		t.Parallel()

		status, err := client.Status(ctx, &kmsapi.StatusRequest{})
		if err != nil {
			t.Fatal(err)
		}

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

func newClient(t *testing.T, address string) kmsapi.KeyManagementServiceClient {
	t.Helper()

	cnn, err := grpc.Dial(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDialer(func(addr string, t time.Duration) (net.Conn, error) {
			return net.Dial("unix", addr)
		}),
	)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { _ = cnn.Close() })

	return kmsapi.NewKeyManagementServiceClient(cnn)
}

type testService struct {
	decrypt func(ctx context.Context, uid string, req *util.DecryptRequestBody) ([]byte, error)
	encrypt func(ctx context.Context, uid string, data []byte) (*util.EncryptResponseBody, error)
	status  func(ctx context.Context) (*util.StatusResponseBody, error)
}

var _ util.Service = (*testService)(nil)

func (s *testService) Decrypt(ctx context.Context, uid string, req *util.DecryptRequestBody) ([]byte, error) {
	return s.decrypt(ctx, uid, req)
}

func (s *testService) Encrypt(ctx context.Context, uid string, data []byte) (*util.EncryptResponseBody, error) {
	return s.encrypt(ctx, uid, data)
}

func (s *testService) Status(ctx context.Context) (*util.StatusResponseBody, error) {
	return s.status(ctx)
}

// Creates a random keyId
func makeID(rand func([]byte) (int, error)) (string, error) {
	b := make([]byte, 12)
	if _, err := rand(b); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(b), nil
}

// // Used to a mock KMS service for testing purposes
// // The mock KMS service is used to create a the gRPC server
// func newBase64Service(keyId string) *testService {
// 	decrypt := func(_ context.Context, _ string, req *util.DecryptRequestBody) ([]byte, error) {
// 		if req.KeyId != keyId {
// 			return nil, fmt.Errorf("keyId mismatch. want: %q, have: %q", keyId, req.KeyId)
// 		}

// 		return base64.StdEncoding.DecodeString(string(req.Ciphertext))
// 	}

// 	encrypt := func(_ context.Context, _ string, data []byte) (*util.EncryptResponseBody, error) {
// 		return &util.EncryptResponseBody{
// 			Ciphertext: []byte(base64.StdEncoding.EncodeToString(data)),
// 			KeyId:      keyId,
// 		}, nil
// 	}

// 	status := func(_ context.Context) (*util.StatusResponseBody, error) {
// 		return &util.StatusResponseBody{
// 			Version: version,
// 			Healthz: "ok",
// 			KeyId:   keyId,
// 		}, nil
// 	}

// 	return &testService{
// 		decrypt: decrypt,
// 		encrypt: encrypt,
// 		status:  status,
// 	}
// }
