## KMS v2beta1
This document contains explanation of how to make use of the KMS v2beta1 and an explanation of the files in the ```scripts``` directory

### Content Page
- [Directory Structure](#Directory-Strucuture)
- [KMS v2beta1 Resources](#KMS-v2beta1-Resources)
- [Requirement](#Requirements)
- [Kind Resources](#Kind-Resources)
- [Quick Start](#Quick-Start)
- [Explanation](#Explanation)
  * [Explanation of kind.yaml](#Explanation-of-kind.yaml)
  * [Explanation of plugin.yaml](#Explanation-of-plugin.yaml)
  * [Explanation of encryption-config.yaml](#Explanation-of-encryption-config.yaml)
- [Enabling KMS Features](#Enabling-KMS-Features)
  * [Key Rotation](#Key-Rotation)
  * [Decryption](#Decryption)

### Directory Structure
```
.
├── 01_Decrypt_KMS_Encryption.sh
├── 02_Rotate_Keys.sh
├── KMS.md
├── Kind.md
├── encryption-conf.yaml
├── kind.yaml
├── kube-apiserver.yaml
├── plugin.yaml
└── run-e2e.sh        # Script from Kubernetes source code https://github.com/kubernetes/kubernetes/tree/master/test/e2e/testing-manifests/auth/encrypt. For reference
```

### KMS v2beta1 Resources
- Encryption At Rest doc: https://kubernetes.io/docs/tasks/administer-cluster/encrypt-data/
- KMS Provider doc: https://kubernetes.io/docs/tasks/administer-cluster/kms-provider/#configuring-the-kms-provider-kms-v2

<p align="center">
    <img src="./images/Decrypt\ Request.png">
    <img src="./images/Encrypt\ Request.png">
    <img src="./images/Status\ Request.png">
    <img src="./images/Generate\ DEK.png">
</p>

### Requirements
- To use KMS v2beta1, Kubernetes v1.27
- Comments:
  - To check for the Kubernetes version, run ```kubectl version```
  - At the time of writing, the latest version for Minikube is Kubernetes v1.25 which is too old
  - Kind is the recommended Kubenetes tool to use. To install Kind on macOS, run ```brew install kind```

### Kind Resources
- Official doc: https://kind.sigs.k8s.io/docs/user/quick-start/#creating-a-cluster
- https://medium.com/@talhakhalid101/creating-a-kubernetes-cluster-for-development-with-kind-189df2cb0792
- https://blog.kubesimplify.com/getting-started-with-kind-creating-a-multi-node-local-kubernetes-cluster#heading-creating-multi-node-cluster

### Quick Start
```
<!-- Create the cluster based on config -->
kind create cluster --config ./pkg/scripts/kind.yaml

<!-- Check that the pods are running properly -->
kubectl get pods -A

<!-- Create secret -->
kubectl create secret generic my-secret --from-literal=key1=supersecret

<!-- Enter the etcd pod -->
kubectl exec -n kube-system -it etcd-kind-control-plane -- sh

<!-- In the etcd shell, view the saved secret in the etcd shell -->
ETCDCTL_API=3 etcdctl \
 --cacert=/etc/kubernetes/pki/etcd/ca.crt   \
 --cert=/etc/kubernetes/pki/etcd/server.crt \
 --key=/etc/kubernetes/pki/etcd/server.key  \
 get /registry/secrets/default/my-secret

<!-- Login to control panel shell -->
docker ps
docker exec -t <container-id> bash

<!-- Rotate keys -->
cd /mount
./02_Rotate_Keys.sh
kubectl get secrets --all-namespaces -o json | kubectl replace -f -

<!-- Decrypt secrets -->
cd /mount
./01_Decrypt_KMS_Encryption.sh
kubectl get secrets --all-namespaces -o json | kubectl replace -f -
```

### Explanation
#### Explanation of ```kind.yaml```
- Specify the node image to use to build the node. In this case, we specify the image for Kubernetes 1.27
  ```
  - role: control-plane
    image: kindest/node:v1.27.3@sha256:3966ac761ae0136263ffdb6cfd4db23ef8a83cba8a463690e98317add2c9ba72
  ```
- Mount files from the host OS to the control-plane
  - ```hostPath```: Specifies the file path on the host OS
  - ```containerPath```: Specifies the destination to mount the file in the node
  - ```readOnly```: If true, the file cannot be edited from within the node
  - To check if the mount is correct, you can enter node shell using ```docker exec -it <container-id>```
  ```
  extraMounts:
  - containerPath: /mount
    hostPath: /Users/kaijun/Go/src/kms/pkg/scripts/
    propagation: None
  - containerPath: /etc/kubernetes/encryption/encryption-conf.yaml
    hostPath: /Users/kaijun/Go/src/kms/pkg/scripts/encryption-conf.yaml
    # readOnly: true
    propagation: None
  - containerPath: /etc/kubernetes/manifests/plugin.yaml
    hostPath: /Users/kaijun/Go/src/kms/pkg/scripts/plugin.yaml
    # readOnly: true
    propagation: None
  ```
  - Here, there are 3 files mounted into the control-plane.
    1. Scripts directory: The directory is mounted so that the bash scripts can be run to make the config changes more convenient
    2. Encryption config yaml file: The config is needed in order for the apiserver to know how to encrypt the data
    3. Plugin yaml file: By mounting the plugin yaml into the ```/etc/kubernetes/manifests/``` directory, you are running the plugin as a static pod. A static pod is a pod that is always kept running; if it stops, the apiserver will bring it up again
  - Note: ```readOnly``` is commented out to make the config changes when testing out the KMS functionalities easier
- Make config changes to the apiserver
  ```
    kubeadmConfigPatches:
    - |
      kind: ClusterConfiguration
      apiServer:
        extraArgs:
          feature-gates: "KMSv2=true"
          encryption-provider-config: "/etc/kubernetes/enc/encryption-conf.yaml"
          encryption-provider-config-automatic-reload: "true"
          v: "5"
        extraVolumes:
        - name: enc-vol
          hostPath: "/etc/kubernetes/encryption/encryption-conf.yaml"
          mountPath: "/etc/kubernetes/enc/encryption-conf.yaml"
          readOnly: true
          pathType: File
        - name: sock-path
          hostPath: "/tmp"
          mountPath: "/tmp"
          type: DirectoryOrCreate
  ```
  - ```extraVolumes```: Allows the Pods that run on the Node to have access to certain files on the Node through Volume mounts.
    - ```hostPath```: Specifies the path on the Node
    - ```mountPath```: Specifies the destination to mount the file in the Pod
    - For the encryption config, we are mounting the config found at "/etc/kubernetes/encryption/encryption-conf.yaml" on the Node to "/etc/kubernetes/enc/encryption-conf.yaml" in the Pod
    - For the socket, we need to mount the unix socket which is on the Node into the Pod, so that the apiserver can make the gRPC calls to the unix socket
  - ```extraArgs```:
    - ```feature-gates: "KMSv2=true"```: Need to set the KMSv2 flag to true to enable KMSv2
    - ```encryption-provider-config```: Specify the path to which the apiserver will search for the encryption config file. The path must be the path of the file in the Pod
    - ```encryption-provider-config-automatic-reload```: When set as true, there is no need to restart the apiServer everytime there is a change in the encryption config (Eg Key Rotation, Decryption, etc)
- A sample of the apiserver config can be seen in ```./kube-apiserver.yaml```

#### Explanation of ```plugin.yaml```
- To be added to ```/etc/kubernetes/manifests``` to run the plugin as a static pods
- Specify the container name, image name and pull policy. Image is by default pulled from Dockerhub
  ```
  containers:
    - name: mock-kmsv2-provider
      image: kaijun123/kubernetes-kms:v9
      imagePullPolicy: IfNotPresent
  ```
- Mount the unix socket on the Node into the plugin so that the gRPC server can listen to the unix socket
  ```
        volumeMounts:
          - name: sock
            mountPath: /tmp
    volumes:
      - name: sock
        hostPath:
          path: /tmp
  ```

#### Explanation of ```encryption-config.yaml```
- Contains Encryption Config
- Specify the resource to encrypt
  ```
    - resources:
      - secrets
  ```
- Specify apiVersion. ```Name``` should correspond to that in ```plugin.yaml```. ```endpoint``` refers to the unix socket endpoint to listen to; should correspond to the socket listened to in the code and the ```plugin.yaml```
  ```
  providers:
      - kms:
          apiVersion: v2
          name: mock-kmsv2-provider
          endpoint: unix:///tmp/kms.socket
          timeout: 3s
  ```

### Enabling KMS Features

#### Key Rotation
- In order for key rotation to take place, you will need to deploy another plugin on the control panel as a static pod
- First entry tells the apiserver what to encrypt new secrets with
- Second entry tells the apiserver what to decrypt the old secrets with
  ```
  providers:
    - kms:
      ...
    - kms:
      ...
  ```
- As the change in the encryption method only applies for subsequent secrets, in order to enforce the decryption and re-encryption of existing secrets, you will need to run the following code:
  ```kubectl get secrets --all-namespaces -o json | kubectl replace -f -```

#### Decryption
- To decrypt all existing secrets, you will need to edit the encryption config and apply the decryption to all secrets
- ```identity: {}``` means that no encryption used
  ```
  providers:
    - identity: {}
    - kms:
      ...
  ```