#!/bin/sh

cd /etc/kubernetes

sudo mkdir encryption

sudo cp /mount/pkg/plugins/plugin.yaml /etc/kubernetes/manifests/plugin.yaml

sudo cp /mount/pkg/plugins/kms-encryption-conf.yaml /etc/kubernetes/encryption/encryption-conf.yaml

sudo cp /mount/pkg/plugins/kms-kube-apiserver.yaml /etc/kubernetes/manifests/kube-apiserver.yaml