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

package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kaijun123/kubernetes-kms/pkg/healthz"
	"github.com/kaijun123/kubernetes-kms/pkg/qrng"
	"github.com/kaijun123/kubernetes-kms/pkg/service"
	"github.com/kaijun123/kubernetes-kms/pkg/util"
	"k8s.io/klog/v2"
)

var (
	listenAddr = flag.String("listen-addr", "unix:///tmp/kms.socket", "gRPC listen address")
	timeout    = flag.Duration("timeout", 5*time.Second, "gRPC timeout")
)

func main() {
	flag.Parse()

	addr, err := util.ParseEndpoint(*listenAddr)
	if err != nil {
		klog.ErrorS(err, "failed to parse endpoint")
		os.Exit(1)
	}

	remoteKMSService, err := qrng.NewQrngRemoteService()
	if err != nil {
		klog.ErrorS(err, "failed to create remote service")
		os.Exit(1)
	}

	ctx := withShutdownSignal(context.Background())

	grpcService := service.NewGRPCService(
		addr,
		*timeout,
		remoteKMSService,
	)

	klog.InfoS("starting server", "listenAddr", *listenAddr)
	go func() {
		if err := grpcService.ListenAndServe(); err != nil {
			klog.ErrorS(err, "failed to serve")
			os.Exit(1)
		}
	}()

	httpServer := healthz.InitHttpServer(remoteKMSService)
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			klog.ErrorS(err, "failed to serve")
			os.Exit(1)
		}
	}()

	<-ctx.Done()

	klog.InfoS("shutting down server")
	grpcService.Shutdown()
}

// withShutdownSignal returns a copy of the parent context that will close if
// the process receives termination signals.
func withShutdownSignal(ctx context.Context) context.Context {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)

	nctx, cancel := context.WithCancel(ctx)

	go func() {
		<-signalChan
		klog.InfoS("received shutdown signal")
		cancel()
	}()
	return nctx
}
