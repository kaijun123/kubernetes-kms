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
	"context"
	"net"
	"time"

	"google.golang.org/grpc"

	kmsapi "github.com/kaijun123/kubernetes-kms/apis/v2"
	"github.com/kaijun123/kubernetes-kms/pkg/util"
	"k8s.io/klog/v2"
)

// GRPCService is a grpc server that runs the kms v2 alpha1 API.
type GRPCService struct {
	addr       string
	timeout    time.Duration
	server     *grpc.Server
	kmsService util.Service
}

// Asserts that the GRPC implements kmsapi.KeyManagementServiceServer: ie Status, Encrypt, Decrypt
var _ kmsapi.KeyManagementServiceServer = (*GRPCService)(nil)

// NewGRPCService creates an instance of GRPCService.
func NewGRPCService(
	address string,
	timeout time.Duration,

	kmsService util.Service,
) *GRPCService {
	klog.V(4).InfoS("KMS plugin configured", "address", address, "timeout", timeout)

	return &GRPCService{
		addr:       address,
		timeout:    timeout,
		kmsService: kmsService,
	}
}

// ListenAndServe accepts incoming connections on a Unix socket. It is a blocking method.
// Returns non-nil error unless Close or Shutdown is called.
func (s *GRPCService) ListenAndServe() error {
	ln, err := net.Listen("unix", s.addr)
	if err != nil {
		return err
	}
	defer ln.Close()

	gs := grpc.NewServer(
		grpc.ConnectionTimeout(s.timeout),
	)
	s.server = gs

	kmsapi.RegisterKeyManagementServiceServer(gs, s)

	klog.V(4).InfoS("kms plugin serving", "address", s.addr)
	return gs.Serve(ln)
}

// Shutdown performs a graceful shutdown. Doesn't accept new connections and
// blocks until all pending RPCs are finished.
func (s *GRPCService) Shutdown() {
	klog.V(4).InfoS("kms plugin shutdown", "address", s.addr)
	if s.server != nil {
		s.server.GracefulStop()
	}
}

// Close stops the server by closing all connections immediately and cancels
// all active RPCs.
func (s *GRPCService) Close() {
	klog.V(4).InfoS("kms plugin close", "address", s.addr)
	if s.server != nil {
		s.server.Stop()
	}
}

// Status sends a status request to specified kms service.
func (s *GRPCService) Status(ctx context.Context, _ *kmsapi.StatusRequest) (*kmsapi.StatusResponse, error) {
	res, err := s.kmsService.Status(ctx)
	if err != nil {
		return nil, err
	}

	return &kmsapi.StatusResponse{
		Version: res.Version,
		Healthz: res.Healthz,
		KeyId:   res.KeyId,
	}, nil
}

// Decrypt sends a decryption request to specified kms service.
func (s *GRPCService) Decrypt(ctx context.Context, req *kmsapi.DecryptRequest) (*kmsapi.DecryptResponse, error) {
	klog.V(4).InfoS("decrypt request received", "id", req.Uid)

	plaintext, err := s.kmsService.Decrypt(ctx, req.Uid, &util.DecryptRequestBody{
		KeyId:      req.KeyId,
		Ciphertext: req.Ciphertext,
	})
	if err != nil {
		return nil, err
	}

	return &kmsapi.DecryptResponse{
		Plaintext: plaintext,
	}, nil
}

// Encrypt sends an encryption request to specified kms service.
func (s *GRPCService) Encrypt(ctx context.Context, req *kmsapi.EncryptRequest) (*kmsapi.EncryptResponse, error) {
	klog.V(4).InfoS("encrypt request received", "id", req.Uid)

	encryptResponseBody, err := s.kmsService.Encrypt(ctx, req.Uid, req.Plaintext)
	if err != nil {
		return nil, err
	}

	return &kmsapi.EncryptResponse{
		Ciphertext:  encryptResponseBody.Ciphertext,
		KeyId:       encryptResponseBody.KeyId,
		Annotations: encryptResponseBody.Annotations,
	}, nil
}
