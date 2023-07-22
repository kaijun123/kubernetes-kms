#!/bin/bash
# Edits the config file to decrypt the secrets

echo "apiVersion: apiserver.config.k8s.io/v1
kind: EncryptionConfiguration
resources:
  - resources:
      - secrets
    providers:
      - identity: {}
      - kms:
          apiVersion: v2
          name: mock-kmsv2-provider
          endpoint: unix:///tmp/kms.socket
          timeout: 3s
" > /etc/kubernetes/encryption/encryption-conf.yaml