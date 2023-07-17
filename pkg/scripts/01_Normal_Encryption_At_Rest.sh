#!/bin/sh

cd /etc/kubernetes

sudo mkdir encryption

sudo cp /mount/pkg/plugins/encryption-conf.yaml /etc/kubernetes/encryption/encryption-conf.yaml

sudo cp /mount/pkg/plugins/kube-apiserver.yaml /etc/kubernetes/manifests/kube-apiserver.yaml