#!/bin/bash
# Rotate the KEK

# Create a new plugin
echo "apiVersion: v1
kind: Pod
metadata:
  name: mock-kmsv2-provider-2
  namespace: kube-system
  labels:
    tier: control-plane
    component: mock-kmsv2-provider
spec:
  hostNetwork: true
  containers:
    - name: mock-kmsv2-provider
      image: kaijun123/kubernetes-kms:v9
      imagePullPolicy: IfNotPresent
      livenessProbe:
        httpGet:
          path: /healthz
          port: 8087
        failureThreshold: 2
        periodSeconds: 10
      volumeMounts:
        - name: sock
          mountPath: /tmp
  volumes:
    - name: sock
      hostPath:
        path: /tmp
" > /etc/kubernetes/manifests/plugin2.yaml

# Edit the encryption config
echo "apiVersion: apiserver.config.k8s.io/v1
kind: EncryptionConfiguration
resources:
  - resources:
      - secrets
    providers:
      - kms:
          apiVersion: v2
          name: mock-kmsv2-provider-2
          endpoint: unix:///tmp/kms.socket
          timeout: 3s
      - kms:
          apiVersion: v2
          name: mock-kmsv2-provider
          endpoint: unix:///tmp/kms.socket
          timeout: 3s
" > /etc/kubernetes/encryption/encryption-conf.yaml